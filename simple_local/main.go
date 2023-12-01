package simple

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {

	url := "nats://localhost:4222"
	if url == "" {
		url = nats.DefaultURL
	}

	nc, _ := nats.Connect(url)

	defer nc.Drain()

	nc.Publish("greet.joe", []byte("hello"))

	sub, _ := nc.SubscribeSync("greet.*")

	msg, _ := sub.NextMsg(10 * time.Millisecond)
	fmt.Println("subscribed after a publish...")
	fmt.Printf("msg is nil? %v\n", msg == nil)

	nc.Publish("greet.joe", []byte("hello"))
	nc.Publish("greet.pam", []byte("hello"))

	msg, _ = sub.NextMsg(10 * time.Millisecond)
	fmt.Printf("msg data: %q on subject %q\n", string(msg.Data), msg.Subject)

	msg, _ = sub.NextMsg(10 * time.Millisecond)
	fmt.Printf("msg data: %q on subject %q\n", string(msg.Data), msg.Subject)

	nc.Publish("greet.bob", []byte("hello"))

	msg, _ = sub.NextMsg(10 * time.Millisecond)
	fmt.Printf("msg data: %q on subject %q\n", string(msg.Data), msg.Subject)
}
