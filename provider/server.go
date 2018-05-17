package provider

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"strings"

	"github.com/spf13/viper"

	"github.com/oceanoverflow/sidecar/codec"
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
	connPool    pool.Pool
	leakyBuffer *pool.LeakyBuffer
)

func ServeCommunicate(size, port, dubbo string) {
	client := registry.New()
	client.Put("com.some.package.IHelloService", size, port)

	factory := func() (net.Conn, error) { return net.Dial("tcp", dubbo) }

	initial = viper.GetInt("connpool.initial")
	max = viper.GetInt("connpool.max")
	connPool, err := pool.New(initial, max, factory)
	if err != nil {
		log.Fatal("error creating connection pool")
	}
	n = viper.GetInt("leakybuffer.n")
	bufSize = viper.GetInt("leakybuffer.bufSize")
	leakyBuffer = pool.NewLeakyBuffer(n, bufSize)

	listenPort := viper.GetString("provider." + size)
	ln, err := net.Listen("tcp", listenPort)
	if err != nil {
		log.Fatalln("Unable to bind to the specific port")
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalln("Unable to accept connection")
		}
		go handleConnection(conn, dubbo, connPool)
	}
}

func handleConnection(conn net.Conn, dubbo string, connPool pool.Pool) {
	info, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println("error reading from the net.Conn")
		return
	}
	info = info[:len(info)-1]
	slices := strings.Split(info, "-")
	r := codec.NewRequest()
	r.Arguments = []byte(slices[3])
	payload := r.Encode()
	result := call(dubbo, payload, connPool)
	dubboResponse, err := codec.Read(bytes.NewReader(result))
	if err != nil {
		log.Println("error decoding result, can not parse")
		return
	}
	conn.Write(dubboResponse.Value)
}

func call(dubbo string, payload []byte, connPool pool.Pool) (result []byte) {
	conn, err := connPool.Get()
	if err != nil {
		log.Printf("error get connection from conn pool")
		return nil
	}
	defer conn.Close()

	conn.Write(payload)
	buf := leakyBuffer.Get()
	defer leakyBuffer.Put(buf)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("error calling RPC")
		return nil
	}
	result = buf[:n]
	return
}