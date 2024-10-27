package policies

import (
	"context"
	"fmt"
	"strings"

	"github.com/nbd-wtf/go-nostr"
)

type TagDefinition struct {
	Required bool
	Validate func(value string) error
}

var definedTags = map[string]TagDefinition{
	"title": {
		Required: true, //Required
		Validate: func(value string) error {
			if len(value) > 100 {
				return fmt.Errorf("title must not exceed 100 characters")
			}
			return nil
		},
	},
	"location": {
		Required: true, //Required
		Validate: func(value string) error {
			if len(value) > 50 {
				return fmt.Errorf("location must not exceed 50 characters")
			}
			return nil
		},
	},
	"employment-time": {
		Required: true, //Required
		Validate: func(value string) error {
			validTypes := map[string]bool{"full-time": true, "part-time": true, "contract": true, "temporary": true}
			if !validTypes[value] {
				return fmt.Errorf("invalid employment-time: must be full-time, part-time, contract, or temporary")
			}
			return nil
		},
	},
	"industry": {
		Required: true, //Required
		Validate: func(value string) error {
			if len(value) > 50 {
				return fmt.Errorf("industry must not exceed 50 characters")
			}
			return nil
		},
	},
	"key-words": {
		Required: true, //Required
		Validate: func(value string) error {
			keywords := strings.Split(value, ",")
			if len(keywords) > 5 {
				return fmt.Errorf("key-words must not exceed 5 items")
			}
			for _, kw := range keywords {
				if len(strings.TrimSpace(kw)) > 20 {
					return fmt.Errorf("each key-word must not exceed 20 characters")
				}
			}
			return nil
		},
	},
	"salary": {
		Required: false, //Optional
		Validate: func(value string) error {
			return nil
		},
	},
}

func VerifyMessagePattern() func(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
	return func(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
		missingTags := []string{}
		invalidTags := []string{}

		for tagName, definition := range definedTags {
			found := false
			for _, tag := range event.Tags {
				if len(tag) >= 2 && tag[0] == tagName {
					found = true
					if err := definition.Validate(tag[1]); err != nil {
						invalidTags = append(invalidTags, fmt.Sprintf("%s: %s", tagName, err.Error()))
					}
					break
				}
			}
			if !found && definition.Required {
				missingTags = append(missingTags, tagName)
			}
		}

		if len(missingTags) > 0 || len(invalidTags) > 0 {
			errorMsg := ""
			if len(missingTags) > 0 {
				errorMsg += "Missing required tags: " + strings.Join(missingTags, ", ") + ". "
			}
			if len(invalidTags) > 0 {
				errorMsg += "Invalid tags: " + strings.Join(invalidTags, "; ") + "."
			}
			return true, "Event rejected: " + errorMsg
		}

		return false, ""
	}
}
