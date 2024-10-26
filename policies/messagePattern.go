package policies

import (
	"context"
	"strings"

	"github.com/nbd-wtf/go-nostr"
)

func VerifyMessagePattern() func(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
	return func(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
		lines := strings.Split(strings.TrimSpace(event.Content), "\n")
		if len(lines) != 2 ||
			!strings.HasPrefix(lines[0], "Job Title: ") ||
			!strings.HasPrefix(lines[1], "description: ") {
			return true, "Event rejected: invalid format. Use 'Job Title:' on the first line and 'description:' on the second line."
		}
		return false, ""
	}
}
