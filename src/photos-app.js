/*
    Copyright (c) 2016-2017, Randy Westlund and Jacqueline Kory Westlund.
    All rights reserved.
    This code is under the BSD-2-Clause license.
*/
import '@polymer/app-layout/app-layout.js';
import '@polymer/app-route/app-location.js';
import '@polymer/app-route/app-route.js';
import '@polymer/iron-ajax/iron-ajax.js';
import '@polymer/iron-media-query/iron-media-query.js';
import '@polymer/iron-pages/iron-pages.js';
import '@polymer/paper-button/paper-button.js';
import '@polymer/paper-icon-button/paper-icon-button.js';
import '@polymer/paper-item/paper-item.js';
import '@polymer/paper-listbox/paper-listbox.js';
import '@polymer/paper-spinner/paper-spinner.js';
import '@polymer/polymer/lib/elements/dom-if.js';
import '@material/mwc-icon/mwc-icon.js';
import { PolymerElement, html } from '@polymer/polymer/polymer-element.js';

import { PhotosMixin } from './photos-mixin.js';
import './album-collection.js';
import './cookie-display.js';
import './global-styles.js';
import './item-collection.js';
import './photo-collection.js';
import './photos-forms.js';
import './photos-uploads.js';

class PhotosApp extends PhotosMixin(PolymerElement) {
    static get template() {
        return html`
        <style include="global-styles"></style>
        <style>
            :host {
                display: block;
                /* Primary background color is also defined in index.html */
                --primary-background-color: #f7f0ed;
                --medium-background-color: #f2dfd9;
                --dark-background-color: #d18e7a;
                --header-background-color: #610505;
                --light-header-text-color: #fbf2ea;
                --dark-header-text-color: #2a1a09;

                --paper-fab-background: var(--header-background-color);

                /* These apply to the application menu drawer. */
                --app-drawer-content-container: {
                    /* There is no border by default. */
                    border-right: 1px solid var(--dark-background-color);
                    background-color: var(--medium-background-color);
                }

                /* Make the drawer a little smaller than default. */
                --app-drawer-width: 10em;
            }

            app-drawer-layout:not([narrow]) [drawer-toggle] {
                  display: none;
            }
            paper-listbox {
                background-color: var(--medium-background-color);
            }

            /* This applies to the main application title. */
            div[main-title] {
                font-size: x-large;
                font-weight: bold;
                color: var(--light-header-text-color);
                text-shadow: 1px 1px var(--dark-header-text-color),
                    0 0 5px var(--dark-header-text-color);
            }
            /* This applies to the application header/toolbar. */
            app-toolbar {
                background-color: var(--header-background-color);
                background-image: url("images/header_leaves.jpg");
                background-repeat: repeat-x;
            }
            @media(max-width: 1400px) {
                app-toolbar {
                    background-size: cover;
                }
            }
            hr {
                border-color: var(--dark-background-color);
            }
            iron-pages {
                margin: 15px;
            }
            paper-listbox a {
                text-decoration: none;
                color: #111111;
            }
            /* For selected menu items in the drawer. */
            paper-listbox .iron-selected paper-item {
                background-color: var(--dark-background-color);
                font-weight: bold;
            }

            span.username {
                padding-left: 5px;
                font-size: large;
                color: var(--light-header-text-color);
                text-shadow: 1px 1px var(--dark-header-text-color),
                    0 0 5px var(--dark-header-text-color);
            }
            paper-icon-button.header {
                --paper-icon-button-ink-color: var(--light-header-text-color);
                color: var(--light-header-text-color);
            }
            app-header mwc-icon {
                color: var(--light-header-text-color);
            }
            paper-listbox mwc-icon {
                margin-right: 5px;
            }
        </style>

        <iron-ajax id="post_ajax"
                    method="POST"
                    url="/api/albums"
                    body="[[item_to_edit]]"
                    content-type="application/json"
                    handle-as="json"
                    last-response="{{item_to_edit}}"
                    on-response="post_item_successful"
                    on-error="post_item_failed"
                    loading="{{loading}}">
        </iron-ajax>

        <cookie-display cookie-name="username" cookie-value="{{user_name}}">
        </cookie-display>
        <cookie-display cookie-name="role" cookie-value="{{user_role}}">
        </cookie-display>

        <iron-media-query query="(max-width: 500px)" query-matches="{{mobile}}">
        </iron-media-query>

        <!-- App drawer -->
        <app-drawer-layout id="drawer_layout" fullbleed responsive-width="900px">
            <app-drawer id="drawer" slot="drawer" swipe-open>
                <app-toolbar></app-toolbar>
                <paper-listbox selected="[[route_data.page]]"
                        attr-for-selected="name"
                        on-tap="toggle_drawer">
                    <a name="" href="/">
                        <paper-item>
                            <mwc-icon>photo</mwc-icon>
                            Recent
                        </paper-item>
                    </a>
                    <a name="albums" href="/albums/">
                        <paper-item>
                            <mwc-icon>photo_library</mwc-icon>
                            Albums
                        </paper-item>
                    </a>
                    <a name="about" href="/about/">
                        <paper-item>
                            <mwc-icon>info</mwc-icon>
                            About
                        </paper-item>
                    </a>
                    <template is="dom-if" if="[[user_is_admin(user_role)]]">
                        <a name="users" href="/users/">
                            <paper-item>
                                <mwc-icon>people</mwc-icon>
                                Users
                            </paper-item>
                        </a>
                        <a name="uploads" href="/uploads/">
                            <paper-item>
                                <mwc-icon>cloud_upload</mwc-icon>
                                Uploads
                            </paper-item>
                        </a>
                    </template>
                </paper-listbox>
                <template is="dom-if" if="[[user_is_admin(user_role)]]">
                    <hr>
                    <paper-button raised on-tap="create_album">
                        <mwc-icon>add_photo_alternate</mwc-icon>
                        Album
                    </paper-button>
                    <paper-spinner active="[[loading]]">
                    </paper-spinner>
                </template>

                <hr>
                <paper-listbox>
                <template is="dom-if" if="[[!user_name]]">
                    <a href="/api/auth/google/login">
                        <paper-item>
                            <mwc-icon>account_circle</mwc-icon>
                            Sign In
                        </paper-item>
                    </a>
                </template>
                <template is="dom-if" if="[[user_name]]">
                    <a href="/api/auth/logout">
                        <paper-item>
                           <mwc-icon>directions_run</mwc-icon>
                           Logout
                        </paper-item>
                    </a>
                </template>
                </paper-listbox>

            </app-drawer>

            <!-- App header -->
            <app-header-layout fullbleed>
                <app-header slot="header" fixed>
                    <app-toolbar>
                        <paper-icon-button class="header" icon="menu" drawer-toggle>
                        </paper-icon-button>
                        <div main-title>[[get_page_name(route_data.page)]]</div>
                        <template is="dom-if" if="[[!mobile]]">
                            <template is="dom-if" if="[[user_name]]">
                                <mwc-icon>account_circle</mwc-icon>
                                <span class="username">Hi, [[user_name]]!</span>
                            </template>
                        </template>
                    </app-toolbar>
                </app-header>

                <!-- Ignore /api/ and /s/ routes; let the browser have them. This
                     is also defined in sw-precache-config.js. -->
                <app-location route="{{route}}" url-space-regex="^/(?!(api|s)/)">
                </app-location>
                <!-- The top-level router. -->
                <app-route
                    route="{{route}}"
                    pattern="/:page"
                    data="{{route_data}}"
                    tail="{{tail}}">
                </app-route>

                <!-- If every page used the same subrouter, they'd all be bound
                 to the same id (i.e. loading albums/test would also try to
                 load otherurl/test). Instead, use a separate router for each path
                 so that a non-active path will not match the id from the
                 current path. There may be a better way to do this. -->
                <app-route
                     route="{{route}}"
                     pattern="/albums/:album"
                     data="{{albums_route_data}}">
                </app-route>

                <!-- App pages -->
                <iron-pages selected="[[route_data.page]]"
                        attr-for-selected="name">

                    <section name="">
                         <!-- Recent photos -->
                        <photo-collection>
                        </photo-collection>
                    </section>

                    <section name="albums">
                         <!-- Albums -->
                        <template is="dom-if" if="[[albums_route_data.album]]">
                            <photo-collection album-name="[[albums_route_data.album]]">
                            </photo-collection>
                        </template>

                        <template is="dom-if" if="[[!albums_route_data.album]]">
                            <album-collection></album-collection>
                        </template>
                    </section>

                    <section name="about">
                         <!-- About -->
                        <p>Welcome to our photo site!</p>

                        <p>Family and friends will have the ability to log in and
                        view photos. If you don't have an account but think you
                        should, let us know and we can sort it out.</p>

                        <p>The source for this site is on <a
                            href="https://github.com/rwestlund/photos">Github</a>
                        under the BSD-2-Clause license.</p>
                    </section>

                    <template is="dom-if" if="[[user_is_admin(user_role)]]">
                        <section name="users">
                             <!-- Users -->
                            <item-collection item-name="users">
                            </item-collection>
                        </section>

                        <section name="uploads">
                             <!-- Uploads -->
                            <photos-uploads></photos-uploads>
                        </section>
                    </template>
                </iron-pages>
            </app-header-layout>
        </app-drawer-layout>

        <photos-forms></photos-forms>
        `;
    }
    // Define importPath, which makes the CSS header path rewrite properly.
    static get importMeta() {
        return import.meta;
    }
    toggle_drawer() {
        if (this.$.drawer_layout.narrow) this.$.drawer.toggle();
    }
    // Opens the create album modal.
    create_album() {
        this.set('item_to_edit', {});
        // Ask for the form to be opened.
        window.dispatchEvent(new CustomEvent("open-form", {
            detail: {
                name: "create_album_form",
                album: this.item_to_edit,
                callback: "resolve_create_album",
                that: this,
            },
        }));
    }
    // Handle response from dialog. Reason is either confirmed or canceled.
    resolve_create_album(e, reason) {
        if(!reason.confirmed) return;
        // Override dirty checking; let Polymer know it changed.
        var tmp = this.item_to_edit;
        this.set("item_to_edit", {});
        this.set("item_to_edit", tmp);
        this.$.post_ajax.generateRequest();
    }
    post_item_successful() { console.log("updated"); }
    post_item_failed() { console.log("failed to update") }
    get_page_name(page) {
        switch (page) {
            case "": return "Recent Photos";
            case "albums": return "Albums";
            case "about": return "About";
            case "users": return "Users";
            case "uploads": return "Uploads";
            default: return "Photos";
        }
    }
}
customElements.define("photos-app", PhotosApp);
