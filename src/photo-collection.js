/*
    Copyright (c) 2016-2019, Randy Westlund and Jacqueline Kory Westlund.
    All rights reserved.
    This code is under the BSD-2-Clause license.
*/
import '@polymer/app-layout/app-grid/app-grid-style.js';
import '@polymer/iron-ajax/iron-ajax.js';
import '@polymer/iron-flex-layout/iron-flex-layout-classes.js';
import '@polymer/iron-media-query/iron-media-query.js';
import '@polymer/paper-icon-button/paper-icon-button.js';
import '@polymer/paper-spinner/paper-spinner.js';
import '@polymer/polymer/lib/elements/dom-if.js';
import '@polymer/polymer/lib/elements/dom-repeat.js';
import { PolymerElement, html } from '@polymer/polymer/polymer-element.js';

import './cookie-display.js';
import { PhotosMixin } from './photos-mixin.js';
import './global-styles.js';
import './photo-display.js';

class PhotoCollection extends PhotosMixin(PolymerElement) {
    static get template() {
        return html`
        <style include="app-grid-style"></style>
        <style include="iron-flex iron-flex-alignment"></style>
        <style include="global-styles"></style>
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
            span.album_image_count {
                color: gray;
                font-size: large;
                font-style: italic;
            }
        </style>

        <cookie-display cookie-name="role" cookie-value="{{user_role}}">
        </cookie-display>

        <iron-ajax
                auto=""
                method="GET"
                url="/api/photos"
                params="[[search_filter]]"
                handle-as="json"
                last-response="{{photos}}">
        </iron-ajax>

        <iron-ajax id="put_ajax"
                method="PUT"
                url="/api/albums/[[albumName]]"
                body="[[item_to_edit]]"
                content-type="application/json"
                handle-as="json"
                last-response="{{item_to_edit}}"
                on-response="put_item_successful"
                on-error="put_item_failed"
                loading="{{loading}}">
        </iron-ajax>

         <template is="dom-if" if="[[albumName]]">
            <iron-ajax
                    auto="[[albumName]]"
                    method="GET"
                    url="/api/albums/[[albumName]]"
                    handle-as="json"
                    last-response="{{album}}">
            </iron-ajax>


            <div class="layout horizontal justified">
                <div>
                    <span class="page_name">[[albumName]]</span>
                    <template is="dom-if" if="[[user_is_admin(user_role)]]">
                        <paper-icon-button icon="icons:create" on-tap="edit_album">
                        </paper-icon-button>
                        <paper-spinner active="[[loading]]"></paper-spinner>
                    </template>
                </div>
                <span class="album_image_count">
                    [[pretty_image_count(album.image_count)]]
                </span>
            </div>
        </template>

        <template is="dom-if" if="[[!photos.length]]">
            <h4>This album has no photos!</h4>
        </template>

        <div class="app-grid">
            <template is="dom-if" if="[[albumName]]">
                <template is="dom-repeat" items="[[photos]]">
                    <photo-display photo="[[item]]" album="" on-set-album-image="set_album_image" on-expand-photo="show_big_thumbnail">
                    </photo-display>
                </template>
            </template>
            <template is="dom-if" if="[[!albumName]]">
                <template is="dom-repeat" items="[[photos]]">
                    <photo-display photo="[[item]]" on-expand-photo="show_big_thumbnail">
                    </photo-display>
                </template>
            </template>
        </div>
        `;
    }
    static get properties() {
        return {
            albumName: { type: String, value: "" },
            search_filter: {
                type: Object,
                computed: "compute_search_filter(albumName)",
            }
        };
    }

    compute_search_filter(album) {
        var o = {};
        if (album) o.album = album;
        return o;
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

    edit_album() {
        this.set('item_to_edit', JSON.parse(JSON.stringify(this.album)));
        // Ask for the form to be opened.
        window.dispatchEvent(new CustomEvent("open-form", {
            detail: {
                name: "edit_album_form",
                album: this.item_to_edit,
                callback: "resolve_edit_album",
                that: this,
            },
        }));
    }

    resolve_edit_album(e, reason) {
        if(!reason.confirmed) return;
        // Override dirty checking; let Polymer know it changed.
        var tmp = this.item_to_edit;
        this.set("item_to_edit", {});
        this.set("item_to_edit", tmp);
        this.$.put_ajax.generateRequest();
    }

    put_item_successful() {
        console.log("updated");
        if (this.albumName !== this.item_to_edit.name) {
            window.history.pushState({}, null, "/albums/" +
                this.item_to_edit.name);
            window.dispatchEvent(new CustomEvent("location-changed"));
        }
    }

    put_item_failed() { console.log("failed to update"); }

    set_album_image(e) {
        this.set('item_to_edit', JSON.parse(JSON.stringify(this.album)));
        this.item_to_edit.cover_image_id = e.detail;
        this.$.put_ajax.generateRequest();
    }

    // Opens the big thumbnail display.
    show_big_thumbnail(e) {
        // Ask for the big thumbnail display to be opened.
        window.dispatchEvent(new CustomEvent("open-form", {
            detail: {
                name: "expand_photo_form",
                photoId: e.detail,
                photoList: this.photos,
            },
        }));
    }
}
customElements.define("photo-collection", PhotoCollection);
