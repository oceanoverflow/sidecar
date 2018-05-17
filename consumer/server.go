package consumer

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/spf13/viper"

	"github.com/oceanoverflow/sidecar/pool"
	"github.com/oceanoverflow/sidecar/registry"
)

var (
	n       int
	bufSize int
	initial int
	max     int
)

var (
	lb *pool.LeakyBuffer
)

// ListenAndServe is the gateway between HTTP and TCP
func ListenAndServe(port string) error {
	n = viper.GetInt("leakybuffer.n")
	bufSize = viper.GetInt("leakybuffer.bufSize")
	initial = viper.GetInt("connpool.initial")
	max = viper.GetInt("connpool.max")
	client := registry.New()
	client.Connect("com.some.package.IHelloService")
	lb = pool.NewLeakyBuffer(n, bufSize)
	http.Handle("/hash", &handler{
		client: client,
		pools:  make(map[string]pool.Pool),
	})
	return http.ListenAndServe(port, nil)
}

type handler struct {
	client *registry.Client
	pools  map[string]pool.Pool
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Println("Cannot handle method other than POST")
	}
	r.ParseForm()

	itfc := r.Form.Get("interface")
	method := r.Form.Get("method")
	paramTypes := r.Form.Get("parameterTypesString")
	param := r.Form.Get("parameter")

	result := h.communicate(itfc, method, paramTypes, param)
	w.Write([]byte(result))
	return
}

func (h *handler) communicate(interfaceName, method, parameterTypesString, parameter string) string {
	target := h.client.Next()
	if h.pools[target] == nil {
		factory := func() (net.Conn, error) { return net.Dial("tcp", target) }
		p, err := pool.New(initial, max, factory)
		if err != nil {
			log.Fatal("error creating connection pool")
		}
		h.pools[target] = p
	}
	conn, err := h.pools[target].Get()
	if err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%s-%s-%s-%s\n", interfaceName, method, parameterTypesString, parameter)
	_, err = conn.Write([]byte(s))
	if err != nil {
		panic(err)
	}
	b := lb.Get()
	defer lb.Put(b)
	n, err = conn.Read(b)
	if err != nil {
		panic(err)
	}
	return string(b[:n])
}
