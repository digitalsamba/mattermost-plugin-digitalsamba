import {Dispatch} from 'redux';
import {GetStateFunc} from 'mattermost-redux/types/actions';

import Client from '../client';
import {UserConfig} from '../types';
import {RECEIVED_USER_CONFIG, OPEN_MEETING, CLOSE_MEETING} from '../action_types';

export function startMeeting(channelId: string, topic = '', rootId = '') {
    return async (dispatch: Dispatch, getState: GetStateFunc) => {
        try {
            const result = await Client.startMeeting(channelId, topic, rootId);
            return {data: result};
        } catch (error) {
            return {error};
        }
    };
}

export function loadConfig() {
    return async (dispatch: Dispatch) => {
        try {
            const config = await Client.getUserConfig();
            dispatch({
                type: RECEIVED_USER_CONFIG,
                data: config,
            });
            return {data: config};
        } catch (error) {
            return {error};
        }
    };
}

export function updateConfig(config: UserConfig) {
    return async (dispatch: Dispatch) => {
        try {
            await Client.updateUserConfig(config);
            dispatch({
                type: RECEIVED_USER_CONFIG,
                data: config,
            });
            return {data: config};
        } catch (error) {
            return {error};
        }
    };
}

export function openMeeting(meetingInfo: any) {
    return {
        type: OPEN_MEETING,
        data: meetingInfo,
    };
}

export function closeMeeting(meetingId: string) {
    return {
        type: CLOSE_MEETING,
        data: meetingId,
    };
}