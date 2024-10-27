package test

import (
	"context"
	"testing"
	"time"

	"github.com/nbd-wtf/go-nostr"
)

func TestRelay(t *testing.T) {
	relay, err := nostr.RelayConnect(context.Background(), "ws://localhost:3334")
	if err != nil {
		t.Fatalf("Failed to connect to relay: %v", err)
	}
	defer relay.Close()

	t.Logf("Connected to relay successfully")

	// Test valid event publication
	t.Run("PublishValidEvent", func(t *testing.T) {
		sk := nostr.GeneratePrivateKey()
		pk, _ := nostr.GetPublicKey(sk)
		testPublishEvent(t, relay, sk, pk, 1, validTags(), "Valid content", true)
	})

	// Test VerifyMessagePattern policy
	t.Run("VerifyMessagePattern", func(t *testing.T) {
		sk := nostr.GeneratePrivateKey()
		pk, _ := nostr.GetPublicKey(sk)

		// Test missing required tag
		invalidTags := validTags()
		invalidTags = append(invalidTags[:1], invalidTags[2:]...) // Remove second tag (location)
		testPublishEvent(t, relay, sk, pk, 1, invalidTags, "Missing required tag", false)

		// Test invalid tag value
		invalidTags = validTags()
		invalidTags[0] = nostr.Tag{"title", string(make([]byte, 101))} // Title too long
		testPublishEvent(t, relay, sk, pk, 1, invalidTags, "Invalid tag value", false)

		// Test valid tags and content
		testPublishEvent(t, relay, sk, pk, 1, validTags(), "Valid content", true)
	})
}

func testPublishEvent(t *testing.T, relay *nostr.Relay, sk string, pk string, kind int, tags nostr.Tags, content string, expectSuccess bool) {
	ev := nostr.Event{
		PubKey:    pk,
		CreatedAt: nostr.Timestamp(time.Now().Unix()),
		Kind:      kind,
		Tags:      tags,
		Content:   content,
	}

	t.Logf("Publishing event: Kind=%d, Content='%s', Tags=%v", ev.Kind, ev.Content, ev.Tags)

	err := ev.Sign(sk)
	if err != nil {
		t.Fatalf("Failed to sign event: %v", err)
	}

	err = relay.Publish(context.Background(), ev)
	if expectSuccess {
		if err != nil {
			t.Errorf("Expected successful publish, but got error: %v", err)
			t.Logf("Failed event: %+v", ev)
		} else {
			t.Logf("Successfully published event: Kind=%d", ev.Kind)
		}
	} else {
		if err == nil {
			t.Errorf("Expected publish to fail, but it succeeded for Kind=%d", ev.Kind)
		} else {
			t.Logf("Expected failure occurred for Kind=%d: %v", ev.Kind, err)
		}
	}

	sub, err := relay.Subscribe(context.Background(), nostr.Filters{
		{
			Authors: []string{pk},
			Kinds:   []int{kind},
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
		} else {
			t.Logf("Successfully received published event")
		}
	case <-time.After(5 * time.Second):
		if expectSuccess {
			t.Errorf("Timeout: No event received")
		} else {
			t.Logf("No event received as expected")
		}
	}
}

func validTags() nostr.Tags {
	return nostr.Tags{
		{"title", "Software Developer"},
		{"location", "Remote"},
		{"employment-time", "full-time"},
		{"industry", "Technology"},
		{"key-words", "golang,nostr,development"},
		{"salary", "Competitive"},
	}
}
