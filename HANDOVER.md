# DigitalSamba Mattermost Plugin - Development Handover

## Current Status (Version 1.0.6)

### What's Working
1. **API Integration**: Server successfully creates rooms and tokens via DigitalSamba API
2. **Token Generation**: Each user gets their own token with identity (name, email, avatar) and moderator role
3. **Meeting Creation**: Meetings can be created from channel header or slash commands
4. **URL Construction**: Meeting URLs are properly constructed from dashboard URL

### Current Issues
1. **Embedded Mode**: Meeting opens but shows "Join meeting" screen (token not being recognized)
2. **Multiple Windows**: Sometimes shows duplicate meeting windows
3. **SDK Initialization**: Errors with DigitalSamba SDK initialization

### Key Files to Review

#### Server-side
- `/server/digitalsamba_client.go` - API client for DigitalSamba
- `/server/meeting.go` - Meeting creation logic, returns MeetingInfo struct
- `/server/api.go` - REST endpoints including token generation
- `/server/configuration.go` - Plugin configuration

#### Client-side
- `/webapp/src/components/conference/conference.tsx` - Embedded meeting component (MAIN ISSUE HERE)
- `/webapp/src/components/post_type_digitalsamba/` - Meeting post UI
- `/webapp/src/actions/index.ts` - Redux actions
- `/webapp/src/reducers/index.ts` - State management
- `/webapp/src/client/client.ts` - API client

### Critical Code Areas

#### Current SDK Initialization (conference.tsx:40-55)
```typescript
const initOptions = {
    url: `${activeMeeting.room_url}?token=${encodeURIComponent(activeMeeting.token)}`,
    frame: {
        container: containerRef.current,
    },
};
sambaRef.current = new DigitalSambaEmbedded(initOptions);
```

#### Token Flow
1. User clicks "Join Meeting"
2. Client fetches token via `/api/v1/token` endpoint
3. Token includes user identity and moderator role
4. Token should be passed to SDK but isn't working

### Debug Information Needed
1. Check browser console for DigitalSamba SDK errors
2. Inspect iframe `src` attribute - does it include token parameter?
3. Network tab - check requests to digitalsamba.com
4. Check if token JWT is valid (decode at jwt.io)

### DigitalSamba SDK Notes
- Version: 0.0.48 (from package.json)
- Documentation: https://docs.digitalsamba.com/reference/sdk/digitalsambaembedded-class
- Expected URL format: `https://TEAM.digitalsamba.com/ROOM?token=TOKEN`

### Next Steps
1. Debug why token in URL isn't authenticating the user
2. Consider trying alternative SDK initialization methods
3. Check if there's a CSP issue preventing iframe from loading properly
4. Verify token is valid and contains correct claims

### Environment
- Plugin Version: 1.0.6
- Mattermost instance: https://mattermost.digitalsamba.com
- API endpoint: https://api.digitalsamba.com/api/v1

### Console Error Logs
[TO BE ADDED BY USER]