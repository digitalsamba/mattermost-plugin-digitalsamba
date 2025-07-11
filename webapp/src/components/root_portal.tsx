import React from 'react';
import ReactDOM from 'react-dom';
import {Provider} from 'react-redux';

import Conference from './conference';

export default class RootPortal {
    private registry: any;
    private store: any;
    private portalNode: HTMLElement | null = null;

    constructor(registry: any, store: any) {
        this.registry = registry;
        this.store = store;
    }

    render() {
        console.log('[DigitalSamba RootPortal] Rendering...');
        if (!this.portalNode) {
            const rootPortalId = 'digitalsamba-root-portal';
            this.portalNode = document.getElementById(rootPortalId);
            
            if (!this.portalNode) {
                this.portalNode = document.createElement('div');
                this.portalNode.id = rootPortalId;
                document.body.appendChild(this.portalNode);
                console.log('[DigitalSamba RootPortal] Created portal node');
            }
        }

        ReactDOM.render(
            <Provider store={this.store}>
                <Conference/>
            </Provider>,
            this.portalNode
        );
        console.log('[DigitalSamba RootPortal] Conference component rendered');
    }

    cleanup() {
        if (this.portalNode) {
            ReactDOM.unmountComponentAtNode(this.portalNode);
            if (this.portalNode.parentNode) {
                this.portalNode.parentNode.removeChild(this.portalNode);
            }
            this.portalNode = null;
        }
    }
}