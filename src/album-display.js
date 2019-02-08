/*
    Copyright (c) 2016-2019, Randy Westlund and Jacqueline Kory Westlund.
    All rights reserved.
    This code is under the BSD-2-Clause license.
*/
import '@polymer/iron-ajax/iron-ajax.js';
import '@polymer/iron-flex-layout/iron-flex-layout-classes.js';
import '@polymer/iron-image/iron-image.js';
import '@polymer/paper-styles/element-styles/paper-material-styles.js';
import '@polymer/polymer/lib/elements/dom-if.js';
import { PolymerElement, html } from '@polymer/polymer/polymer-element.js';

import { PhotosMixin } from './photos-mixin.js';
import './global-styles.js';

class AlbumDisplay extends PhotosMixin(PolymerElement) {
    static get template() {
        return html`
        <style include="iron-flex iron-flex-alignment"></style>
        <style include="paper-material-styles"></style>
        <style include="global-styles"></style>
        <style>
            :host {
                --iron-image-placeholder: {
                    background: #dddddd;
                }
            }
            iron-image {
                --iron-image-width: 100%;
                margin: 0px 5px;
            }
            span.album_name {
                font-weight: bold;
                font-size: large;
            }
            span.album_image_count {
                color: gray;
                font-style: italic;
            }
        </style>

        <iron-ajax auto="[[album.cover_image_id]]" method="GET" url="/api/photos/[[album.cover_image_id]]" handle-as="json" last-response="{{photo}}">
        </iron-ajax>

        <iron-ajax id="put_ajax" method="PUT" url="/api/photos/[[photo.id]]" body="[[item_to_edit]]" content-type="application/json" handle-as="json" last-response="{{item_to_edit}}" on-response="put_item_successful" on-error="put_item_failed" loading="{{loading}}">
        </iron-ajax>

        <a href\$="/albums/[[album.name]]" class="plain">
            <div class="paper-material card-item" elevation="1">
                <template is="dom-if" if="[[photo.id]]">
                    <iron-image src="/api/photos/[[photo.id]]/thumbnail" preload="" fade="">
                    </iron-image>
                </template>
                <div class="layout horizontal justified">
                    <span class="album_name">[[album.name]]</span>
                    <span class="album_image_count">
                        [[pretty_image_count(album.image_count)]]
                    </span>
                </div>
            </div>
        </a>
        `;
    }
    static get properties() {
        return {
            album: { type: String },
            photo: { type: Object, value: {} },
        };
    }

    // Opens the edit modal.
    edit_photo() {
        // Deep copy the object so we don't change the card's display
        // until the save is successful.
        this.set('item_to_edit', JSON.parse(JSON.stringify(this.photo)));
        // Ask for the form to be opened.
        window.dispatchEvent(new CustomEvent("open-form", {
            detail: {
                name: "edit_photo_form",
                photo: this.item_to_edit,
                callback: "resolve_edit_photo",
                that: this,
            },
        }));
    }

    // Handle response from dialog. Reason is either confirmed or canceled.
    resolve_edit_photo(e, reason) {
        if(!reason.confirmed) return;
        console.log(this.photo);
        // Override dirty checking; let Polymer know it changed.
        var tmp = this.item_to_edit;
        this.set("item_to_edit", {});
        this.set("item_to_edit", tmp);
        this.$.put_ajax.generateRequest();
    }

    // Copy response from PUT to update the display. The rationale for
    // not loading the PUT response directly into photo is to prevent
    // a failed return status from clearing the display.
    put_item_successful() {
        this.set("photo", this.item_to_edit);
        console.log("updated");
    }

    put_item_failed() {
        console.log("failed to update")
    }
}
customElements.define("album-display", AlbumDisplay);
