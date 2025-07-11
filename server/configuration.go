package main

import (
	"fmt"
	"strings"
)

type configuration struct {
	DigitalSambaAPIKey          string
	DigitalSambaDashboardURL    string
	DigitalSambaTeamName        string
	DigitalSambaEmbedded        bool
	DigitalSambaShowPrejoinPage bool
	DigitalSambaNamingScheme    string
	DigitalSambaRoomExpiry      int
	DigitalSambaMaxParticipants int
	DigitalSambaEnableRecording bool
	DigitalSambaEnableBreakoutRooms bool
}

func (c *configuration) IsValid() error {
	if c.DigitalSambaAPIKey == "" {
		return fmt.Errorf("DigitalSamba API Key is required")
	}

	if c.DigitalSambaDashboardURL == "" {
		return fmt.Errorf("DigitalSamba Dashboard URL is required")
	}

	// Validate URL format
	dashboardURL := strings.TrimSpace(c.DigitalSambaDashboardURL)
	if !strings.HasPrefix(dashboardURL, "http://") && !strings.HasPrefix(dashboardURL, "https://") {
		return fmt.Errorf("DigitalSamba Dashboard URL must start with http:// or https://")
	}

	// Validate room expiry
	if c.DigitalSambaRoomExpiry < 0 {
		return fmt.Errorf("room expiry time cannot be negative")
	}

	// Validate max participants
	if c.DigitalSambaMaxParticipants < 1 || c.DigitalSambaMaxParticipants > 2000 {
		return fmt.Errorf("maximum participants must be between 1 and 2000")
	}

	// Validate naming scheme
	validSchemes := []string{"words", "uuid", "mattermost", "ask"}
	valid := false
	for _, scheme := range validSchemes {
		if c.DigitalSambaNamingScheme == scheme {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid naming scheme: %s", c.DigitalSambaNamingScheme)
	}

	return nil
}

func (c *configuration) GetDashboardURL() string {
	url := strings.TrimSpace(c.DigitalSambaDashboardURL)
	return strings.TrimRight(url, "/")
}