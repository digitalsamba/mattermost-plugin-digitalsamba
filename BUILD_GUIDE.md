# DigitalSamba Mattermost Plugin Build Guide

## Prerequisites

Before building the plugin, ensure you have:

- **Go** 1.16 or higher
- **Node.js** 14 or higher  
- **npm** 6 or higher
- **Make** build tool
- **Git** (optional, for version tracking)

## Project Structure

```
mattermost-digitalsamba-plugin/
├── server/              # Go backend code
│   ├── plugin.go        # Main plugin logic
│   ├── api.go          # HTTP endpoints
│   ├── command.go      # Slash command handler
│   ├── configuration.go # Plugin settings
│   ├── digitalsamba_client.go  # DigitalSamba API client
│   └── meeting.go      # Meeting management logic
├── webapp/             # React frontend code
│   ├── src/
│   │   ├── components/ # React components
│   │   ├── actions/    # Redux actions
│   │   ├── reducers/   # Redux reducers
│   │   └── client/     # API client
│   └── package.json    # npm dependencies
├── assets/             # Plugin assets
├── build/              # Build scripts
├── plugin.json         # Plugin manifest
├── Makefile           # Build configuration
└── go.mod             # Go dependencies
```

## Build Process

### 1. Install Dependencies

First, install the webapp dependencies:

```bash
cd webapp
npm install --legacy-peer-deps
cd ..
```

The `--legacy-peer-deps` flag is needed due to some peer dependency conflicts in the Mattermost ecosystem.

### 2. Update Go Dependencies

Ensure Go dependencies are up to date:

```bash
go mod tidy
```

### 3. Build Commands

#### Full Build (Recommended)
```bash
make dist
```

This command:
- Builds server binaries for all platforms (Linux, macOS, Windows)
- Builds the webapp JavaScript bundle
- Packages everything into a `.tar.gz` file

#### Individual Build Steps

Build only the server:
```bash
make server
```

Build only the webapp:
```bash
make webapp
```

Create the plugin bundle (after building server and webapp):
```bash
make bundle
```

### 4. Using Docker (Alternative)

If you have Docker installed, you can use the included `docker-make` script:

```bash
./docker-make dist
```

This builds the plugin in a Docker container with all dependencies pre-installed.

## Version Management

Before building a new release:

1. Update the version in `plugin.json`:
```json
{
    "version": "1.0.6",
    ...
}
```

2. The build process will automatically:
   - Include the version in the filename
   - Update the manifest in the bundled plugin

## Build Output

The built plugin will be located at:
```
dist/digitalsamba-{version}.tar.gz
```

For example: `dist/digitalsamba-1.0.5.tar.gz`

## Common Build Issues

### 1. Git Repository Warnings
You may see warnings about "not a git repository". These are harmless and occur because the build scripts try to include git information in the build.

### 2. npm Dependency Conflicts
If you encounter npm errors, try:
```bash
cd webapp
rm -rf node_modules package-lock.json
npm install --legacy-peer-deps
```

### 3. Missing Build Tools
If `make` commands fail, ensure the build tools are compiled:
```bash
cd build/manifest && go build -o ../bin/manifest
cd ../pluginctl && go build -o ../bin/pluginctl
```

## Development Workflow

1. Make your code changes
2. Test locally if possible
3. Update version in `plugin.json`
4. Run `make dist`
5. Upload the generated `.tar.gz` file to Mattermost

## Quick Build Commands Reference

```bash
# Full build
make dist

# Clean and rebuild
make clean && make dist

# Development build with file watching
make watch

# Run tests
make test

# Check code style
make check-style

# Deploy to local Mattermost (requires configuration)
make deploy
```

## Integration with DigitalSamba

This plugin integrates with DigitalSamba's API. Key integration points:

- **API Authentication**: Uses Bearer token with API key
- **Room Creation**: POST to `/api/v1/rooms`
- **Token Generation**: POST to `/api/v1/rooms/{room}/token`
- **Embedded SDK**: Uses `@digitalsamba/embedded-sdk` npm package

Ensure your DigitalSamba API key is configured in the plugin settings after installation.