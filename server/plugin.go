package main

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/mattermost/mattermost/server/public/pluginapi/experimental/telemetry"
	"github.com/mattermost/mattermost/server/public/pluginapi/i18n"
	"github.com/pkg/errors"
)

const digitalSambaNameSchemeAsk = "ask"
const digitalSambaNameSchemeWords = "words"
const digitalSambaNameSchemeUUID = "uuid"
const digitalSambaNameSchemeMattermost = "mattermost"
const configChangeEvent = "config_update"

type UserConfig struct {
	NamingScheme    string `json:"naming_scheme"`
	Embedded        bool   `json:"embedded"`
	ShowPrejoinPage bool   `json:"show_prejoin_page"`
}

type Plugin struct {
	plugin.MattermostPlugin

	client *pluginapi.Client

	telemetryClient telemetry.Client
	tracker         telemetry.Tracker

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration.
	configuration *configuration

	b *i18n.Bundle

	botID string

	// DigitalSamba API client
	digitalSambaClient *DigitalSambaClient
}

func (p *Plugin) OnActivate() error {
	config := p.getConfiguration()
	if err := config.IsValid(); err != nil {
		return err
	}

	command, err := p.createDigitalSambaCommand()
	if err != nil {
		return err
	}

	if err = p.API.RegisterCommand(command); err != nil {
		return err
	}

	i18nBundle, err := i18n.InitBundle(p.API, filepath.Join("assets", "i18n"))
	if err != nil {
		return err
	}
	p.b = i18nBundle

	digitalSambaBot := &model.Bot{
		Username:    "digitalsamba",
		DisplayName: "DigitalSamba",
		Description: "A bot account created by the DigitalSamba plugin",
	}
	options := []pluginapi.EnsureBotOption{
		pluginapi.ProfileImagePath("assets/icon.png"),
	}

	p.client = pluginapi.NewClient(p.API, p.Driver)
	botID, ensureBotError := p.client.Bot.EnsureBot(digitalSambaBot, options...)
	if ensureBotError != nil {
		return errors.Wrap(ensureBotError, "failed to ensure DigitalSamba bot user")
	}

	p.botID = botID

	// Initialize DigitalSamba client
	p.digitalSambaClient = NewDigitalSambaClient(config.GetDashboardURL(), config.DigitalSambaAPIKey)

	p.telemetryClient, err = telemetry.NewRudderClient()
	if err != nil {
		p.API.LogWarn("telemetry client not started", "error", err.Error())
	}

	return nil
}

func (p *Plugin) OnDeactivate() error {
	if p.telemetryClient != nil {
		_ = p.telemetryClient.Close()
	}
	return nil
}

func (p *Plugin) getConfiguration() *configuration {
	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()

	if p.configuration == nil {
		return &configuration{}
	}

	return p.configuration
}

func (p *Plugin) setConfiguration(configuration *configuration) {
	p.configurationLock.Lock()
	defer p.configurationLock.Unlock()

	if configuration != nil && p.configuration == configuration {
		panic("setConfiguration called with the existing configuration")
	}

	p.configuration = configuration
}

func (p *Plugin) OnConfigurationChange() error {
	var configuration = new(configuration)

	if err := p.API.LoadPluginConfiguration(configuration); err != nil {
		return errors.Wrap(err, "failed to load plugin configuration")
	}

	if err := configuration.IsValid(); err != nil {
		return errors.Wrap(err, "configuration is invalid")
	}

	p.setConfiguration(configuration)

	// Update DigitalSamba client with new configuration
	if p.digitalSambaClient != nil {
		p.digitalSambaClient = NewDigitalSambaClient(configuration.GetDashboardURL(), configuration.DigitalSambaAPIKey)
	}

	return nil
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/v1/meetings":
		p.handleStartMeeting(w, r)
	case "/api/v1/token":
		p.handleGetToken(w, r)
	case "/api/v1/config":
		p.handleConfig(w, r)
	case "/api/v1/user-config":
		p.handleUserConfig(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (p *Plugin) getUserConfig(userID string) (*UserConfig, error) {
	data, appErr := p.API.KVGet("config_" + userID)
	if appErr != nil {
		return nil, appErr
	}

	if data == nil {
		return &UserConfig{
			Embedded:        p.getConfiguration().DigitalSambaEmbedded,
			NamingScheme:    p.getConfiguration().DigitalSambaNamingScheme,
			ShowPrejoinPage: p.getConfiguration().DigitalSambaShowPrejoinPage,
		}, nil
	}

	var userConfig UserConfig
	err := json.Unmarshal(data, &userConfig)
	if err != nil {
		return nil, err
	}

	return &userConfig, nil
}

func (p *Plugin) setUserConfig(userID string, config *UserConfig) error {
	b, err := json.Marshal(config)
	if err != nil {
		return err
	}

	appErr := p.API.KVSet("config_"+userID, b)
	if appErr != nil {
		return appErr
	}

	p.API.PublishWebSocketEvent(configChangeEvent, nil, &model.WebsocketBroadcast{UserId: userID})
	return nil
}