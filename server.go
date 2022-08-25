package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Printf("usage %s [bind_addr:bind_port]\n", os.Args[0])
		return
	}

	bindAddr := ":5555"
	if len(os.Args) == 2 {
		bindAddr = os.Args[1]
	}

	conn, err := net.ListenPacket("udp", bindAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	log.Println("Serving on", bindAddr)

	for {
		_, addr, err := conn.ReadFrom(nil)
		if err != nil {
			log.Printf("could not read from addr (%s): %v", addr.String(), err)
			continue
		}
		log.Println("received connection from", addr.String())

		_, err = conn.WriteTo([]byte(addr.String()), addr)
		if err != nil {
			log.Printf("could not write to addr (%s): %v", addr.String(), err)
		}
	}
}
