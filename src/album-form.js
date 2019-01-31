/*
    Copyright (c) 2016-2017, Randy Westlund and Jacqueline Kory Westlund.
    All rights reserved.
    This code is under the BSD-2-Clause license.
*/
import '@polymer/paper-input/paper-input.js';
import '@rwestlund/responsive-dialog/responsive-dialog.js';
import { PolymerElement, html } from '@polymer/polymer/polymer-element.js';

import { FormMixin } from './form-mixin.js';

class AlbumForm extends FormMixin(PolymerElement) {
    static get template() {
        return html`
        <responsive-dialog id="dialog" title="[[title]]" dismiss-text="Cancel" confirm-text="Save" on-iron-overlay-closed="resolve_dialog">
            <paper-input type="text" label="Album" value="{{album.name}}" autocapitalize="words" char-counter="" maxlength="20">
            </paper-input>
        </responsive-dialog>
        `;
    }
    static get properties() {
        return {
            album: { type: Object },
            title: { type: String, value: "Create album" },
        };
    }
}
customElements.define("album-form", AlbumForm);
