import { dedupingMixin } from '@polymer/polymer/lib/utils/mixin.js';

let m = (base) => class extends base {
    open() {
        // Call an open hook function if it exists. A form may use this to do
        // setup or trigger AJAX calls.
        if (this.open_hook) this.open_hook();
        this.$.dialog.open();
    }
    close() { this.$.dialog.close(); }
    resolve_dialog(e, reason) {
        this.dispatchEvent(new CustomEvent("closed", {
            bubbles: true,
            composed: true,
            detail: reason,
        }));
    }
    // Handle click on the X suffix for a paper-input. This crawls up the DOM
    // until it finds a paper-input and clears it.
    clear_field(e) {
        var elem = e.target;
        while (elem = elem.parentElement) {
            if (elem.localName === "paper-input")
                return elem.value = null;
        }
    }
    // Same, but for a number field.
    clear_number_field(e, p) {
        var elem = e.target;
        while (elem = elem.parentElement) {
            if (elem.localName === "paper-input")
                return elem.value = 0;
        }
    }
};
export const FormMixin = dedupingMixin(m);
