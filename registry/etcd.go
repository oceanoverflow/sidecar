package registry

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/spf13/viper"

	"github.com/oceanoverflow/sidecar/loadbalancing"
)

var (
	rootPath = "dubbomesh"
)

// Client used for service registry and discovery, also has the ability for load balancing
type Client struct {
	etcdClient client.Client
	connected  bool
	nodes      loadbalancing.WeightedServers
}

var c *Client
var once sync.Once

// New return an instance of etcd client
func New() *Client {
	once.Do(func() {
		c = &Client{}
	})

	if c.connected {
		log.Println("Can't connect twice")
		return nil
	}

	var endpoints []string
	single := viper.GetString("etcd")
	if single == "" {
		log.Println("Please specify the address of etcd")
		return nil
	}
	endpoints = append(endpoints, single)

	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	var err error
	c.etcdClient, err = client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

// Connect let the consumer agent know the ip address of the three provider agent
// /dubbomesh/com.some.package.IHelloService/
func (c *Client) Connect(serviceName string) error {
	kapi := client.NewKeysAPI(c.etcdClient)
	path := fmt.Sprintf("/%s/%s", rootPath, serviceName)

	resp, err := kapi.Get(context.Background(), path, nil)
	if err != nil {
		return err
	} else {
		if resp.Node.Dir {
			for _, peer := range resp.Node.Nodes {
				s := peer.Value
				switch strings.Split(s, ":")[0] {
				case "provider-small":
					c.nodes.Add(s, 1)
				case "provider-medium":
					c.nodes.Add(s, 2)
				case "provider-large":
					c.nodes.Add(s, 3)
				}
			}
		}
	}

	watcher := kapi.Watcher(path, &client.WatcherOptions{Recursive: true})
	go c.watch(watcher)
	c.connected = true
	return nil
}

// Register register the provider's service address on etcd
// /dubbomesh/com.some.package.IHelloService/provider-small:20000
func (c *Client) Register(serviceName, host, port string) error {
	kapi := client.NewKeysAPI(c.etcdClient)
	s := fmt.Sprintf("/%s/%s/%s:%s", rootPath, serviceName, host, port)

	resp, err := kapi.Set(context.Background(), s, "", nil)
	if err != nil {
		return err
	}
	log.Printf("serviceName is registered on etcd, resp is %q\n", resp)
	return nil
}

// Next is a simple wrapper for returning the next server for load balancing
func (c *Client) Next() string {
	return c.nodes.Next()
}

func (c *Client) watch(watcher client.Watcher) {
	for {
		resp, err := watcher.Next(context.Background())
		if err == nil {
			if resp.Action == "set" {
				s := resp.Node.Value
				switch strings.Split(s, ":")[0] {
				case "provider-small":
					c.nodes.Add(s, 1)
				case "provider-medium":
					c.nodes.Add(s, 2)
				case "provider-large":
					c.nodes.Add(s, 3)
				}
			}
		}

		// stop watch when all nodes is get
		if c.nodes.Len() == 3 {
			return
		}
	}
}
