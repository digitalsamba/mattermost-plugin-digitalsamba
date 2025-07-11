import * as React from 'react';

import {Channel} from 'mattermost-redux/types/channels';
import {Post} from 'mattermost-redux/types/posts';
import {getConfig} from 'mattermost-redux/selectors/entities/general';
import {GlobalState} from 'mattermost-redux/types/store';

import Icon from './components/icon';
import PostTypeDigitalSamba from './components/post_type_digitalsamba';
import I18nProvider from './components/i18n_provider';
import RootPortal from './components/root_portal';
import reducer from './reducers';
import {startMeeting, loadConfig, openMeeting} from './actions';
import manifest from './manifest';
import Client from './client';

class PluginClass {
    rootPortal?: RootPortal;

    initialize(registry: any, store: any) {
        this.rootPortal = new RootPortal(registry, store);
        if (this.rootPortal) {
            this.rootPortal.render();
        }

        registry.registerReducer(reducer);

        const action = async (channel: Channel) => {
            const result = await store.dispatch(startMeeting(channel.id));
            if (result.data && store.getState()['plugins-digitalsamba']?.userConfig?.embedded) {
                // Fetch token for current user
                try {
                    const token = await Client.getToken(result.data.room_id);
                    const meetingInfo = {
                        ...result.data,
                        token: token,
                    };
                    store.dispatch(openMeeting(meetingInfo));
                } catch (error) {
                    console.error('[DigitalSamba] Failed to get token:', error);
                    // Fallback to external link
                    if (result.data.room_url) {
                        window.open(result.data.room_url, '_blank');
                    }
                }
            } else if (result.data && result.data.room_url && result.data.room_id) {
                // Open in new tab for non-embedded mode with token
                try {
                    const token = await Client.getToken(result.data.room_id);
                    const urlWithToken = `${result.data.room_url}?token=${encodeURIComponent(token)}`;
                    window.open(urlWithToken, '_blank');
                } catch (error) {
                    console.error('[DigitalSamba] Failed to get token:', error);
                    window.open(result.data.room_url, '_blank');
                }
            }
        };
        const helpText = 'Start DigitalSamba Meeting';

        // Channel header icon
        registry.registerChannelHeaderButtonAction(<Icon/>, action, helpText);

        // App Bar icon
        if (registry.registerAppBarComponent) {
            const config = getConfig(store.getState());
            const siteUrl = (config && config.SiteURL) || '';
            const iconURL = `${siteUrl}/plugins/${manifest.id}/public/app-bar-icon.png`;
            registry.registerAppBarComponent(iconURL, action, helpText);
        }

        Client.setServerRoute(getServerRoute(store.getState()));
        registry.registerPostTypeComponent('custom_digitalsamba', (props: { post: Post }) => (
            <I18nProvider><PostTypeDigitalSamba post={props.post}/></I18nProvider>));
        registry.registerWebSocketEventHandler('custom_digitalsamba_config_update', () => {
            console.log('[DigitalSamba] WebSocket config update received');
            store.dispatch(loadConfig());
        });
        console.log('[DigitalSamba] Plugin initialized, loading config...');
        store.dispatch(loadConfig());
    }

    uninitialize() {
        if (this.rootPortal) {
            this.rootPortal.cleanup();
        }
    }
}

(global as any).window.registerPlugin('digitalsamba', new PluginClass());

function getServerRoute(state: GlobalState) {
    const config = getConfig(state);
    let basePath = '';
    if (config && config.SiteURL) {
        basePath = config.SiteURL;
        if (basePath && basePath[basePath.length - 1] === '/') {
            basePath = basePath.substr(0, basePath.length - 1);
        }
    }
    return basePath;
}