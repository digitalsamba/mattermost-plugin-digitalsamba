import {combineReducers} from 'redux';

import {RECEIVED_USER_CONFIG, OPEN_MEETING, CLOSE_MEETING} from '../action_types';
import {UserConfig, MeetingInfo} from '../types';

function userConfig(state: UserConfig | null = {embedded: true, show_prejoin_page: true, naming_scheme: 'words'}, action: any) {
    switch (action.type) {
    case RECEIVED_USER_CONFIG:
        return action.data;
    default:
        return state;
    }
}

function embeddedMeetings(state: MeetingInfo[] = [], action: any) {
    switch (action.type) {
    case OPEN_MEETING:
        console.log('[DigitalSamba Reducer] OPEN_MEETING action:', action);
        console.log('[DigitalSamba Reducer] Current meetings in state:', state);
        // Only keep one meeting at a time to avoid multiple windows
        const newState = [action.data];
        console.log('[DigitalSamba Reducer] Replaced all meetings with new meeting:', newState);
        return newState;
    case CLOSE_MEETING:
        console.log('[DigitalSamba Reducer] CLOSE_MEETING action:', action);
        const filteredState = state.filter((meeting) => meeting.meeting_id !== action.data);
        console.log('[DigitalSamba Reducer] New embedded meetings state after close:', filteredState);
        return filteredState;
    default:
        return state;
    }
}

export default combineReducers({
    userConfig,
    embeddedMeetings,
});