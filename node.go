package main

import (
	"fmt"
	"net"
	"os"
	"net/rpc"
	"log"
	"strconv"
)

func main() {
	fmt.Println("hello world")
	var ip_addr string
	if (len(os.Args)>1){
		ip_addr = os.Args[1]
	}else{
		ip_addr = ":8080"
	}
	// Start Listener
	udp_addr, _ := net.ResolveUDPAddr("udp", ip_addr)
    client, err := net.ListenUDP("udp", udp_addr)
    if err != nil {
     	panic(err)
	}
	defer client.Close()

	go RunListener(client)


	// Connect to server
	serverConn, err := rpc.Dial("tcp", ":8081")
	if err != nil {
		panic(err)
	}
	localIP := GetLocalMinerIP()
	var response []string

	// Get IP from server
	err = serverConn.Call("GServer.Register", localIP, &response)
	if err != nil{
		panic(err)
	}
	if len(response)>0{
		for ind, val := range response {
			fmt.Println(strconv.Itoa(ind) + ": " + val)
		}
	}

	// Connect to other node
	// Exchange messages with other node
	select {}
}

func GetLocalMinerIP() string {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}

	// https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
	for _, address := range addresses {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return ""
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
