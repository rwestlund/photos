/*
    Copyright (c) 2016-2017, Randy Westlund and Jacqueline Kory Westlund.
    All rights reserved.
    This code is under the BSD-2-Clause license.
*/
/* This module shows a list of items, with server-side pagination and search. */
import '@polymer/iron-ajax/iron-ajax.js';
import '@polymer/iron-flex-layout/iron-flex-layout-classes.js';
import '@polymer/iron-media-query/iron-media-query.js';
import '@polymer/paper-icon-button/paper-icon-button.js';
import '@polymer/paper-fab/paper-fab.js';
import '@polymer/paper-icon-button/paper-icon-button.js';
import '@polymer/paper-input/paper-input.js';
import '@polymer/paper-spinner/paper-spinner.js';
import '@polymer/polymer/lib/elements/dom-if.js';
import '@polymer/polymer/lib/elements/dom-repeat.js';
import { PolymerElement, html } from '@polymer/polymer/polymer-element.js';

import { FormMixin } from './form-mixin.js';
import { PhotosMixin } from './photos-mixin.js';
import './global-styles.js';
import './user-display.js';

class ItemCollection extends FormMixin(PhotosMixin(PolymerElement)) {
    static get template() {
        return html`
        <style include="iron-flex iron-flex-alignment iron-flex-reverse"></style>
        <style include="global-styles"></style>
        <style>
            :host {
                display: block;
            }
            .nav-buttons {
                margin-top: 5px;
                margin-left: 2em;
            }
            .search-box {
                min-width: 10em;
            }
            .nav-buttons paper-icon-button {
                padding-bottom: 0px;
            }
        </style>

        <iron-media-query query="(max-width: 600px)" query-matches="{{is_mobile}}">
        </iron-media-query>

        <!-- AJAX requests. -->
        <iron-ajax id="get_items_ajax"
                auto=""
                method="GET"
                url="/api/[[itemName]]"
                params="[[search_filter]]"
                handle-as="json"
                last-response="{{items}}"
                debounce-duration="100"
                on-error="loading_data_failed"
                on-response="new_items_received"
                loading="{{loading.get_items}}">
        </iron-ajax>
        <iron-ajax id="create_item_ajax"
                method="POST"
                url="/api/[[itemName]]"
                body="[[new_item]]"
                content-type="application/json"
                handle-as="json"
                last-response="{{new_item}}"
                on-error="creating_item_failed"
                on-response="creating_item_succeeded"
                loading="{{loading.post_item}}">
        </iron-ajax>

        <div class="layout horizontal-reverse center end wrap">
            <!-- Navigation icons. -->
            <div class="layout horizontal center nav-buttons">
                <paper-spinner class="nav-button" active="[[loading.get_items]]">
                </paper-spinner>
                <paper-icon-button icon="icons:refresh" on-tap="refresh">
                </paper-icon-button>
                <paper-icon-button icon="icons:arrow-back" disabled\$="[[!skip]]" on-tap="previous">
                </paper-icon-button>
                <paper-icon-button icon="icons:arrow-forward" disabled\$="[[disable_next]]" on-tap="next">
                </paper-icon-button>
                <div class="nav-button">Page&nbsp;[[page_number]]</div>
            </div>
            <!-- Search box. -->
            <paper-input type="text" class="flex search-box" label="Search" no-label-float="" value="{{search_text}}">
                <paper-icon-button slot="suffix" icon="icons:clear" on-tap="clear_field">
                </paper-icon-button>
            </paper-input>
        </div>

        <!-- Display a list of whichever item we're showing. -->
        <template is="dom-if" if="[[equal(itemName, 'albums')]]">
            <template is="dom-repeat" items="[[items]]">
                <!--TODO merge album-collection into here. -->
            </template>
        </template>

        <template is="dom-if" if="[[equal(itemName, 'users')]]">
            <template is="dom-repeat" items="[[items]]">
                <!-- This is bound with {{}} because it can edit/PUT data. -->
                <user-display user="{{item}}" on-delete-item="remove_item">
                </user-display>
            </template>
        </template>

        <!-- FAB to trigger creation forms. -->
        <template is="dom-if" if="[[show_fab(itemName)]]">
            <paper-fab icon="icons:add" on-tap="create_item">
            </paper-fab>
        </template>


        <!-- Show bottom buttons if there are several items listed. -->
        <template is="dom-if" if="[[show_bottom_buttons(items)]]">
            <div class="layout horizontal end-justified center nav-buttons">
                <paper-spinner class="nav-button" active="[[loading.get_items]]">
                </paper-spinner>
                <paper-icon-button icon="icons:arrow-back" disabled\$="[[!skip]]" on-tap="previous">
                </paper-icon-button>
                <paper-icon-button icon="icons:arrow-forward" disabled\$="[[disable_next]]" on-tap="next">
                </paper-icon-button>
                <div class="nav-button">Page&nbsp;[[page_number]]</div>
            </div>
        </template>
        `;
    }
    static get properties() {
        return {
            // The parent provides 'users' or 'albums' here.
            itemName: { type: String },
            items: { type: Array },
            count: { type: Number, value: 20 },
            skip: { type: Number, value: 0 },
            page_number: {
                type: Number,
                computed: "compute_page_number(skip)"
            },
            search_text: { type: String },
            search_filter: {
                type: Object,
                computed: "compute_search_filter(count, skip, "
                + "search_text)"
            },
            disable_next: {
                type: Boolean,
                computed: "compute_disable_next(items, count)"
            },
        };
    }

