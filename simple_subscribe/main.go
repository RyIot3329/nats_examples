package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/nats-io/nats.go"
)

func printMsg(m *nats.Msg, i int) {
	fmt.Printf("[#%d] Received on [%s]: '%s'", i, m.Subject, string(m.Data))
}

func main() {

	url := "nats://127.0.0.1:4222"
	if url == "" {
		url = nats.DefaultURL
	}

	nc, _ := nats.Connect(url)

	defer nc.Drain()

	i := 0
	sbj := "hello.world"

	nc.Subscribe(sbj, func(msg *nats.Msg) {
		fmt.Println("got message")
		i += 1
		printMsg(msg, i)
	})
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on [%s]", sbj)
	// if *showTime {
	// 	log.SetFlags(log.LstdFlags)
	// }

	runtime.Goexit()
}
