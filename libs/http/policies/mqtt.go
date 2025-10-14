package mqttPolicies

import (
	"fmt"
	"strings"
)

type Policy struct {
	AllowedSubscriptionTopics []string `json:"allowedSubscriptionTopics" validate:"required,dive,mqtt_topic"`
	AllowedPublishTopics      []string `json:"allowedPublishTopics" validate:"required,dive,mqtt_topic"`
}

// Checks if the policy permits the action on the topic
func (p Policy) Assert(action string, topic string) error {
	// Check action and topic against the policy
	allow := false
	if action == "subscribe" {
		for _, allowedTopic := range p.AllowedSubscriptionTopics {
			if topicMatches(topic, allowedTopic) {
				allow = true
				break
			}
		}
	}

	if action == "publish" {
		for _, allowedTopic := range p.AllowedPublishTopics {
			if topicMatches(topic, allowedTopic) {
				allow = true
				break
			}
		}
	}

	if !allow {
		return fmt.Errorf("the certificate is not allowed to %s to topic %s", action, topic)
	}

	return nil
}

// Checks if a topic matches a pattern according to MQTT wildcards rules
func topicMatches(topic, pattern string) bool {
	// Split both strings into parts
	topicParts := strings.Split(topic, "/")
	patternParts := strings.Split(pattern, "/")

	// If pattern ends with #, remove it and check if topic starts with remaining pattern
	if strings.HasSuffix(pattern, "/#") {
		patternParts = patternParts[:len(patternParts)-1]
		return len(topicParts) >= len(patternParts) &&
			topicMatchesPattern(topicParts[:len(patternParts)], patternParts)
	}

	// For exact matching (including + wildcards), both must have same number of parts
	if len(topicParts) != len(patternParts) {
		return false
	}

	return topicMatchesPattern(topicParts, patternParts)
}

// Compares topic parts with pattern parts
func topicMatchesPattern(topicParts, patternParts []string) bool {
	for i := range patternParts {
		// + matches any single level
		if patternParts[i] == "+" {
			continue
		}

		if patternParts[i] != topicParts[i] {
			return false
		}
	}
	return true
}
