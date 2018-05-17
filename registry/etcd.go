package registry

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/spf13/viper"

	"github.com/oceanoverflow/sidecar/loadbalancing"
	"github.com/oceanoverflow/sidecar/utils"
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

// New return an instance of etcd client
func New() *Client {
	c := &Client{}
	if c.connected {
		log.Println("can't connect twice")
		return nil
	}

	ep := []string{}
	single := viper.GetString("etcd")
	if single == "" {
		log.Println("Please specify the address of etcd")
		return nil
	}
	ep = append(ep, single)

	cfg := client.Config{
		Endpoints:               ep,
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
				ss := strings.Split(s, "-")
				switch ss[0] {
				case "small":
					c.nodes.Add(ss[1], 1)
				case "medium":
					c.nodes.Add(ss[1], 2)
				case "large":
					c.nodes.Add(ss[1], 3)
				}
			}
		}
	}

	watcher := kapi.Watcher(path, &client.WatcherOptions{Recursive: true})
	go c.watch(watcher)
	c.connected = true
	return nil
}

// Put let the provider agent publish their ip address
// /dubbomesh/com.some.package.IHelloService/192.168.100.100:2000
func (c *Client) Put(serviceName, size, port string) error {
	kapi := client.NewKeysAPI(c.etcdClient)
	s := fmt.Sprintf("/%s/%s/%s-%s:%s", rootPath, serviceName, size, utils.GetHostIP(), port)
	resp, err := kapi.Set(context.Background(), s, "", nil)
	if err != nil {
		return err
	}
	log.Printf("Set is done. Metadata is %q\n", resp)
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
				ss := strings.Split(s, "-")
				switch ss[0] {
				case "small":
					c.nodes.Add(ss[1], 1)
				case "medium":
					c.nodes.Add(ss[1], 2)
				case "large":
					c.nodes.Add(ss[1], 3)
				}
			}
		}
		// stop watching is all the key is get
		// here the max number is 3
		// this is hacky, modify this later
		if c.nodes.Len() == 3 {
			return
		}
	}
}
