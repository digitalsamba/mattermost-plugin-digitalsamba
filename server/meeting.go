package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/pluginapi/i18n"
)

type MeetingInfo struct {
	MeetingID string `json:"meeting_id"`
	RoomID    string `json:"room_id"`
	RoomURL   string `json:"room_url"`
	Token     string `json:"token"`
	TeamName  string `json:"team_name"`
	RoomName  string `json:"room_name"`
}

func boolPtr(b bool) *bool {
	return &b
}

func (p *Plugin) startMeeting(user *model.User, channel *model.Channel, meetingID string, meetingTopic string, _ bool, rootID string) (*MeetingInfo, error) {
	l := p.b.GetServerLocalizer()
	
	// Generate meeting ID if not provided
	if meetingID == "" {
		meetingID = p.generateMeetingID(user, channel, meetingTopic)
	}
	
	// Set default meeting topic if not provided
	defaultMeetingTopic := p.b.LocalizeDefaultMessage(l, &i18n.Message{
		ID:    "digitalsamba.start_meeting.default_meeting_topic",
		Other: "DigitalSamba Meeting",
	})
	
	if meetingTopic == "" {
		meetingTopic = defaultMeetingTopic
	}

	// Create room in DigitalSamba
	config := p.getConfiguration()
	roomExpiry := time.Now().Add(time.Duration(config.DigitalSambaRoomExpiry) * time.Minute)
	
	// Ensure friendly URL doesn't exceed 32 character limit
	friendlyURL := meetingID
	if len(friendlyURL) > 32 {
		friendlyURL = friendlyURL[:32]
	}
	
	createRoomReq := &CreateRoomRequest{
		Topic:             meetingTopic,
		FriendlyURL:       friendlyURL,
		Privacy:           "public",
		MaxParticipants:   config.DigitalSambaMaxParticipants,
		RecordingsEnabled: config.DigitalSambaEnableRecording,
		ChatEnabled:       true,
		
		// Join Settings - optimized for internal teams
		JoinScreenEnabled: boolPtr(false),  // Skip prejoin screen for faster access
		MuteOnJoin:        boolPtr(false),  // Don't mute by default for internal meetings
		CameraOffOnJoin:   boolPtr(false),  // Camera on by default for better engagement
		
		// Disable consent messages for internal use
		ConsentMessage:    "",  // Empty string disables consent message
		
		// Layout Settings
		DefaultLayout:     "auto",       // Smart layout based on content
		
		// Enable collaboration features
		EnableWhiteboard:  true,
		EnablePolling:     true,
		EnableQA:          true,
	}
	
	if config.DigitalSambaRoomExpiry > 0 {
		createRoomReq.ExpiresAt = &roomExpiry
	}
	
	room, err := p.digitalSambaClient.CreateRoom(createRoomReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create room: %w", err)
	}

	// Create moderator token for the meeting creator
	tokenReq := &CreateTokenRequest{
		RoomID:    room.ID,
		UserID:    user.Id,
		UserName:  user.GetDisplayName(model.ShowNicknameFullName),
		UserEmail: user.Email,
		Role:      "moderator", // All internal users are moderators
	}
	
	// Add avatar URL if available
	siteURL := p.API.GetConfig().ServiceSettings.SiteURL
	if siteURL != nil && *siteURL != "" {
		tokenReq.AvatarURL = fmt.Sprintf("%s/api/v4/users/%s/image?_=%d", *siteURL, user.Id, user.LastPictureUpdate)
	}
	
	hostToken, err := p.digitalSambaClient.CreateToken(tokenReq)
	if err != nil {
		// Clean up room if token creation fails
		_ = p.digitalSambaClient.DeleteRoom(room.ID)
		return nil, fmt.Errorf("failed to create host token: %w", err)
	}
	
	// Debug logging
	p.API.LogDebug("DigitalSamba token created", 
		"room_id", room.ID,
		"room_friendly_url", room.FriendlyURL,
		"token_room_url", hostToken.RoomURL,
		"dashboard_url", config.GetDashboardURL())

	// Construct the meeting URL
	// DigitalSamba meeting URLs follow the pattern: https://TEAM.digitalsamba.com/ROOM
	meetingURL := ""
	
	// Extract team name from the API URL
	dashboardURL := config.GetDashboardURL()
	
	// Try to extract team name from URL like https://myteam.digitalsamba.com/api/v1
	if strings.Contains(dashboardURL, ".digitalsamba.com") {
		parts := strings.Split(dashboardURL, ".")
		if len(parts) > 0 && strings.HasPrefix(dashboardURL, "https://") {
			teamName := strings.TrimPrefix(parts[0], "https://")
			meetingURL = fmt.Sprintf("https://%s.digitalsamba.com/%s", teamName, room.FriendlyURL)
		}
	}
	
	// Fallback if we couldn't extract team name
	if meetingURL == "" {
		// Use the room ID as a fallback - user will need to configure team name properly
		p.API.LogWarn("Could not extract team name from dashboard URL, using placeholder", 
			"dashboard_url", dashboardURL)
		meetingURL = fmt.Sprintf("https://CONFIGURE_TEAM_NAME.digitalsamba.com/%s", room.FriendlyURL)
	}
	
	// Create meeting post
	meetingTypeString := p.b.LocalizeWithConfig(l, &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "digitalsamba.start_meeting.meeting_id",
			Other: "Meeting ID",
		},
	})

	slackAttachment := model.SlackAttachment{
		Fallback: p.b.LocalizeWithConfig(l, &i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "digitalsamba.start_meeting.fallback_text",
				Other: `Video Meeting started at [{{.MeetingID}}]({{.MeetingURL}}).

[Join Meeting]({{.MeetingURL}})`,
			},
			TemplateData: map[string]string{
				"MeetingID":  meetingID,
				"MeetingURL": meetingURL,
			},
		}),
		Title: meetingTopic,
		Text: p.b.LocalizeWithConfig(l, &i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "digitalsamba.start_meeting.slack_attachment_text",
				Other: `{{.MeetingType}}: [{{.MeetingID}}]({{.MeetingURL}})

[Join Meeting]({{.MeetingURL}})`,
			},
			TemplateData: map[string]string{
				"MeetingType": meetingTypeString,
				"MeetingID":   meetingID,
				"MeetingURL":  meetingURL,
			},
		}),
	}

	if config.DigitalSambaRoomExpiry > 0 {
		expiryText := p.b.LocalizeWithConfig(l, &i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "digitalsamba.start_meeting.room_expires",
				Other: "Room expires at: {{.ExpiryTime}}",
			},
			TemplateData: map[string]string{
				"ExpiryTime": roomExpiry.Format("15:04 MST"),
			},
		})
		slackAttachment.Text += "\n\n" + expiryText
	}

	post := &model.Post{
		UserId:    user.Id,
		ChannelId: channel.Id,
		Type:      "custom_digitalsamba",
		Props: map[string]interface{}{
			"attachments":     []*model.SlackAttachment{&slackAttachment},
			"meeting_id":      meetingID,
			"room_id":         room.ID,
			"meeting_url":     meetingURL,
			"meeting_topic":   meetingTopic,
			"room_expires_at": roomExpiry.Unix(),
		},
		RootId: rootID,
	}

	if _, err := p.API.CreatePost(post); err != nil {
		// Clean up room if post creation fails
		_ = p.digitalSambaClient.DeleteRoom(room.ID)
		return nil, err
	}

	p.trackMeeting(nil)
	// Extract team name and room name from the URL
	// Format: https://TEAM.digitalsamba.com/ROOM
	teamName := ""
	roomName := room.FriendlyURL
	if parts := strings.Split(strings.TrimPrefix(meetingURL, "https://"), "."); len(parts) > 0 {
		teamName = parts[0]
	}
	if idx := strings.LastIndex(meetingURL, "/"); idx != -1 {
		roomName = meetingURL[idx+1:]
	}

	return &MeetingInfo{
		MeetingID: meetingID,
		RoomID:    room.ID,
		RoomURL:   meetingURL,
		Token:     "", // Don't include token - each user should fetch their own
		TeamName:  teamName,
		RoomName:  roomName,
	}, nil
}

