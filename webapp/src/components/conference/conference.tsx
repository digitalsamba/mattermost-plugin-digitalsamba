import React, {useEffect, useRef, useState} from 'react';
import {useSelector, useDispatch} from 'react-redux';
import {GlobalState} from 'mattermost-redux/types/store';
import DigitalSambaEmbedded from '@digitalsamba/embedded-sdk';

import {closeMeeting} from '../../actions';
import {MeetingInfo} from '../../types';
import './conference.scss';

export default function Conference() {
    console.log('[DigitalSamba Conference] Component rendered');
    const dispatch = useDispatch();
    const containerRef = useRef<HTMLDivElement>(null);
    const sambaRef = useRef<DigitalSambaEmbedded | null>(null);
    const [isMinimized, setIsMinimized] = useState(false);
    
    const embeddedMeetings: MeetingInfo[] = useSelector((state: GlobalState) => 
        (state as any)['plugins-digitalsamba']?.embeddedMeetings || []
    );
    const userConfig = useSelector((state: GlobalState) => 
        (state as any)['plugins-digitalsamba']?.userConfig
    );

    const activeMeeting = embeddedMeetings[embeddedMeetings.length - 1];
    console.log('[DigitalSamba Conference] Component state - Active meeting:', activeMeeting);
    console.log('[DigitalSamba Conference] Component state - All meetings:', embeddedMeetings);
    console.log('[DigitalSamba Conference] Component state - Container ref:', containerRef.current);
    console.log('[DigitalSamba Conference] Component state - Samba ref:', sambaRef.current);

    useEffect(() => {
        console.log('[DigitalSamba Conference] useEffect triggered', {
            activeMeeting: activeMeeting,
            containerRef: containerRef.current,
            sambaRef: sambaRef.current
        });
        
        const initializeSdk = async () => {
            if (activeMeeting && containerRef.current && !sambaRef.current) {
                console.log('[DigitalSamba Conference] Initializing SDK...');
                try {
                    // Initialize DigitalSamba embedded SDK with all options
                    const initOptions = {
                        url: activeMeeting.room_url,
                        token: activeMeeting.token,
                        root: containerRef.current,
                    };
                    console.log('[DigitalSamba Conference] Creating DigitalSamba with token as separate parameter');
                    console.log('[DigitalSamba Conference] Init options:', {
                        ...initOptions,
                        token: 'TOKEN_HIDDEN', // Hide token in logs
                    });
                    
                    sambaRef.current = DigitalSambaEmbedded.createControl(initOptions);
                    await sambaRef.current.load();
                
                console.log('[DigitalSamba Conference] SDK instance created successfully');
                
                // Listen for events only if SDK initialized successfully
                if (sambaRef.current) {
                    sambaRef.current.on('meetingEnded', () => {
                        console.log('[DigitalSamba Conference] Meeting ended event');
                        handleClose();
                    });

                    sambaRef.current.on('leftMeeting', () => {
                        console.log('[DigitalSamba Conference] Left meeting event');
                        handleClose();
                    });
                }
                } catch (error) {
                    console.error('[DigitalSamba Conference] Failed to initialize SDK:', error);
                    // Clean up on error
                    sambaRef.current = null;
                }
            }
        };

        initializeSdk();

        return () => {
            if (sambaRef.current) {
                sambaRef.current.destroy();
                sambaRef.current = null;
            }
        };
    }, [activeMeeting]);

    const handleClose = () => {
        if (activeMeeting) {
            dispatch(closeMeeting(activeMeeting.meeting_id));
        }
        if (sambaRef.current) {
            sambaRef.current.destroy();
            sambaRef.current = null;
        }
    };

    const handleMinimize = () => {
        setIsMinimized(!isMinimized);
    };

    if (!activeMeeting) {
        console.log('[DigitalSamba Conference] No active meeting, returning null');
        return null;
    }

    return (
        <div className={`digitalsamba-conference ${isMinimized ? 'minimized' : ''}`}>
            <div className='digitalsamba-conference-header'>
                <span>DigitalSamba Meeting - {activeMeeting.room_name}</span>
                <div className='digitalsamba-conference-controls'>
                    <button
                        onClick={handleMinimize}
                        className='digitalsamba-minimize-btn'
                        title={isMinimized ? 'Maximize' : 'Minimize'}
                    >
                        {isMinimized ? '□' : '_'}
                    </button>
                    <button
                        onClick={handleClose}
                        className='digitalsamba-close-btn'
                        title='Close'
                    >
                        ×
                    </button>
                </div>
            </div>
            {!isMinimized && (
                <div
                    ref={containerRef}
                    className='digitalsamba-conference-container'
                />
            )}
        </div>
    );
}