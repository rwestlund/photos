/*
    Copyright (c) 2016-2017, Randy Westlund and Jacqueline Kory Westlund.
    All rights reserved.
    This code is under the BSD-2-Clause license.
*/
import '@material/mwc-icon/mwc-icon.js';
import '@polymer/paper-button/paper-button.js';
import '@polymer/polymer/lib/elements/dom-if.js';
import '@polymer/polymer/lib/elements/dom-repeat.js';
import { PolymerElement, html } from '@polymer/polymer/polymer-element.js';
import '@vaadin/vaadin-upload/vaadin-upload.js';

import './global-styles.js';

class PhotosUploads extends PolymerElement {
    static get template() {
        return html`
        <style include="global-styles"></style>
        <style>
            :host {
                display: block;

                --vaadin-upload-drop-label: {
                    color: var(--secondary-text-color);
                }
            }
            vaadin-upload {
                --primary-color: var(--header-background-color);
                --light-primary-color: var(--medium-background-color);
                background-color: #ffffff;
                margin-top: 1em;
            }
            vaadin-upload[nodrop] {
                background-color: inherit;
            }
        </style>

        <paper-button raised="" on-tap="select_albums">
            <mwc-icon>photo_album</mwc-icon>
            Set albums for uploads
        </paper-button>

        <template is="dom-if" if="[[albums_selected.length]]">
            <p><strong>Selected albums:</strong>
                </p><ul>
                    <template is="dom-repeat" items="[[albums_selected]]">
                        <li>[[item]]</li>
                    </template>
                </ul>
            <p></p>
        </template>

        <vaadin-upload
                target="/api/photos"
                accept="image/*,video/*"
                on-upload-request="handle_upload_request">
        </vaadin-upload>
        `;
    }

    static get properties() {
        return {
            albums_selected: { type: Array, value: [] },
        };
    }

    select_albums() {
        // Deep copy the object so we don't change the card's
        // display until the save is successful.
        this.set("edit_albums_selected",
            JSON.parse(JSON.stringify(this.albums_selected)));
        // Ask for the form to be opened.
        window.dispatchEvent(new CustomEvent("open-form", {
            detail: {
                name:               "select_album_form",
                selectedAlbums:     this.edit_albums_selected,
                that:               this,
                callback:           "resolve_album_select_dialog",
            },
        }));
    }

    handle_upload_request(e, data) {
        data.formData.append("albums", JSON.stringify(this.albums_selected));
    }

    // Handle response from dialog. Reason is either confirmed or canceled.
    resolve_album_select_dialog(e, reason) {
        if (!reason.confirmed) return;
        this.set("albums_selected", this.edit_albums_selected);
    }
}
customElements.define("photos-uploads", PhotosUploads);
