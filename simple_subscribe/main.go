package main

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func printMsg(m *nats.Msg, i int) {
	log.Printf("[#%d] Received on [%s]: '%s'", i, m.Subject, string(m.Data))
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// listen for interrupts to exit gracefully
		sigChannel := make(chan os.Signal, 1)
		signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)
		<-sigChannel
		close(sigChannel)
		cancel()
	}()

	// register the consumer
	go consumer(ctx, "hello.world")
	go consumer(ctx, "goodbye.world")

	<-ctx.Done()

	log.Println("exiting gracefully")
}

func consumer(ctx context.Context, subj string) {

	nc, err := nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		log.Fatal("Failed to connect to NATS server:", err)
	}
	defer nc.Close()

	fmt.Println("Connected to NATS server on port 4222")

	messages := make(chan *nats.Msg, 1000)

	// we're subscribing to the subject
	// and assigning our channel as reference to receive messages there
	subscription, err := nc.ChanSubscribe(subj, messages)
	if err != nil {
		log.Fatal("Failed to subscribe to subject:", err)
	}

	defer func() {
		subscription.Unsubscribe()
		close(messages)
	}()

	log.Println("Subscribed to", subj)

	for {
		select {
		case <-ctx.Done():
			log.Println("exiting from consumer")
			return
		case msg := <-messages:
			log.Println("received", string(msg.Data))
		}
	}
}
