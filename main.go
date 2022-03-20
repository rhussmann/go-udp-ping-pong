package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {

	pc, err := net.ListenPacket("udp", ":1053")
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	fmt.Printf("Listening...\n")
	for {
		buf := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			continue
		}

		fmt.Printf("Before go routine\n")
		go serve(pc, addr, buf[:n])
		fmt.Printf("After go routine...\n")
	}

	return
}

func serve(pc net.PacketConn, addr net.Addr, buf []byte) {

	fmt.Printf("Sleeping 10 seconds...\n")
	time.Sleep(10 * time.Second)
	fmt.Printf("Received token %s\n", string(buf))
	buf[2] |= 0x80
	pc.WriteTo(buf, addr)
	fmt.Printf("Send response\n")
}
