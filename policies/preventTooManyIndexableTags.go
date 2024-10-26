package policies

import (
	"context"

	"github.com/nbd-wtf/go-nostr"
)

const maxIndexableTags = 10

func PreventTooManyIndexableTags() func(context.Context, *nostr.Event) (bool, string) {
	return func(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
		ntags := 0
		for _, tag := range event.Tags {
			if len(tag) > 0 && len(tag[0]) == 1 {
				ntags++
			}
		}
		if ntags > maxIndexableTags {
			return true, "too many indexable tags"
		}
		return false, ""
	}
}
