package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage %s server_addr:server_port ...\n", os.Args[0])
		return
	}

	conn, err := net.ListenPacket("udp", "")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	var serverResolvedAddr net.Addr
	const retries = 10
	var addressBuf [21]byte // xxx.xxx.xxx.xxx:yyyyy

	for _, serverStrAddr := range os.Args[1:] {
		serverResolvedAddr, err = net.ResolveUDPAddr("udp", serverStrAddr)
		if err != nil {
			log.Printf("could not resolve addr (%s): %v", serverStrAddr, err)
			continue
		}

		var i int
		for i = 0; i < retries; i++ {
			_, err = conn.WriteTo(nil, serverResolvedAddr)
			if err != nil {
				log.Printf("could not write to addr (%s): %v", serverStrAddr, err)
				continue
			}

			conn.SetReadDeadline(time.Now().Add(2 * time.Second))
			n, senderAddr, err := conn.ReadFrom(addressBuf[:])
			if err != nil {
				continue
			}

			if senderAddr.String() != serverResolvedAddr.String() {
				continue
			}

			fmt.Printf("%s (from %s)\n", string(addressBuf[:n]), serverStrAddr)
			break
		}

		if i == retries {
			fmt.Println("could not receive from", serverStrAddr)
		}
	}
}
