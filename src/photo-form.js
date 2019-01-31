/*
    Copyright (c) 2016-2017, Randy Westlund and Jacqueline Kory Westlund.
    All rights reserved.
    This code is under the BSD-2-Clause license.
*/
import '@polymer/iron-ajax/iron-ajax.js';
import '@polymer/iron-flex-layout/iron-flex-layout-classes.js';
import '@polymer/iron-icons/iron-icons.js';
import '@polymer/paper-dropdown-menu/paper-dropdown-menu.js';
import '@polymer/paper-icon-button/paper-icon-button.js';
import '@polymer/paper-input/paper-input.js';
import '@polymer/paper-item/paper-item.js';
import '@polymer/paper-listbox/paper-listbox.js';
import '@polymer/paper-styles/element-styles/paper-material-styles.js';
import '@polymer/polymer/lib/elements/dom-repeat.js';
import { PolymerElement, html } from '@polymer/polymer/polymer-element.js';
import '@rwestlund/responsive-dialog/responsive-dialog.js';

import { FormMixin } from './form-mixin.js';

class PhotoForm extends FormMixin(PolymerElement) {
    static get template() {
        return html`
        <style include="iron-flex iron-flex-alignment"></style>
        <style include="paper-material-styles"></style>
        <style>
            .photo-item-label {
                padding-top: 1.2em;
                display: block;
            }
            .album-item {
                margin-right: 1em;
                padding-right: 1em;
            }
        </style>

        <iron-ajax id="get_albums_ajax" method="GET" url="/api/albums" handle-as="json" last-response="{{albums}}">
        </iron-ajax>

        <responsive-dialog id="dialog" title="Edit photo" dismiss-text="Cancel" confirm-text="Save" on-iron-overlay-closed="resolve_dialog">
            <paper-input type="text" label="Caption" value="{{photo.caption}}" autocapitalize="sentences" char-counter="" maxlength="140">
            </paper-input>

            <strong class="photo-item-label">Albums</strong>
                <div class="layout horizontal wrap">
                    <template is="dom-repeat" items="[[photo.albums]]">
                        <div class="paper-material album-item" elevation="1">
                            <paper-icon-button icon="icons:cancel" data-index\$="[[index]]" on-tap="remove_album">
                            </paper-icon-button>
                            <span>[[item]]</span>
                        </div>
                    </template>
                </div>

            <div class="layout vertical">
                <paper-dropdown-menu label="Select Albums" vertical-align="top" horizontal-align="right">
                    <paper-listbox slot="dropdown-content" id="album_dd" attr-for-selected="value" on-iron-select="add_selected_album">
                        <template is="dom-repeat" items="[[albums]]">
                            <paper-item value="[[item.name]]">
                                [[item.name]]
                            </paper-item>
                        </template>
                    </paper-listbox>
                </paper-dropdown-menu>
            </div>
        </responsive-dialog>
        `;
    }
    static get properties() {
        return {
            photo: { type: Object },
        };
    }
    open_hook() {
        this.$.get_albums_ajax.generateRequest();
        this.$.album_dd.selected = null;
    }
    // Add an album to the album list, checking for duplicates first.
    add_selected_album(e, item) {
        for (var i in this.photo.albums)
            if (this.photo.albums[i] === this.$.album_dd.selected)
                return;
        this.push('photo.albums', this.$.album_dd.selected);
    }
    remove_album(e) {
        var index = Number(e.currentTarget.attributes['data-index'].value);
        this.splice('photo.albums', index, 1);
    }
}
customElements.define("photo-form", PhotoForm);
