import React from 'react';
import {Post} from 'mattermost-redux/types/posts';
import {useDispatch, useSelector} from 'react-redux';
import {GlobalState} from 'mattermost-redux/types/store';

import {openMeeting} from '../../actions';
import Client from '../../client';

interface Props {
    post: Post;
}

export default function PostTypeDigitalSamba(props: Props) {
    const dispatch = useDispatch();
    const userConfig = useSelector((state: GlobalState) => (state as any)['plugins-digitalsamba']?.userConfig);
    const pluginState = useSelector((state: GlobalState) => (state as any)['plugins-digitalsamba']);
    
    console.log('[DigitalSamba PostType] Component rendered with userConfig:', userConfig);
    console.log('[DigitalSamba PostType] Full plugin state:', pluginState);
    
    const meetingId = props.post.props?.meeting_id;
    const meetingUrl = props.post.props?.meeting_url;
    const roomId = props.post.props?.room_id;
    const meetingTopic = props.post.props?.meeting_topic || 'DigitalSamba Meeting';
    
    const handleJoinMeeting = async () => {
        console.log('[DigitalSamba] Join meeting clicked', {
            userConfig,
            meetingId,
            meetingUrl,
            roomId,
            embedded: userConfig?.embedded
        });
        if (userConfig?.embedded && roomId && meetingUrl) {
            // Extract team name and room name from URL
            let teamName = '';
            let roomName = meetingId;
            
            try {
                console.log('[DigitalSamba] Parsing meeting URL:', meetingUrl);
                const url = new URL(meetingUrl);
                const hostParts = url.hostname.split('.');
                if (hostParts.length > 2) {
                    teamName = hostParts[0];
                }
                const pathParts = url.pathname.split('/').filter(p => p);
                if (pathParts.length > 0) {
                    roomName = pathParts[pathParts.length - 1];
                }
            } catch (e) {
                console.error('Failed to parse meeting URL:', e);
            }
            
            try {
                // Fetch token for this room
                console.log('[DigitalSamba] Fetching token for room:', roomId);
                const token = await Client.getToken(roomId);
                console.log('[DigitalSamba] Token received:', token ? 'yes' : 'no');
                
                // Create meeting info object
                const meetingInfo = {
                    meeting_id: meetingId,
                    room_id: roomId,
                    room_url: meetingUrl,
                    token: token,
                    team_name: teamName,
                    room_name: roomName,
                };
                
                console.log('[DigitalSamba] Dispatching openMeeting with:', meetingInfo);
                dispatch(openMeeting(meetingInfo));
            } catch (error) {
                console.error('[DigitalSamba] Failed to get token:', error);
                // Fallback to external link
                window.open(meetingUrl, '_blank');
            }
        } else if (meetingUrl && roomId) {
            // For external meetings, fetch token and append to URL
            console.log('[DigitalSamba] Opening meeting in new tab, fetching token first');
            try {
                const token = await Client.getToken(roomId);
                const urlWithToken = `${meetingUrl}?token=${encodeURIComponent(token)}`;
                console.log('[DigitalSamba] Opening meeting with token URL');
                window.open(urlWithToken, '_blank');
            } catch (error) {
                console.error('[DigitalSamba] Failed to get token, opening without token:', error);
                window.open(meetingUrl, '_blank');
            }
        }
    };

    if (!meetingId) {
        return null;
    }

    return (
        <div className='digitalsamba-post-type'>
            <div className='digitalsamba-post-header'>
                <h4>{meetingTopic}</h4>
                <p>Meeting ID: {meetingId}</p>
            </div>
            <button
                className='btn btn-primary'
                onClick={handleJoinMeeting}
            >
                Join Meeting
            </button>
        </div>
    );
}