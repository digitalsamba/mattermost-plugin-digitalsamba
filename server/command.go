package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/pluginapi/i18n"
)

const commandHelp = `* |/digitalsamba| - Start a meeting with a random name
* |/digitalsamba [topic]| - Start a meeting with specified topic
* |/digitalsamba settings| - View your current settings
* |/digitalsamba settings [setting] [value]| - Update your settings
  * |setting| can be "naming_scheme" or "embed"
  * |naming_scheme| values: "words", "uuid", "mattermost", "ask"
  * |embed| values: "true", "false"
* |/digitalsamba help| - Show this help text`

func (p *Plugin) createDigitalSambaCommand() (*model.Command, error) {
	iconData := ""

	return &model.Command{
		Trigger:              "digitalsamba",
		DisplayName:          "DigitalSamba",
		Description:          "Start and manage DigitalSamba meetings",
		AutoComplete:         true,
		AutoCompleteDesc:     "Available commands: start, settings, help",
		AutoCompleteHint:     "[command]",
		AutocompleteData:     getAutocompleteData(),
		AutocompleteIconData: iconData,
	}, nil
}

func getAutocompleteData() *model.AutocompleteData {
	command := model.NewAutocompleteData("digitalsamba", "[command]", "Available commands: start, settings, help")

	start := model.NewAutocompleteData("start", "[topic]", "Start a meeting")
	start.AddTextArgument("Topic of the meeting", "[topic]", "")
	command.AddCommand(start)

	settings := model.NewAutocompleteData("settings", "[setting] [value]", "Update your personal settings")
	settings.AddStaticListArgument("setting", true, []model.AutocompleteListItem{
		{Item: "naming_scheme", HelpText: "Set the naming scheme for meetings"},
		{Item: "embed", HelpText: "Set whether to embed meetings"},
	})
	command.AddCommand(settings)

	help := model.NewAutocompleteData("help", "", "Display usage information")
	command.AddCommand(help)

	return command
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	p.trackCommand(args)

	fields := strings.Fields(args.Command)
	if len(fields) == 0 || fields[0] != "/digitalsamba" {
		return &model.CommandResponse{}, nil
	}

	if len(fields) == 1 {
		// Just "/digitalsamba" - start a meeting
		return p.runStartMeetingCommand(args, "")
	}

	subcommand := fields[1]

	switch subcommand {
	case "help":
		return p.runHelpCommand(args)
	case "settings":
		if len(fields) == 2 {
			return p.runShowSettingsCommand(args)
		}
		if len(fields) >= 4 {
			return p.runUpdateSettingsCommand(args, fields[2], strings.Join(fields[3:], " "))
		}
		return p.sendEphemeralResponse(args, "Invalid settings command. Use `/digitalsamba settings` to view or `/digitalsamba settings [setting] [value]` to update.")
	case "start":
		topic := ""
		if len(fields) > 2 {
			topic = strings.Join(fields[2:], " ")
		}
		return p.runStartMeetingCommand(args, topic)
	default:
		// Treat everything else as a meeting topic
		topic := strings.Join(fields[1:], " ")
		return p.runStartMeetingCommand(args, topic)
	}
}

func (p *Plugin) runHelpCommand(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	return p.sendEphemeralResponse(args, commandHelp)
}

func (p *Plugin) runShowSettingsCommand(args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	userConfig, err := p.getUserConfig(args.UserId)
	if err != nil {
		return p.sendEphemeralResponse(args, "Failed to get user settings")
	}

	l := p.b.GetUserLocalizer(args.UserId)
	message := p.b.LocalizeWithConfig(l, &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "digitalsamba.command.settings.current",
			Other: `Current DigitalSamba Settings:
* Naming Scheme: {{.NamingScheme}}
* Embed Video: {{.Embed}}
* Show Pre-join Page: {{.ShowPrejoin}}`,
		},
		TemplateData: map[string]string{
			"NamingScheme": userConfig.NamingScheme,
			"Embed":        fmt.Sprintf("%v", userConfig.Embedded),
			"ShowPrejoin":  fmt.Sprintf("%v", userConfig.ShowPrejoinPage),
		},
	})

	return p.sendEphemeralResponse(args, message)
}

func (p *Plugin) runUpdateSettingsCommand(args *model.CommandArgs, setting, value string) (*model.CommandResponse, *model.AppError) {
	userConfig, err := p.getUserConfig(args.UserId)
	if err != nil {
		return p.sendEphemeralResponse(args, "Failed to get user settings")
	}

	switch setting {
	case "naming_scheme":
		validSchemes := []string{"words", "uuid", "mattermost", "ask"}
		valid := false
		for _, scheme := range validSchemes {
			if value == scheme {
				valid = true
				userConfig.NamingScheme = value
				break
			}
		}
		if !valid {
			return p.sendEphemeralResponse(args, fmt.Sprintf("Invalid naming scheme. Valid values are: %s", strings.Join(validSchemes, ", ")))
		}
	case "embed":
		if value == "true" {
			userConfig.Embedded = true
		} else if value == "false" {
			userConfig.Embedded = false
		} else {
			return p.sendEphemeralResponse(args, "Invalid embed value. Use 'true' or 'false'")
		}
	default:
		return p.sendEphemeralResponse(args, "Invalid setting. Valid settings are: naming_scheme, embed")
	}

	if err := p.setUserConfig(args.UserId, userConfig); err != nil {
		return p.sendEphemeralResponse(args, "Failed to update settings")
	}

	return p.sendEphemeralResponse(args, "Settings updated successfully")
}

func (p *Plugin) runStartMeetingCommand(args *model.CommandArgs, topic string) (*model.CommandResponse, *model.AppError) {
	user, appErr := p.API.GetUser(args.UserId)
	if appErr != nil {
		return p.sendEphemeralResponse(args, "Failed to get user information")
	}

	channel, appErr := p.API.GetChannel(args.ChannelId)
	if appErr != nil {
		return p.sendEphemeralResponse(args, "Failed to get channel information")
	}

	if topic == "" {
		userConfig, err := p.getUserConfig(args.UserId)
		if err != nil {
			return p.sendEphemeralResponse(args, "Failed to get user settings")
		}

		if userConfig.NamingScheme == digitalSambaNameSchemeAsk {
			if err := p.askMeetingType(user, channel, args.RootId); err != nil {
				return p.sendEphemeralResponse(args, "Failed to display meeting options")
			}
			return &model.CommandResponse{}, nil
		}
	}

	meetingID, err := p.startMeeting(user, channel, "", topic, false, args.RootId)
	if err != nil {
		return p.sendEphemeralResponse(args, fmt.Sprintf("Failed to start meeting: %v", err))
	}

	return &model.CommandResponse{
		Text: fmt.Sprintf("Meeting started: %s", meetingID),
	}, nil
}

func (p *Plugin) sendEphemeralResponse(args *model.CommandArgs, message string) (*model.CommandResponse, *model.AppError) {
	post := &model.Post{
		UserId:    p.botID,
		ChannelId: args.ChannelId,
		Message:   message,
	}
	_ = p.API.SendEphemeralPost(args.UserId, post)
	return &model.CommandResponse{}, nil
}

func (p *Plugin) trackCommand(args *model.CommandArgs) {
	if p.tracker == nil {
		return
	}

	fields := strings.Fields(args.Command)
	if len(fields) < 2 {
		return
	}

	var event string
	switch fields[1] {
	case "help":
		event = "help_command"
	case "settings":
		event = "settings_command"
	default:
		event = "start_meeting_command"
	}

	_ = p.tracker.TrackEvent(event, map[string]interface{}{})
}