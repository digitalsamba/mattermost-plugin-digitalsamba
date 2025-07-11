import {Client4} from 'mattermost-redux/client';
import {UserConfig} from '../types';

class Client {
    private serverRoute = '';

    setServerRoute(route: string) {
        this.serverRoute = route + '/plugins/digitalsamba';
    }

    startMeeting = async (channelId: string, topic = '', rootId = '') => {
        const url = `${this.serverRoute}/api/v1/meetings`;
        const body = {
            channel_id: channelId,
            meeting_topic: topic,
            root_id: rootId,
        };

        const response = await fetch(url, Client4.getOptions({
            method: 'POST',
            body: JSON.stringify(body),
        }));

        if (!response.ok) {
            throw new Error('Failed to start meeting');
        }

        return response.json();
    };

    getUserConfig = async (): Promise<UserConfig> => {
        const url = `${this.serverRoute}/api/v1/user-config`;
        console.log('[DigitalSamba Client] Getting user config from:', url);
        
        const response = await fetch(url, Client4.getOptions({
            method: 'GET',
        }));

        if (!response.ok) {
            console.error('[DigitalSamba Client] Failed to get user config:', response.status);
            throw new Error('Failed to get user config');
        }

        const config = await response.json();
        console.log('[DigitalSamba Client] User config received:', config);
        return config;
    };

    updateUserConfig = async (config: UserConfig) => {
        const url = `${this.serverRoute}/api/v1/user-config`;
        
        const response = await fetch(url, Client4.getOptions({
            method: 'POST',
            body: JSON.stringify(config),
        }));

        if (!response.ok) {
            throw new Error('Failed to update user config');
        }
    };

    getToken = async (roomId: string): Promise<string> => {
        const url = `${this.serverRoute}/api/v1/token`;
        console.log('[DigitalSamba Client] Getting token for room:', roomId, 'URL:', url);
        
        const response = await fetch(url, Client4.getOptions({
            method: 'POST',
            body: JSON.stringify({ room_id: roomId }),
        }));

        console.log('[DigitalSamba Client] Token response status:', response.status);
        
        if (!response.ok) {
            const errorText = await response.text();
            console.error('[DigitalSamba Client] Token error response:', errorText);
            throw new Error(`Failed to get token: ${response.status} ${errorText}`);
        }

        const data = await response.json();
        console.log('[DigitalSamba Client] Token received');
        return data.token;
    };
}

const client = new Client();
export default client;