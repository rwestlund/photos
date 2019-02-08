/*
    Copyright (c) 2016-2017, Randy Westlund and Jacqueline Kory Westlund.
    All rights reserved.
    This code is under the BSD-2-Clause license.
*/
import '@polymer/iron-ajax/iron-ajax.js';
import '@polymer/iron-image/iron-image.js';
import '@polymer/paper-icon-button/paper-icon-button.js';
import '@polymer/paper-spinner/paper-spinner.js';
import '@polymer/paper-styles/element-styles/paper-material-styles.js';
import '@polymer/polymer/lib/elements/dom-if.js';
import '@polymer/polymer/lib/elements/dom-repeat.js';
import { PolymerElement, html } from '@polymer/polymer/polymer-element.js';

import './cookie-display.js';
import { PhotosMixin } from './photos-mixin.js';
import './global-styles.js';
class PhotoDisplay extends PhotosMixin(PolymerElement) {
    static get template() {
        return html`
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
            p.album_list {
                color: gray;
                font-style: italic;
                word-spacing: 2px;
                margin: 0;
            }
        </style>

        <iron-ajax id="put_ajax" method="PUT" url="/api/photos/[[photo.id]]" body="[[item_to_edit]]" content-type="application/json" handle-as="json" last-response="{{item_to_edit}}" on-response="put_item_successful" on-error="put_item_failed" loading="{{loading}}">
        </iron-ajax>

        <cookie-display cookie-name="role" cookie-value="{{user_role}}">
        </cookie-display>

        <div class="paper-material card-item" elevation="1">
            <template is="dom-if" if="[[user_is_admin(user_role)]]">
                <paper-icon-button icon="icons:create" on-tap="edit_photo">
                </paper-icon-button>
                <template is="dom-if" if="[[album]]">
                    <paper-icon-button icon="image:photo-album" on-tap="set_album_image">
                    </paper-icon-button>
                </template>
                <paper-spinner active="[[loading]]"></paper-spinner>
            </template>
            <iron-image src="/api/photos/[[photo.id]]/thumbnail" preload="" fade="" on-tap="display_big_photo">
            </iron-image>
            <p>[[photo.caption]]</p>
            <p class="album_list">
                Albums:
                <template is="dom-repeat" items="[[photo.albums]]">
                    <a href="/albums/[[item]]" class="plain">[[item]]</a>
                </template>
            </p>
        </div>
        `;
    }
    static get properties() {
        return {
            photo: { type: Object, value: () => ({}) },
            album: { type: Boolean, value: false },
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
    put_item_failed() { console.log("failed to update"); }

    set_album_image() {
        this.dispatchEvent(new CustomEvent("set-album-image", {
            bubbles: true,
            composed: true,
            detail: this.photo.id,
        }));
    }

    // Opens the big thumbnail display.
    display_big_photo() {
        this.dispatchEvent(new CustomEvent("expand-photo", {
            bubbles: true,
            composed: true,
            detail: this.photo.id,
        }));
    }
}
customElements.define("photo-display", PhotoDisplay);
