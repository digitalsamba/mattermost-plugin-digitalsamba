package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type DigitalSambaClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

type Room struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	FriendlyURL       string    `json:"friendly_url"`
	Privacy           string    `json:"privacy"`
	MaxParticipants   int       `json:"max_participants"`
	SessionDuration   int       `json:"session_duration"`
	EnableRecording   bool      `json:"enable_recording"`
	EnableBreakoutRooms bool    `json:"enable_breakout_rooms"`
	ExpiresAt         *time.Time `json:"expires_at,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

type CreateRoomRequest struct {
	Topic             string     `json:"topic,omitempty"`
	Description       string     `json:"description,omitempty"`
	FriendlyURL       string     `json:"friendly_url,omitempty"`
	Privacy           string     `json:"privacy,omitempty"`
	MaxParticipants   int        `json:"max_participants,omitempty"`
	RecordingsEnabled bool       `json:"recordings_enabled,omitempty"`
	ChatEnabled       bool       `json:"chat_enabled,omitempty"`
	ExpiresAt         *time.Time `json:"expires_at,omitempty"`
	
	// Join Settings
	JoinScreenEnabled *bool      `json:"join_screen_enabled,omitempty"`
	MuteOnJoin        *bool      `json:"mute_on_join,omitempty"`
	CameraOffOnJoin   *bool      `json:"camera_off_on_join,omitempty"`
	
	// Consent Settings - disabled for internal use
	ConsentMessage    string     `json:"consent_message,omitempty"`
	
	// Layout Settings
	DefaultLayout     string     `json:"default_layout,omitempty"` // "auto" or "tiled"
	
	// Feature Settings
	EnableWhiteboard  bool       `json:"enable_whiteboard,omitempty"`
	EnablePolling     bool       `json:"enable_polling,omitempty"`
	EnableQA          bool       `json:"enable_qa,omitempty"`
}

type RoomToken struct {
	Token       string    `json:"token"`
	RoomURL     string    `json:"room_url"`
	ParticipantID string  `json:"participant_id"`
	Role        string    `json:"role"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

type CreateTokenRequest struct {
	RoomID    string `json:"room_id"`
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email,omitempty"`
	Role      string `json:"role"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

func NewDigitalSambaClient(baseURL, apiKey string) *DigitalSambaClient {
	return &DigitalSambaClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *DigitalSambaClient) doRequest(method, path string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	// Ensure the path starts with /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	fullURL := c.baseURL + path
	req, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API error: status=%d, url=%s, body=%s", resp.StatusCode, fullURL, string(body))
	}

	return resp, nil
}

func (c *DigitalSambaClient) CreateRoom(req *CreateRoomRequest) (*Room, error) {
	resp, err := c.doRequest("POST", "/rooms", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var room Room
	if err := json.NewDecoder(resp.Body).Decode(&room); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &room, nil
}

func (c *DigitalSambaClient) GetRoom(roomID string) (*Room, error) {
	resp, err := c.doRequest("GET", "/rooms/"+roomID, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var room Room
	if err := json.NewDecoder(resp.Body).Decode(&room); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &room, nil
}

func (c *DigitalSambaClient) DeleteRoom(roomID string) error {
	resp, err := c.doRequest("DELETE", "/rooms/"+roomID, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *DigitalSambaClient) CreateToken(req *CreateTokenRequest) (*RoomToken, error) {
	// The endpoint is /rooms/{room}/token
	endpoint := fmt.Sprintf("/rooms/%s/token", req.RoomID)
	
	// Create a new request body without the room_id field
	tokenBody := map[string]interface{}{
		"role": req.Role,
		"name": req.UserName,
		"u":    req.UserName, // Username parameter to bypass login dialog
	}
	
	// Add external user ID for cross-referencing
	if req.UserID != "" {
		tokenBody["ud"] = req.UserID // External user identifier (Mattermost user ID)
	}
	
	if req.UserEmail != "" {
		tokenBody["email"] = req.UserEmail
	}
	
	if req.AvatarURL != "" {
		tokenBody["avatar_url"] = req.AvatarURL
		tokenBody["avatar"] = req.AvatarURL // Also set 'avatar' as per documentation
	}
	
	resp, err := c.doRequest("POST", endpoint, tokenBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var token RoomToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &token, nil
}