package policies

import (
	"context"
	"fmt"
	"slices"

	"github.com/nbd-wtf/go-nostr"
)

func RestrictToSpecifiedKinds(kinds ...uint16) func(context.Context, *nostr.Event) (bool, string) {
	slices.Sort(kinds)

	return func(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
		if _, allowed := slices.BinarySearch(kinds, uint16(event.Kind)); allowed {
			return false, ""
		}

		return true, fmt.Sprintf("received event kind %d not allowed", event.Kind)
	}
}
