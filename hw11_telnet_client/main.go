package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"
)

var (
	ErrUnknownConnectionData = errors.New("unknown connection data")
	ErrConnectionToServer    = errors.New("error connection to server")
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()
	host := flag.Arg(0)
	port := flag.Arg(1)

	if host == "" || port == "" {
		log.Fatal(ErrUnknownConnectionData)
	}

	telnetClient := NewTelnetClient(host+":"+port, *timeout, os.Stdin, os.Stdout)

	if err := telnetClient.Connect(); err != nil {
		log.Fatal(ErrConnectionToServer)
		return
	}
	defer telnetClient.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		if err := telnetClient.Send(); err != nil {
			telnetClient.Close()
			cancel()
			log.Printf("error send operation: %s\n", err)
			return
		}
	}()

	go func() {
		if err := telnetClient.Receive(); err != nil {
			os.Stdin.Close()
			cancel()
			log.Printf("error receive operation: %s\n", err)
			return
		}
	}()

	<-ctx.Done()
}
