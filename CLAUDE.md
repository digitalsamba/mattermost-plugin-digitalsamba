# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Mattermost plugin that integrates DigitalSamba video conferencing into Mattermost. It's based on the architecture of the mattermost-plugin-jitsi but uses DigitalSamba's API and SDK instead.

## Development Commands

### Build Commands
```bash
# Build the plugin for all platforms
make dist

# Build server component only
make server

# Build webapp component only  
make webapp

# Run linters and type checks
make check-style

# Run tests
make test

# Deploy to local Mattermost instance
make deploy

# Watch mode for development
make watch
```

### Development Setup
```bash
# Install Go dependencies
go mod download

# Install webapp dependencies
cd webapp && npm install

# Install linting tools
make install-go-tools
```

## Architecture Overview

### Plugin Structure
- **Server (Go)**: Handles API endpoints, slash commands, bot interactions, and DigitalSamba API integration
- **Webapp (React/TypeScript)**: Provides UI components for meeting controls and embedded video
- **Hybrid Architecture**: Server manages authentication and room creation; webapp handles video embedding

### Key Components

#### Server-side (`/server`)
- `plugin.go`: Main plugin logic, handles activation and configuration
- `api.go`: REST API endpoints for meeting management
- `command.go`: Slash command implementation (`/digitalsamba`)
- `configuration.go`: Plugin settings management
- `meeting.go`: DigitalSamba room creation and token generation

#### Client-side (`/webapp`)
- `index.tsx`: Plugin initialization and registration
- `components/conference/`: Embedded video conference component using DigitalSamba SDK
- `components/post_type_digitalsamba/`: Custom post type for meeting announcements
- `actions/`: Redux actions for meeting management
- `client/`: API client for server communication

### DigitalSamba Integration Points

1. **Authentication**: API key stored in plugin configuration
2. **Room Creation**: Server creates rooms via DigitalSamba REST API
3. **Token Generation**: Server generates join tokens for participants
4. **Video Embedding**: Webapp uses @digitalsamba/embedded-sdk for video UI
5. **Meeting Links**: Can be opened externally or embedded within Mattermost

### Key Differences from Jitsi Plugin

1. **API Authentication**: DigitalSamba uses Bearer token auth instead of JWT
2. **SDK Integration**: Uses npm package instead of external script loading
3. **Room Management**: Requires explicit room creation via API
4. **Participant Limits**: Supports up to 2000 participants
5. **Feature Set**: Additional features like polling, Q&A, breakout rooms

## Configuration Schema

The plugin configuration includes:
- `DigitalSambaAPIKey`: API authentication key
- `DigitalSambaDashboardURL`: URL for DigitalSamba dashboard
- `DigitalSambaEmbedded`: Enable embedded video mode
- `DigitalSambaShowPrejoinPage`: Show pre-join configuration
- `DigitalSambaNamingScheme`: Meeting ID generation method
- `DigitalSambaDefaultRoomSettings`: JSON configuration for room defaults

## Testing Strategy

- Unit tests for server-side meeting logic
- Component tests for React components
- Integration tests for DigitalSamba API calls
- Mock DigitalSamba SDK for webapp tests

## Security Considerations

- API keys must be stored securely in plugin settings
- Room tokens should have appropriate expiration times
- Validate all user inputs before API calls
- Use HTTPS for all DigitalSamba communications

## Common Development Tasks

### Adding a New Meeting Feature
1. Update server API to support the feature
2. Modify room creation logic if needed
3. Update webapp components to expose the feature
4. Add appropriate i18n translations
5. Update tests

### Debugging Video Issues
1. Check browser console for SDK errors
2. Verify API key is valid and has appropriate permissions
3. Check room creation response for errors
4. Ensure participant tokens are generated correctly

### Updating DigitalSamba SDK
1. Update package.json with new version
2. Run `npm install` in webapp directory
3. Update any changed API calls
4. Test embedded and external meeting modes