const $_documentContainer = document.createElement('template');

$_documentContainer.innerHTML = `<title>photos</title><style>
        body {
            background-color: #f7f0ed;
        }
    </style><photos-app></photos-app>`;

document.head.appendChild($_documentContainer.content);

/*
    Copyright (c) 2016. Randy Westlund and Jacqueline Kory Westlund.
    All rights reserved.
    This code is under the BSD-2-Clause license.
*/
/*
  FIXME(polymer-modulizer): the above comments were extracted
  from HTML and may be out of place here. Review them and
  then delete this comment!
*/
;
