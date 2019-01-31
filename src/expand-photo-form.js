/*
    Copyright (c) 2016-2017, Randy Westlund and Jacqueline Kory Westlund.
    All rights reserved.
    This code is under the BSD-2-Clause license.
*/
import '@polymer/iron-flex-layout/iron-flex-layout-classes.js';
import '@polymer/iron-icon/iron-icon.js';
import '@polymer/iron-icons/hardware-icons.js';
import '@polymer/iron-image/iron-image.js';
import '@polymer/polymer/lib/elements/dom-if.js';
import { IronOverlayBehavior } from '@polymer/iron-overlay-behavior/iron-overlay-behavior.js';
import { PolymerElement, html } from '@polymer/polymer/polymer-element.js';
import { mixinBehaviors } from '@polymer/polymer/lib/legacy/class.js';

class ExpandPhotoForm extends mixinBehaviors([IronOverlayBehavior], PolymerElement) {
    static get template() {
        return html`
        <style include="iron-flex iron-flex-alignment"></style>
        <style>
            :host {
                --iron-image-placeholder: {
                    background: #dddddd;
                }
                width: 100%;
                height: 100%;
            }
            iron-image {
                --iron-image-width: 100%;
                margin: 0px 5px;
            }
            iron-icon.large-icon {
                color: magenta;
                width: 4em;
                height: 4em;
            }
            div.tall {
                height: 100%;
            }
            div.sidebar {
                cursor: pointer;
            }
        </style>

        <div class="layout horizontal justified center tall">
            <div class="layout horizontal center tall sidebar" on-tap="previous_photo">
                <iron-icon icon="hardware:keyboard-arrow-left" class="large-icon">
                </iron-icon>
            </div>
            <template is="dom-if" if="[[photoId]]">
                <iron-image src="/api/photos/[[photoId]]/big_thumbnail" preload="" fade="">
                </iron-image>
            </template>
            <div class="layout horizontal center tall sidebar" on-tap="next_photo">
                <iron-icon icon="hardware:keyboard-arrow-right" class="large-icon">
                </iron-icon>
            </div>
        </div>
        `;
    }
    static get properties() {
        return {
            photoId: { type: Number },
            photoList: { type: Array },
        };
    }

    connectedCallback() {
        super.connectedCallback();
        this._close = this.close.bind(this)
        this.addEventListener("tap", this._close);
    }
    disconnectedCallback() {
        super.disconnectedCallback();
        this.removeEventListener("tap", this._close);
    }

    previous_photo(e) {
        e.stopPropagation();
        var index = this.photoList.map(o => o.id).indexOf(this.photoId);
        index--;
        if (index < 0) index = this.photoList.length - 1;
        this.set("photoId", this.photoList[index].id);
    }

    next_photo(e) {
        e.stopPropagation();
        var index = this.photoList.map(o => o.id).indexOf(this.photoId);
        index++;
        if (index >= this.photoList.length) index = 0;
        this.set("photoId", this.photoList[index].id);
    }
}
customElements.define("expand-photo-form", ExpandPhotoForm);
