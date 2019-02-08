/*
    Copyright (c) 2016-2017, Randy Westlund and Jacqueline Kory Westlund.
    All rights reserved.
    This code is under the BSD-2-Clause license.
*/
/* This module displays a user editing form. */
import '@polymer/paper-dropdown-menu/paper-dropdown-menu.js';
import '@polymer/paper-icon-button/paper-icon-button.js';
import '@polymer/paper-input/paper-input.js';
import '@polymer/paper-item/paper-item.js';
import '@polymer/paper-listbox/paper-listbox.js';
import '@polymer/polymer/lib/elements/dom-repeat.js';
import { PolymerElement, html } from '@polymer/polymer/polymer-element.js';
import '@rwestlund/responsive-dialog/responsive-dialog.js';

import { FormMixin } from './form-mixin.js';
import { PhotosMixin } from './photos-mixin.js';
import './global-styles.js';

class UserForm extends FormMixin(PhotosMixin(PolymerElement)) {
    static get template() {
        return html`
        <style include="iron-flex"></style>
        <style include="global-styles"></style>
        <responsive-dialog id="dialog" title="[[title]]" dismiss-text="Cancel" confirm-text="Save" on-iron-overlay-closed="resolve_dialog">

            <div class="layout vertical">
                <paper-dropdown-menu label="Account Role" vertical-align="top" horizontal-align="right">
                    <paper-listbox slot="dropdown-content" attr-for-selected="value" selected="{{user.role}}">
                        <template is="dom-repeat" items="[[constants.user_roles]]">
                            <paper-item value="[[item]]">[[item]]</paper-item>
                        </template>
                    </paper-listbox>
                </paper-dropdown-menu>
            </div>

            <paper-input type="email" label="Email" value="{{user.email}}" char-counter="" maxlength="50">
                <paper-icon-button slot="suffix" icon="icons:clear" on-tap="clear_field">
                </paper-icon-button>
            </paper-input>

        </responsive-dialog>
        `;
    }
    static get properties() {
        return {
            // Must be provided by parent object.
            user: { type: Object },
            title: { type: String, value: "Edit User" },
        };
    }
}
customElements.define("user-form", UserForm);
