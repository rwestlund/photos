/*
    Copyright (c) 2016-2017, Randy Westlund and Jacqueline Kory Westlund.
    All rights reserved.
    This code is under the BSD-2-Clause license.
*/
/* This module shows a user card. */
import '@polymer/iron-ajax/iron-ajax.js';
import '@polymer/paper-button/paper-button.js';
import '@polymer/paper-dialog/paper-dialog.js';
import '@material/mwc-icon/mwc-icon.js';
import { PolymerElement, html } from '@polymer/polymer/polymer-element.js';
import { GestureEventListeners } from '@polymer/polymer/lib/mixins/gesture-event-listeners.js';

import { PhotosMixin } from './photos-mixin.js';
import './card-item.js';
import './global-styles.js';

class UserDisplay extends PhotosMixin(GestureEventListeners(PolymerElement)) {
    static get template() {
        return html`
        <style include="global-styles"></style>
        <style>
            :host {
                display: block;
            }
        </style>

        <iron-ajax id="put_ajax"
                method="PUT"
                url="/api/users/[[user.id]]"
                body="[[item_to_edit]]"
                content-type="application/json"
                handle-as="json"
                last-response="{{item_to_edit}}"
                on-response="put_item_successful"
                on-error="put_item_failed"
                loading="{{loading.put_item}}">
        </iron-ajax>
        <iron-ajax id="delete_item_ajax"
                method="DELETE"
                url="/api/users/[[user.id]]"
                handle-as="json"
                on-error="delete_item_failed"
                on-response="delete_item_successful"
                loading="{{loading.delete_item}}">
        </iron-ajax>

        <card-item title="[[first_defined(user.name, user.email)]]" icon="person">
            <paper-listbox slot="dropdown-content">
                <paper-item on-tap="edit_item">
                    <mwc-icon>create</mwc-icon>
                    Edit
                </paper-item>
                <paper-item on-tap="open_delete_item_confirmation">
                    <mwc-icon>delete</mwc-icon>
                    Delete
                </paper-item>
            </paper-listbox>

            <div class="infogrid">
                <div>Role</div>
                <div>[[user.role]]</div>
                <div>Email</div>
                <div class="break-word">
                    <a href$="mailto:[[user.email]]">[[user.email]]</a>
                </div>
                <div>Last Seen</div>
                <div>[[long_date(user.lastlog)]]</div>
                <div>Created</div>
                <div>[[long_date(user.creation_date)]]</div>
            </div>
        </card-item>


        <paper-dialog id="delete_item_confirmation" on-iron-overlay-closed="delete_item">
            <div>
                Delete [[user.role]] [[first_defined(user.name, user.email)]]?
            </div>
            <div class="buttons">
                <paper-button raised="" dialog-dismiss="">Cancel</paper-button>
                <paper-button raised="" dialog-confirm="">Delete</paper-button>
            </div>
        </paper-dialog>
        `;
    }
    static get properties() {
        return {
            // Must be provided by parent object.
            user: { type: Object },
        };
    }

    // Opens the edit modal.
    edit_item() {
        // Deep copy the object so we don't change the card's
        // display until the save is successful.
        this.set('item_to_edit', JSON.parse(JSON.stringify(this.user)));
        // Ask for the form to be opened.
        window.dispatchEvent(new CustomEvent("open-form", {
            detail: {
                name:       "edit_user_form",
                user:       this.item_to_edit,
                that:       this,
                callback:   "resolve_edit_item_dialog",
            },
        }));
    }

    // Handle response from dialog. Reason is either confirmed or canceled.
    resolve_edit_item_dialog(e, reason) {
        if (!reason.confirmed) return;
        // Override dirty checking; let Polymer know it changed.
        var tmp = this.item_to_edit;
        this.set("item_to_edit", {});
        this.set("item_to_edit", tmp);
        this.$.put_ajax.generateRequest();
    }

    open_delete_item_confirmation() { this.$.delete_item_confirmation.open(); }

    // Handle result of delete modal.
    delete_item(e, reason) {
        if (reason.confirmed) this.$.delete_item_ajax.generateRequest();
    }

    delete_item_successful() {
        window.dispatchEvent(new CustomEvent("success-toast", {
            detail: this.user.role + " " + this.first_defined(
                this.user.name, this.user.email) + " deleted",
        }));
        // Ask the parent element to delete this item. It can't be
        // done here because this whole element needs to be removed.
        this.dispatchEvent(new CustomEvent("delete-item", {
            bubbles: true,
            composed: true,
            detail: this.user,
        }));
    }

    delete_item_failed(e, data) {
        if (data.error) this.check_ajax_status(data.request);
        window.dispatchEvent(new CustomEvent("error-toast", {
            detail: "Failed to delete user :(",
        }));
    }

    // Copy response from PUT to update the display. The rationale
    // for not loading the PUT response directly into user is
    // to prevent a failed return status from clearing the display.
    put_item_successful() {
        this.set('user', this.item_to_edit);
        window.dispatchEvent(new CustomEvent("success-toast", {
            detail: this.user.role + " " + this.first_defined(
                this.user.name, this.user.email) + " saved",
        }));
    }

    put_item_failed(e, data) {
        if (data.error) this.check_ajax_status(data.request);
        window.dispatchEvent(new CustomEvent("error-toast", {
            detail: "Failed to save user :(",
        }));
    }
}
customElements.define("user-display", UserDisplay);
