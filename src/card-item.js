/*
    Copyright (c) 2019, Randy Westlund and Jacqueline Kory Westlund.
    All rights reserved.
    This code is under the BSD-2-Clause license.
*/

import '@polymer/iron-flex-layout/iron-flex-layout-classes.js';
import '@polymer/paper-icon-button/paper-icon-button.js';
import '@polymer/paper-menu-button/paper-menu-button.js';
import '@polymer/paper-styles/element-styles/paper-material-styles.js';
import '@polymer/paper-styles/typography.js';
import '@material/mwc-icon/mwc-icon.js';
import { PolymerElement, html } from '@polymer/polymer/polymer-element.js';

class CardItem extends PolymerElement {
    static get template() {
        return html`
        <style include="iron-flex iron-flex-alignment"></style>
        <style include="paper-material-styles"></style>
        <style>
            * { box-sizing: border-box; }
            [hidden] { display: none; }
            :host {
                margin-bottom: 5px;
                border-radius: 3px;
                background-color: white;
                @apply --paper-material;
                @apply --paper-material-elevation-1;
            }
            /* Disabled cards are grayed out. */
            :host([disabled]) {
                background-color: var(--paper-grey-300);
                color: var(--paper-grey-500);
            }
            span.title {
                @apply --paper-font-title;
            }
            div.header {
                padding-top: 7px;
                padding-left: 7px;
                overflow: hidden;
            }
            span.subtitle {
                @apply --paper-font-caption;
                margin-left: 1em;
            }
            paper-menu-button {
                padding: 0;
                /* This is necessary if [disabled]. */
                color: black;
            }
            div.content {
                padding-left: 7px;
                padding-right: 7px;
                padding-bottom: 7px;
            }
            mwc-icon {
                color: gray;
            }
        </style>
        <div class="layout horizontal justified start">
            <div class="header layout horizontal wrap end">
                <slot name="prefix"></slot>
                <mwc-icon hidden="[[!icon]]">[[icon]]</mwc-icon>
                <span class="title">[[title]]</span>
                <span class="subtitle">[[subtitle]]</span>
            </div>
            <paper-menu-button hidden="[[hideMenu]]" horizontal-align="right">
                <paper-icon-button slot="dropdown-trigger" icon="icons:more-vert">
                </paper-icon-button>
                <slot name="dropdown-content" slot="dropdown-content"></slot>
            </paper-menu-button>
        </div>
        <div class="content">
            <slot></slot>
        </div>
        `;
    }
    static get properties() {
        return {
            title: { type: String },
            subtitle: { type: String },
            hideMenu: { type: Boolean, value: false },
            icon: { type: String, value: "" },
        };
    }
}
customElements.define("card-item", CardItem);