    compute_search_filter(count, skip, search_text) {
        var o = { count: count };
        if (skip) o.skip = skip;
        if (search_text) o.query = search_text;
        return o;
    }

    compute_disable_next(items, count) {
        var val = true;
        if (this.items)
            val = this.items.length < this.count;
        return val;
    }

    compute_page_number(skip) { return skip + 1; }

    // Reload the collection of items.
    refresh() { this.$.get_items_ajax.generateRequest(); }

    // Open the appropriate form after a click on the FAB.
    create_item() {
        this.set('new_item', {});
        var data = { that: this, callback: "resolve_create_item" };

        switch (this.itemName) {
            case "users":
                data.name = "create_user_form";
                data.user = this.new_item;
                break;
        }
        // Ask for the form to be opened, whichever one it is.
        window.dispatchEvent(new CustomEvent("open-form", {
            detail: data
        }));
    }

    // Handle response from dialog. Reason is either confirmed or canceled.
    resolve_create_item(e, reason) {
        if (!reason.confirmed) return;
        // Override dirty checking; let Polymer know it changed.
        var tmp = this.new_item;
        this.set("new_item", {});
        this.set("new_item", tmp);
        this.$.create_item_ajax.generateRequest();
    }

    creating_item_succeeded() {
        // This isn't in shorter form like the one below because
        // some pages need to change the route.
        if (this.itemName === "customers") {
            window.dispatchEvent(new CustomEvent("success-toast", {
                detail: "Customer " +
                this.customer_full_name(this.new_item) + " created",
            }));
            // Go to the customers page for the new customer.
            window.history.pushState({}, null,
                "/customers/" + this.new_item.id);
            window.dispatchEvent(new CustomEvent("location-changed"));
        }
        else if(this.itemName === "users") {
            window.dispatchEvent(new CustomEvent("success-toast", {
                detail: this.new_item.role + " " +
                (this.new_item.name || this.new_item.email)
                + " created",
            }));
            this.push("items", this.new_item);
        }
    }

    creating_item_failed(e, data) {
        if (data.error) this.check_ajax_status(data.request);
        var msg;
        if (this.itemName === "customers")
            msg = "Failed to create customer :(";
        else if (this.itemName === "users")
            msg = "Failed to create user :(";

        window.dispatchEvent(new CustomEvent("error-toast", {
            detail: msg,
        }));
    }

    // Simply remove the element from the DOM.
    remove_item(e) { this.splice('items', this.items.indexOf(e.detail), 1); }

    previous() { if (this.skip) this.skip--; }

    next() { this.skip++; }

    // When new items are received, back up a page if the current one is blank.
    new_items_received() {
        if (this.items && !this.items.length && this.skip)
            this.skip--;
    }

    // Whether to show the FAB or not. The item must be an arg for
    // Polymer to evaluate it properly.
    show_fab(item_name) { return true; }

    // Only show bottom nav buttons if there are several items.
    show_bottom_buttons(items) {
        if (!items) return false;
        return items.length > 5;
    }
}
customElements.define("item-collection", ItemCollection);
