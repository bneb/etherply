package pubsub_test

import (
	"testing"
	"time"

	"github.com/bneb/etherply/etherply-sync-server/internal/pubsub"
)

func TestMemoryPubSub(t *testing.T) {
	ps := pubsub.NewMemoryPubSub()
	topic := "topic-1"

	// 1. Subscribe
	rx1, unsub1 := ps.Subscribe(topic)
	rx2, unsub2 := ps.Subscribe(topic)

	defer unsub1()
	defer unsub2()

	// 2. Publish
	msg := pubsub.Message{
		Topic:   topic,
		Payload: []byte("hello"),
	}

	if err := ps.Publish(topic, msg); err != nil {
		t.Fatalf("Publish failed: %v", err)
	}

	// 3. Verify Receipt (with timeout)
	timeout := time.After(1 * time.Second)
	var count int
	for count < 2 {
		select {
		case m := <-rx1:
			if string(m.Payload) != "hello" {
				t.Errorf("Sub1 got wrong payload: %s", m.Payload)
			}
			count++
		case m := <-rx2:
			if string(m.Payload) != "hello" {
				t.Errorf("Sub2 got wrong payload: %s", m.Payload)
			}
			count++
		case <-timeout:
			t.Fatal("Timeout waiting for messages")
		}
	}

	// 4. Test Unsubscribe
	unsub1()

	// Publish again
	msg2 := pubsub.Message{Topic: topic, Payload: []byte("world")}
	ps.Publish(topic, msg2)

	select {
	case <-rx1:
		t.Error("Sub1 received message after unsubscribe")
	case m := <-rx2:
		if string(m.Payload) != "world" {
			t.Error("Sub2 wrong payload")
		}
	case <-time.After(100 * time.Millisecond):
		// Expected for rx1 not to receive, but rx2 SHOULD receive.
		// If rx2 received above, we are good.
	}
}
