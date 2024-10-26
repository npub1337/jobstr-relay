package policies

import (
	"context"

	"github.com/nbd-wtf/go-nostr"
)

const maxTagValueLen = 40

func PreventLargeTags() func(context.Context, *nostr.Event) (bool, string) {
	return func(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
		for _, tag := range event.Tags {
			if len(tag) > 1 && len(tag[0]) == 1 {
				if len(tag[1]) > maxTagValueLen {
					return true, "event contains too large tags"
				}
			}
		}
		return false, ""
	}
}
