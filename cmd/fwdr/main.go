package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	name            string
	source          string
	destinationPort int
	dialTimeout     time.Duration
)

func init() {
	flag.StringVar(&name, "name", "", "the connection name to identify the connection")
	flag.StringVar(&source, "source", "", "the connection string (<host>:<port>)")
	flag.IntVar(&destinationPort, "destination-port", 8080, "the destination port to forward to the host")
	flag.DurationVar(&dialTimeout, "dial-timeout", 30*time.Second, "timeout for connecting to source")
}

func main() {
	flag.Parse()

	if name == "" {
		log.Fatal("name is required")
	}

	if source == "" {
		log.Fatal("source is required")
	}

	if destinationPort <= 0 {
		log.Fatal("destination-port is required")
	}

	incoming, err := net.Listen("tcp", fmt.Sprintf(":%d", destinationPort))
	if err != nil {
		log.Fatalf("could not listen to destination port %d: %v", destinationPort, err)
	}
	log.Printf("connection %s listening on port: %d\n", name, destinationPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("shutting down")
		if err := incoming.Close(); err != nil {
			log.Println("error closing listener:", err)
		}
	}()

	var wg sync.WaitGroup
	for {
		client, err := incoming.Accept()
		if err != nil {
			select {
			case <-quit:
			default:
				log.Println("accept error:", err)
			}
			break
		}
		fmt.Printf("client '%v' connected!\n", client.RemoteAddr())

		wg.Go(func() {
			handleConnection(client, source, name)
		})
	}

	wg.Wait()
	log.Println("all connections closed")
}

func forward(src, dest net.Conn) {
	defer src.Close()  //nolint:errcheck
	defer dest.Close() //nolint:errcheck
	_, _ = io.Copy(src, dest)
}

func handleConnection(client net.Conn, source, name string) {
	defer client.Close() //nolint:errcheck

	remote, err := net.DialTimeout("tcp", source, dialTimeout)
	if err != nil {
		log.Printf("could not connect to source: %v", err)
		return
	}
	defer remote.Close() //nolint:errcheck
	fmt.Printf("connection to remote %s at %v established!\n", name, remote.RemoteAddr())

	var wg sync.WaitGroup
	wg.Go(func() { forward(client, remote) })
	wg.Go(func() { forward(remote, client) })
	wg.Wait()
}
