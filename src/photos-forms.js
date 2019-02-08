/*
    Copyright (c) 2016-2019. Randy Westlund and Jacqueline Kory Westlund.
    All rights reserved.
    This code is under the BSD-2-Clause license.
*/
import { PolymerElement, html } from '@polymer/polymer/polymer-element.js';

import './album-form.js';
import './expand-photo-form.js';
import './photo-form.js';
import './select-album-form.js';
import './user-form.js';

class PhotosForms extends PolymerElement {
    static get template() {
        return html`
        <photo-form id="edit_photo_form"></photo-form>
        <user-form id="create_user_form" title="Create new user"></user-form>
        <user-form id="edit_user_form"></user-form>
        <album-form id="create_album_form"></album-form>
        <album-form id="edit_album_form" title="Edit album name"></album-form>
        <select-album-form id="select_album_form"></select-album-form>
        <expand-photo-form id="expand_photo_form" with-backdrop="">
        </expand-photo-form>
        `;
    }
    static get properties() {
        return {
            // A reference to the element that asked to open a form.
            dialog_parent: { type: Object },
            // The name of the callback function on the requesting element.
            dialog_callback: { type: String },
        };
    }

    connectedCallback() {
        super.connectedCallback();
        this._form_listener = this.open_form.bind(this);
        window.addEventListener("open-form", this._form_listener);
        // Every form fires a "closed" event. Rather than binding a
        // listener to each one, let it bubble up and catch them all here.
        this._closed_listener = this.dialog_closed.bind(this);
        this.addEventListener("closed", this._closed_listener);
    }
    disconnectedCallback() {
        super.disconnectedCallback();
        window.removeEventListener("open-form", this._form_listener);
        window.removeEventListener("closed", this._closed_listener);
    }

    // On close, notify the element that asked for the form.
    dialog_closed(e) {
        if (this.dialog_parent)
            this.dialog_parent[this.dialog_callback](e, e.detail);
    }

    // Open whichever form was requested by the event.
    open_form(e) {
        // Save a reference to the element that sent the request so we
        // can tell it when the form closes.
        this.dialog_parent = e.detail.that;
        // Save the callback.
        this.dialog_callback = e.detail.callback;
        // If any of these properties are given from the event, assign
        // them to the form. This prevents needing a large switch to set
        // up each form.
        var props = [ "photo", "album", "user", "selectedAlbums",
            "photoId", "photoList" ];
        props.forEach(p => {
            if (e.detail[p]) this.$[e.detail.name][p] = e.detail[p];
        });
        this.$[e.detail.name].open();
    }
}
customElements.define("photos-forms", PhotosForms);
