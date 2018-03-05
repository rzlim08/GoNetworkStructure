package main

import (
	"fmt"
	"net"
	"os"
	"net/rpc"
	"strconv"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)


func main() {
	fmt.Println("hello world")
	var ip_addr string
	if (len(os.Args)>1){
		ip_addr = os.Args[1]
	}else{
		ip_addr = "127.0.0.1:0"
	}
	udp_addr, client := startListener(ip_addr)
	defer client.Close()

	go RunListener(client)

	otherNodes := serverRegister(client.LocalAddr().String())

	floodNodes(otherNodes, udp_addr)

	pixelgl.Run(run)
	select {}
}
func run() {
	// all of our code will be fired up from here
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	for !win.Closed() {
		win.Update()
	}
}
func floodNodes(otherNodes []string, udp_addr *net.UDPAddr) {
	for _, ip := range otherNodes {
		node_udp, _ := net.ResolveUDPAddr("udp", ip)
		// Connect to other node
		node_client, err := net.DialUDP("udp", udp_addr, node_udp)
		if err != nil {
			panic(err)
		}
		// Exchange messages with other node
		node_client.Write([]byte("write byte"))
	}
}

func serverRegister(localIP string) []string {
	// Connect to server
	serverConn, err := rpc.Dial("tcp", ":8081")
	if err != nil {
		panic(err)
	}
	var response []string
	// Get IP from server
	err = serverConn.Call("GServer.Register", localIP, &response)
	if err != nil {
		panic(err)
	}
	if len(response) > 0 {
		for ind, val := range response {
			fmt.Println(strconv.Itoa(ind) + ": " + val)
		}
	}
	return response
}

func startListener(ip_addr string) (*net.UDPAddr, *net.UDPConn) {
	// Start Listener
	udp_addr, _ := net.ResolveUDPAddr("udp", ip_addr)
	client, err := net.ListenUDP("udp", udp_addr)
	if err != nil {
		panic(err)
	}
	return udp_addr, client
}

func RunListener(client *net.UDPConn) {
	client.SetReadBuffer(1048576)
	i := 0
	for {
		i++
		buf := make([]byte, 1024)
		rlen, addr, err := client.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(buf[0:rlen]))
		fmt.Println(addr)
		fmt.Println(i)
	}
}

