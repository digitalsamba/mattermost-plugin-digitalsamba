# Mattermost DigitalSamba Plugin

This plugin integrates DigitalSamba video conferencing into Mattermost, allowing users to start and join video meetings directly from their Mattermost channels.

## Features

- Start video meetings with `/digitalsamba` slash command
- Embedded video meetings within Mattermost (optional)
- Multiple meeting naming schemes (random words, UUID, context-based, or user choice)
- Support for up to 2000 participants per meeting
- Meeting recording capabilities (configurable)
- Breakout rooms support (configurable)
- Automatic room expiration

## Requirements

- Mattermost Server v5.2.0 or higher
- DigitalSamba API key (get one at https://dashboard.digitalsamba.com)

## Installation

1. Download the latest plugin file from the releases page
2. In Mattermost, go to **System Console > Plugins > Plugin Management**
3. Upload the plugin file
4. Enable the plugin
5. Configure the plugin settings (see Configuration section)

## Configuration

After installation, configure the plugin in **System Console > Plugins > DigitalSamba**:

### Required Settings

- **DigitalSamba API Key**: Your DigitalSamba API key
- **DigitalSamba Dashboard URL**: API endpoint URL (default: https://api.digitalsamba.com/v1)

### Optional Settings

- **Embed Video Inside Mattermost**: When enabled, meetings open in a floating window
- **Show Pre-join Page**: Display settings page before joining embedded meetings
- **Meeting Names**: Choose how meeting IDs are generated
- **Room Expiry Time**: Minutes before unused rooms expire (0 = no expiry)
- **Maximum Participants**: Max participants per room (1-2000)
- **Enable Recording**: Allow meeting hosts to record
- **Enable Breakout Rooms**: Allow breakout room creation

## Usage

### Starting a Meeting

- `/digitalsamba` - Start a meeting with a random name
- `/digitalsamba [topic]` - Start a meeting with a specific topic

### Managing Settings

- `/digitalsamba settings` - View your personal settings
- `/digitalsamba settings naming_scheme [words|uuid|mattermost|ask]` - Set naming scheme
- `/digitalsamba settings embed [true|false]` - Toggle embedded meetings

### Meeting Features

- Click the video icon in the channel header to start a meeting
- Join meetings by clicking "Join Meeting" in meeting posts
- Embedded meetings appear as a floating window (if enabled)
- External meetings open in a new browser tab

## Development

### Prerequisites

- Go 1.16+
- Node.js 14+
- npm 6+

### Building the Plugin

```bash
# Clone the repository
git clone https://github.com/mattermost-community/mattermost-plugin-digitalsamba.git
cd mattermost-plugin-digitalsamba

# Install dependencies
make deps

# Build the plugin
make dist
```

### Running Tests

```bash
# Run server tests
make test

# Run webapp tests
cd webapp && npm test
```

### Watching for Changes

```bash
make watch
```

## Troubleshooting

### Meeting Creation Fails

1. Verify your API key is correct
2. Check the API endpoint URL
3. Ensure your DigitalSamba account has sufficient quota

### Embedded Meetings Not Working

1. Check browser console for errors
2. Ensure pop-ups are allowed for your Mattermost domain
3. Try disabling embedded mode and using external meetings

### Performance Issues

1. Reduce maximum participants if experiencing lag
2. Disable recording if not needed
3. Use external meetings instead of embedded for large meetings

## Contributing

Contributions are welcome! Please read our [Contributing Guidelines](CONTRIBUTING.md) before submitting PRs.

## License

This plugin is licensed under the [Apache 2.0 License](LICENSE).