func (p *Plugin) generateMeetingID(user *model.User, channel *model.Channel, meetingTopic string) string {
	userConfig, _ := p.getUserConfig(user.Id)
	
	switch userConfig.NamingScheme {
	case digitalSambaNameSchemeWords:
		return generateEnglishTitleName()
	case digitalSambaNameSchemeUUID:
		return generateUUIDName()
	case digitalSambaNameSchemeMattermost:
		if channel.Type == model.ChannelTypeDirect || channel.Type == model.ChannelTypeGroup {
			return generatePersonalMeetingName(user.Username)
		}
		team, _ := p.API.GetTeam(channel.TeamId)
		if team != nil {
			return generateTeamChannelName(team.Name, channel.Name)
		}
		return generateEnglishTitleName()
	default:
		if meetingTopic != "" {
			return encodeDigitalSambaMeetingID(meetingTopic)
		}
		return generateEnglishTitleName()
	}
}

func (p *Plugin) askMeetingType(user *model.User, channel *model.Channel, rootID string) error {
	l := p.b.GetUserLocalizer(user.Id)
	apiURL := *p.API.GetConfig().ServiceSettings.SiteURL + "/plugins/digitalsamba/api/v1/meetings"

	actions := []*model.PostAction{}

	var team *model.Team
	if channel.TeamId != "" {
		team, _ = p.API.GetTeam(channel.TeamId)
	}

	actions = append(actions, &model.PostAction{
		Name: p.b.LocalizeWithConfig(l, &i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "digitalsamba.ask.meeting_name_random_words",
				Other: "Meeting name with random words",
			},
		}),
		Integration: &model.PostActionIntegration{
			URL: apiURL,
			Context: map[string]interface{}{
				"meeting_id":    generateEnglishTitleName(),
				"meeting_topic": "DigitalSamba Meeting",
				"personal":      true,
			},
		},
	})

	actions = append(actions, &model.PostAction{
		Name: p.b.LocalizeWithConfig(l, &i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "digitalsamba.ask.personal_meeting",
				Other: "Personal meeting",
			},
		}),
		Integration: &model.PostActionIntegration{
			URL: apiURL,
			Context: map[string]interface{}{
				"meeting_id":    generatePersonalMeetingName(user.Username),
				"meeting_topic": fmt.Sprintf("%s's Meeting", user.GetDisplayName(model.ShowNicknameFullName)),
				"personal":      true,
			},
		},
	})

	if channel.Type == model.ChannelTypeOpen || channel.Type == model.ChannelTypePrivate {
		actions = append(actions, &model.PostAction{
			Name: p.b.LocalizeWithConfig(l, &i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "digitalsamba.ask.channel_meeting",
					Other: "Channel meeting",
				},
			}),
			Integration: &model.PostActionIntegration{
				URL: apiURL,
				Context: map[string]interface{}{
					"meeting_id":    generateTeamChannelName(team.Name, channel.Name),
					"meeting_topic": fmt.Sprintf("%s Channel Meeting", channel.DisplayName),
					"personal":      false,
				},
			},
		})
	}

	actions = append(actions, &model.PostAction{
		Name: p.b.LocalizeWithConfig(l, &i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "digitalsamba.ask.uuid_meeting",
				Other: "Meeting name with UUID",
			},
		}),
		Integration: &model.PostActionIntegration{
			URL: apiURL,
			Context: map[string]interface{}{
				"meeting_id":    generateUUIDName(),
				"meeting_topic": "DigitalSamba Meeting",
				"personal":      false,
			},
		},
	})

	sa := model.SlackAttachment{
		Title: p.b.LocalizeWithConfig(l, &i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "digitalsamba.ask.title",
				Other: "DigitalSamba Meeting Start",
			},
		}),
		Text: p.b.LocalizeWithConfig(l, &i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "digitalsamba.ask.select_meeting_type",
				Other: "Select type of meeting you want to start",
			},
		}),
		Actions: actions,
	}

	post := &model.Post{
		UserId:    p.botID,
		ChannelId: channel.Id,
		RootId:    rootID,
	}
	post.SetProps(map[string]interface{}{
		"attachments": []*model.SlackAttachment{&sa},
	})
	_ = p.API.SendEphemeralPost(user.Id, post)

	return nil
}

func (p *Plugin) trackMeeting(_ *model.CommandArgs) {
	if p.tracker == nil {
		return
	}
	
	_ = p.tracker.TrackEvent("start_meeting", map[string]interface{}{
		"plugin": "digitalsamba",
	})
}

func encodeDigitalSambaMeetingID(meeting string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9-_]+")
	meeting = strings.ReplaceAll(meeting, " ", "-")
	return reg.ReplaceAllString(meeting, "")
}

func generateUUIDName() string {
	return uuid.New().String()
}

func generatePersonalMeetingName(username string) string {
	return fmt.Sprintf("%s-personal-meeting", username)
}

func generateTeamChannelName(teamName, channelName string) string {
	return fmt.Sprintf("%s-%s-meeting", encodeDigitalSambaMeetingID(teamName), encodeDigitalSambaMeetingID(channelName))
}