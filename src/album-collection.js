/*
    Copyright (c) 2016-2017, Randy Westlund and Jacqueline Kory Westlund.
    All rights reserved.
    This code is under the BSD-2-Clause license.
*/
import '@polymer/app-layout/app-grid/app-grid-style.js';
import '@polymer/iron-ajax/iron-ajax.js';
import '@polymer/polymer/lib/elements/dom-repeat.js';
import { PolymerElement, html } from '@polymer/polymer/polymer-element.js';

import './album-display.js';

class AlbumCollection extends PolymerElement {
    static get template() {
        return html`
        <style include="app-grid-style"></style>
        <style>
            :host {
                display: block;
                    --app-grid-columns: 2;
                    --app-grid-gutter: 15px;
            }
            @media (max-width: 800px) {
                :host {
                    --app-grid-columns: 1;
                }
            }
        </style>

        <iron-ajax auto="" method="GET" url="/api/albums" handle-as="json" last-response="{{albums}}">
        </iron-ajax>

        <div class="app-grid">
            <template is="dom-repeat" items="[[albums]]">
                <album-display album="[[item]]"></album-display>
            </template>
        </div>
        `;
    }

    // Reevaluate app-grid styles on resize.
    connectedCallback() {
        super.connectedCallback();
        this._resize_listener = () => this.updateStyles();
        window.addEventListener("resize", this._resize_listener);
    }
    disconnectedCallback() {
        super.disconnectedCallback();
        window.removeEventListener("resize", this._resize_listener);
    }
}
customElements.define("album-collection", AlbumCollection);
