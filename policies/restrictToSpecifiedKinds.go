package policies

import (
	"context"
	"fmt"
	"log"
	"slices"

	"github.com/nbd-wtf/go-nostr"
)

func RestrictToSpecifiedKinds(allowEphemeral bool, kinds ...uint16) func(context.Context, *nostr.Event) (bool, string) {
	slices.Sort(kinds)
	log.Printf("RestrictToSpecifiedKinds initialized with kinds: %v, allowEphemeral: %v", kinds, allowEphemeral)

	return func(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
		log.Printf("RestrictToSpecifiedKinds: Checking event kind %d", event.Kind)

		if allowEphemeral && event.IsEphemeral() {
			log.Printf("RestrictToSpecifiedKinds: Allowing ephemeral event")
			return false, ""
		}

		for _, kind := range kinds {
			if uint16(event.Kind) == kind {
				log.Printf("RestrictToSpecifiedKinds: Event kind %d is allowed", event.Kind)
				return false, ""
			}
		}

		log.Printf("RestrictToSpecifiedKinds: Event kind %d is not allowed", event.Kind)
		return true, fmt.Sprintf("received event kind %d not allowed", event.Kind)
	}
}
