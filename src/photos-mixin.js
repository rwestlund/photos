import { dedupingMixin } from '@polymer/polymer/lib/utils/mixin.js';

let m = (base) => class extends base {
    static get properties() {
        return {
            constants: {
                type: Object,
                value: { user_roles: ['Admin', 'User'] },
            },
        };
    }

    pretty_image_count(c) { return c + (c === 1 ? " photo" : " photos"); }

    user_is_admin(role) { return (role === 'Admin'); }

    // This is used in the DOM in various places.
    equal(a, b) { return a === b; }
    long_date(d) {
        if (!d) return '';
        var date = new Date(d)
        return date.toDateString() + ' ' + date.toLocaleTimeString();
    }
    // This is used in the DOM in various places.
    first_defined(a, b) { return a || b; }
};
export const PhotosMixin = dedupingMixin(m);
