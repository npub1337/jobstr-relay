package policies

import (
	"context"
	"fmt"
	"slices"

	"github.com/nbd-wtf/go-nostr"
)

var allowedKinds = []uint16{
	1,
}

func RestrictToSpecifiedKinds() func(context.Context, *nostr.Event) (bool, string) {
	slices.Sort(allowedKinds)

	return func(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
		if _, allowed := slices.BinarySearch(allowedKinds, uint16(event.Kind)); allowed {
			return false, ""
		}

		return true, fmt.Sprintf("received event kind %d not allowed", event.Kind)
	}
}
