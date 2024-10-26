package test

import (
	"context"
	"testing"
	"time"

	"github.com/nbd-wtf/go-nostr"
)

func TestRelay(t *testing.T) {
	sk := nostr.GeneratePrivateKey()
	pk, _ := nostr.GetPublicKey(sk)
	t.Logf("Generated key pair. Public key: %s", pk)

	relay, err := nostr.RelayConnect(context.Background(), "ws://localhost:3334")
	if err != nil {
		t.Fatalf("Failed to connect to relay: %v", err)
	}
	defer relay.Close()

	// Test case 1: Invalid format
	testPublishEvent(t, relay, sk, pk, "Invalid format message", false)

	// Test case 2: Valid format
	testPublishEvent(t, relay, sk, pk, "Job Title: Software Developer\ndescription: Exciting opportunity for a skilled developer", true)

	// Test case 3: Another invalid format
	testPublishEvent(t, relay, sk, pk, "Job Title: Software Developer\nInvalid second line", false)
}

func testPublishEvent(t *testing.T, relay *nostr.Relay, sk string, pk string, content string, expectSuccess bool) {
	ev := nostr.Event{
		PubKey:    pk,
		CreatedAt: nostr.Timestamp(time.Now().Unix()),
		Kind:      1,
		Tags:      nil,
		Content:   content,
	}

	err := ev.Sign(sk)
	if err != nil {
		t.Fatalf("Failed to sign event: %v", err)
	}

	err = relay.Publish(context.Background(), ev)
	if expectSuccess {
		if err != nil {
			t.Errorf("Expected successful publish, but got error: %v", err)
		}
	} else {
		if err == nil {
			t.Errorf("Expected publish to fail, but it succeeded")
		}
	}

	sub, err := relay.Subscribe(context.Background(), nostr.Filters{
		{
			Authors: []string{pk},
			Kinds:   []int{1},
			Limit:   1,
		},
	})
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	select {
	case receivedEv := <-sub.Events:
		if !expectSuccess {
			t.Errorf("Received unexpected event: %v", receivedEv)
		}
	case <-time.After(5 * time.Second):
		if expectSuccess {
			t.Errorf("Timeout: No event received")
		}
	}
}
