package main

import (
	"fmt"
	"github.com/songgao/water"
	"log"
	"net"
	"os"
)

func main() {

	if len(os.Args) != 4 {
		fmt.Printf("\nUsage: client IP PORT PASSWORD\n\n")
		return
	}

	//Connect udp
	ip := os.Args[1]
	port := os.Args[2]
	// password := os.Args[3]

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		return
	}
	defer conn.Close()

	log.Printf("Connected to server\n")

	ifce, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Interface Name: %s\n", ifce.Name())

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
			c2 <- buffer[:n]
		}
	}()

	log.Printf("Listening for events\n")

	for {
		select {
		case packet := <-c1:
			conn.Write(packet)
		case buffer := <-c2:
			ifce.Write(buffer)
		}
	}
}
