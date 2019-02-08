const $_documentContainer = document.createElement('template');

$_documentContainer.innerHTML = `<dom-module id="global-styles">
    <template>
        <style>
            a.plain {
                color: inherit;
                text-decoration: none;
            }
            div.paper-material.card-item {
                padding: 10px;
                margin-bottom: 10px;
            }
            div.paper-material {
                background-color: white;
            }
            paper-fab {
                position: fixed;
                bottom: 1em;
                right: 2em;
                z-index: 100;
            }
            paper-button {
                background-color: white;
            }
            .nav-button {
                margin-left: 0.8em;
                margin-right: 0.8em;
            }
            .card-title {
                margin-bottom: 0;
                margin-top: 0;
            }
            iron-icon.large-icon {
                --iron-icon-width: 4em;
                --iron-icon-height: 4em;
                --iron-icon-fill-color: gray;
            }
            /* This makes a table able to wrap long normally unbreakable
               content like email addresses. The two td classes below should be
               used for data. The 80% width is to allow an icon to the left in
               a flexbox layout.  */
            table.fixed-80 {
                table-layout: fixed;
                width: 80%;
            }
            td.td-label {
                vertical-align: top;
                text-align: right;
                padding-right: 0.8em;
                /* 4em is long enough for phone|email|fax. */
                width: 4em;
            }
            td.td-field {
                overflow: hidden;
                word-wrap: break-word;
            }

            .break-word {
                overflow: hidden;
                word-wrap: break-word;
            }

            /* Contains a list of divs. On mobile, stacks with headings and
               values. On wider screens, looks like a table. */
            div.infogrid, div.infogrid-wide {
                display: grid;
                justify-content: start;
                column-gap: 1em;
                /* Firefox 60 still needs the grid- version */
                grid-column-gap: 1em;
            }
            /* These are the headings. */
            div.infogrid > div:nth-of-type(odd),
            div.infogrid-wide > div:nth-of-type(odd) {
                font-weight: bold;
            }
            /* Every heading except the first. */
            div.infogrid > div:nth-of-type(2n+3),
            div.infogrid-wide > div:nth-of-type(2n+3) {
                margin-top: 1em;
            }
            div.infogrid-wide {
                @apply --infogrid-wide;
            }
            div.infogrid-wide > div:nth-of-type(odd) {
                @apply --infogrid-children-wide;
            }
            @media (min-width: 400px) {
                div.infogrid {
                    @apply --infogrid-wide;
                }
                div.infogrid > div:nth-of-type(odd) {
                    @apply --infogrid-children-wide;
                }
            }
            :host {
                --infogrid-wide: {
                    grid-template-columns: auto auto;
                    row-gap: 0.3em;
                    grid-row-gap: 0.3em;
                };
                --infogrid-children-wide: {
                    text-align: right;
                    /* Removing the margin from the first child too doesn't
                       hurt. */
                    margin-top: 0;
                };
            }
        </style>
    </template>
</dom-module>`;

document.head.appendChild($_documentContainer.content);

/*
    Copyright (c) 2016, Randy Westlund. All rights reserved.
    This code is under the BSD-2-Clause license.
*/
/* This defines CSS that is imported by every element. */
/*
  FIXME(polymer-modulizer): the above comments were extracted
  from HTML and may be out of place here. Review them and
  then delete this comment!
*/
;
