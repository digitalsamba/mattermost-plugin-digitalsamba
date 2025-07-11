# DigitalSamba Plugin Debugging Notes

## Recent Changes Made

### 1. URL Construction (server/meeting.go)
- Fixed duplicate `/api/v1` in API URLs
- Extracts team name from dashboard URL to construct meeting URLs
- Format: `https://TEAM.digitalsamba.com/ROOM`

### 2. Token Management
- Server no longer includes creator's token in meeting response
- Each user fetches their own token when joining
- All users join as moderators
- Tokens include full Mattermost identity

### 3. Embedded Meeting Component
- Tried multiple SDK initialization approaches:
  - Constructor with nested options
  - Factory method with createControl()
  - Direct URL with token parameter
- Current approach: Pass full URL with token to SDK

### 4. Redux State
- Changed to only allow one meeting at a time
- Prevents multiple meeting windows

## Known Issues

### SDK Initialization Error
```
TypeError: Cannot set properties of undefined (setting 'display')
```
This suggests the SDK is trying to manipulate DOM elements that don't exist.

### Token Not Working
Despite passing token in URL, users still see "Join meeting" screen asking for name.

## Things to Try

1. **Check Token Format**
   - Decode the JWT token to verify claims
   - Ensure it has required fields for DigitalSamba

2. **Alternative SDK Approaches**
   ```typescript
   // Option 1: Use token property directly
   new DigitalSambaEmbedded({
       team: teamName,
       room: roomName,
       token: token,
       frame: { container: element }
   });

   // Option 2: Use roomSettings
   new DigitalSambaEmbedded({
       url: roomUrl,
       token: token,
       roomSettings: { username: userName }
   });
   ```

3. **Check Network Requests**
   - Is the iframe making authenticated requests?
   - Are cookies being set?
   - Any CORS issues?

4. **DOM Timing**
   - Ensure container exists before SDK init
   - Try setTimeout to delay initialization

## Logging Added
- All major flow points have console.log statements
- Look for `[DigitalSamba]` prefix in console
- Token is partially logged for security

## Build Commands
```bash
make dist  # Builds plugin
# Output: dist/digitalsamba-1.0.6.tar.gz
```