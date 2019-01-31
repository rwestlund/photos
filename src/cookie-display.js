/* Copyright (c) 2016, Randy Westlund. All rights reserved. */
/* This module shows the current value of a cookie. */
import { PolymerElement, html } from '@polymer/polymer/polymer-element.js';

class CookieDisplay extends PolymerElement {
    static get template() {
        return html`
        <style>
            :host {
                display: none;
            }
        </style>
        `;
    }

    static get properties() {
        return {
            cookieName: String,
            cookieValue: {
                type: String,
                computed: "read_cookie(cookieName)",
                notify: true
            },
        };
    }

    read_cookie(name) {
        var cstring = decodeURIComponent(document.cookie);
        var parts = cstring.split(';');
        var search_str = name + '=';
        for (var i = 0; i < parts.length; i++) {
            while (parts[i].charAt(0) === ' ')
                parts[i] = parts[i].substring(1);
            if (parts[i].indexOf(search_str) !== -1)
                return parts[i].substring(search_str.length);
        }
        return '';
    }
}
customElements.define("cookie-display", CookieDisplay);
