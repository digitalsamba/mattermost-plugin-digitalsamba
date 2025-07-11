package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost/server/public/model"
)

type StartMeetingRequest struct {
	ChannelID    string `json:"channel_id"`
	MeetingID    string `json:"meeting_id"`
	MeetingTopic string `json:"meeting_topic"`
	Personal     bool   `json:"personal"`
	RootID       string `json:"root_id"`
}

type TokenRequest struct {
	RoomID string `json:"room_id"`
}

func (p *Plugin) handleStartMeeting(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-Id")
	if userID == "" {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	var req StartMeetingRequest
	
	// Handle both POST body and PostActionIntegrationRequest
	if r.Method == http.MethodPost {
		var actionReq model.PostActionIntegrationRequest
		if err := json.NewDecoder(r.Body).Decode(&actionReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		// Extract data from context
		if context, ok := actionReq.Context["context"].(map[string]interface{}); ok {
			req.ChannelID = actionReq.ChannelId
			req.RootID = ""
			if meetingID, ok := context["meeting_id"].(string); ok {
				req.MeetingID = meetingID
			}
			if meetingTopic, ok := context["meeting_topic"].(string); ok {
				req.MeetingTopic = meetingTopic
			}
			if personal, ok := context["personal"].(bool); ok {
				req.Personal = personal
			}
		}
		
		// Delete the ephemeral post
		if actionReq.PostId != "" {
			p.deleteEphemeralPost(userID, actionReq.PostId)
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	user, appErr := p.API.GetUser(userID)
	if appErr != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	channel, appErr := p.API.GetChannel(req.ChannelID)
	if appErr != nil {
		http.Error(w, "Failed to get channel", http.StatusInternalServerError)
		return
	}

	meetingInfo, err := p.startMeeting(user, channel, req.MeetingID, req.MeetingTopic, req.Personal, req.RootID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to start meeting: %v", err), http.StatusInternalServerError)
		return
	}

	resp := meetingInfo
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (p *Plugin) handleConfig(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-Id")
	if userID == "" {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	config, err := p.getUserConfig(userID)
	if err != nil {
		http.Error(w, "Failed to get user config", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

func (p *Plugin) handleUserConfig(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-Id")
	if userID == "" {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodGet {
		config, err := p.getUserConfig(userID)
		if err != nil {
			http.Error(w, "Failed to get user config", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(config)
		return
	}

	if r.Method == http.MethodPost {
		var config UserConfig
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := p.setUserConfig(userID, &config); err != nil {
			http.Error(w, "Failed to save user config", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (p *Plugin) handleGetToken(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-Id")
	if userID == "" {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	var req TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, appErr := p.API.GetUser(userID)
	if appErr != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Create moderator token for the user
	tokenReq := &CreateTokenRequest{
		RoomID:    req.RoomID,
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
	
	token, err := p.digitalSambaClient.CreateToken(tokenReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create token: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token.Token,
	})
}

func (p *Plugin) deleteEphemeralPost(userID, postID string) {
	p.API.DeleteEphemeralPost(userID, postID)
}