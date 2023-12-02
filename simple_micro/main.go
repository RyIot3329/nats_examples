package main

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"log"
	"os"
	"os/signal"
	"syscall"
)

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
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Close()

	echoHandler := func(req micro.Request) {
		req.Respond(req.Data())
	}

	config := micro.Config{
		Name:        "EchoService",
		Version:     "1.0.0",
		Description: "Send back what you receive",
		// DoneHandler can be set to customize behavior on stopping a service.
		DoneHandler: func(srv micro.Service) {
			info := srv.Info()
			fmt.Printf("stopped service %q with ID %q\n", info.Name, info.ID)
		},

		// ErrorHandler can be used to customize behavior on service execution error.
		ErrorHandler: func(srv micro.Service, err *micro.NATSError) {
			info := srv.Info()
			fmt.Printf("Service %q returned an error on subject %q: %s", info.Name, err.Subject, err.Description)
		},

		// optional base handler
		Endpoint: &micro.EndpointConfig{
			Subject: "echo",
			Handler: micro.HandlerFunc(echoHandler),
		},
	}

	srv, err := micro.AddService(nc, config)
	if err != nil {
		log.Fatal(err)
	}
	defer srv.Stop()
	<-ctx.Done()

}

func ExampleAddService() {
	nc, err := nats.Connect("127.0.0.1:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	echoHandler := func(req micro.Request) {
		req.Respond(req.Data())
	}

	config := micro.Config{
		Name:        "EchoService",
		Version:     "1.0.0",
		Description: "Send back what you receive",
		// DoneHandler can be set to customize behavior on stopping a service.
		DoneHandler: func(srv micro.Service) {
			info := srv.Info()
			fmt.Printf("stopped service %q with ID %q\n", info.Name, info.ID)
		},

		// ErrorHandler can be used to customize behavior on service execution error.
		ErrorHandler: func(srv micro.Service, err *micro.NATSError) {
			info := srv.Info()
			fmt.Printf("Service %q returned an error on subject %q: %s", info.Name, err.Subject, err.Description)
		},

		// optional base handler
		Endpoint: &micro.EndpointConfig{
			Subject: "echo",
			Handler: micro.HandlerFunc(echoHandler),
		},
	}

	srv, err := micro.AddService(nc, config)
	if err != nil {
		log.Fatal(err)
	}
	defer srv.Stop()
}
