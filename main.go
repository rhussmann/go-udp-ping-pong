package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
)

var tokenFound bool = false

func main() {

	h := os.Getenv("OTHER_PLAYER")
	if h == "" {
		fmt.Printf("Need to define OTHER_PLAYER\n")
		os.Exit(1)
	}

	lIP := os.Getenv("LISTEN_IP")
	lH := fmt.Sprintf("%s:1053", lIP)
	fmt.Printf("Listening on %s\n", lH)
	pc, err := net.ListenPacket("udp", lH)
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	host := fmt.Sprintf("%s:1053", os.Getenv("OTHER_PLAYER"))
	udpAddr, err := net.ResolveUDPAddr("udp4", host)
	if err != nil {
		panic(err)
	}

	g, _ := errgroup.WithContext(context.Background())
	g.Go(func() error {

		for {
			buf := make([]byte, 1024)
			fmt.Printf("Blocking read of UDP socket...\n")
			n, _, err := pc.ReadFrom(buf)
			if err != nil {
				fmt.Printf("Found error during read: %s\n", err)
				continue
			}
			fmt.Printf("Running serve...\n")
			go serve(pc, udpAddr, buf[:n])
			fmt.Printf("After go routine...\n")
		}

		return nil
	})

	g.Go(func() error {
		for !tokenFound {
			s1 := rand.NewSource(time.Now().UnixNano())
			r1 := rand.New(s1)
			wait := r1.Intn(100) % 10
			fmt.Printf("Waiting %d seconds for ping...\n", wait)
			time.Sleep(time.Duration(wait) * time.Second)

			if tokenFound {
				fmt.Printf("Token found during sleep, breaking\n")
				break
			}

			fmt.Printf("Token not found\n")
			fmt.Printf("Sending token to %s\n", host)
			if _, err = pc.WriteTo([]byte("Hello world!"), udpAddr); err != nil {
				panic(err)
			}
		}

		fmt.Printf("Token found!\n")
		return nil
	})

	if err := g.Wait(); err != nil {
		panic(err)
	}
}

func serve(pc net.PacketConn, addr net.Addr, buf []byte) {
	tokenFound = true
	fmt.Printf("Received token %s\n", string(buf))
	buf[2] |= 0x80
	pc.WriteTo(buf, addr)
	fmt.Printf("Send response\n")
}
