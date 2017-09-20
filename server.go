package main

import (
	"fmt"
	"github.com/songgao/water"
	"log"
	"net"
	"os"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Printf("\nUsage: server PORT PASSWORD\n\n")
		return
	}

	//Connect udp
	port := os.Args[1]
	// password := os.Args[3]

	ServerAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Printf("Err\n")
	}

	conn, err := net.ListenUDP("udp", ServerAddr)
	if err != nil {
		log.Printf("Err\n")
	}
	defer conn.Close()

	ifce, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Interface Name: %s\n", ifce.Name())

	hibuffer := make([]byte, 1024)
	_, addr, _ := conn.ReadFromUDP(hibuffer)

	c1 := make(chan []byte)

	go func() {
		packet := make([]byte, 20000)
		for {
			n, err := ifce.Read(packet)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Packet Received: % x\n", packet[:n])
			c1 <- packet[:n]

		}
	}()

	c2 := make(chan []byte)
	go func() {
		for {
			buffer := make([]byte, 20000)
			n, err := conn.Read(buffer)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("UDP Packet Received: % x\n", buffer[:n])
			c2 <- buffer[:n]
		}
	}()

	log.Printf("Listening for events\n")

	for {
		select {
		case packet := <-c1:
			conn.WriteToUDP(packet, addr)
		case buffer := <-c2:
			ifce.Write(buffer)
		}
	}
}
