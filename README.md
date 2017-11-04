# Photos
Photos is a web application for sharing family photos, built with Polymer and
Go. Our deployment is at
[https://photos.textplain.net](https://photos.textplain.net).

# Run Dependencies
- PostgreSQL >= 9.5
- NGINX

# Build Dependencies
- Go
- Node.js and NPM

# Installation
- Check out this repo under your `$GOPATH`
- Copy `config.toml.example` to `config.toml` and set your values.
- Run `govendor sync`
- Run `npm install`
- Run `npm run bower install`
- Run `npm run polymer build`
- Run `go build`
- Run `go run tools/createdb/main.go`
- Run `go run tools/resetdb/main.go`
- Run `./photos`
- Add a block like this to your NGINX config:
    ```nginx
    server {
        listen 443 ssl;
        server_name photos.textplain.net;
        ssl_certificate /usr/local/etc/letsencrypt/live/photos.textplain.net/fullchain.pem;
        ssl_certificate_key /usr/local/etc/letsencrypt/live/photos.textplain.net/privkey.pem;

        client_max_body_size 100M;
        proxy_request_buffering off;
        add_header Cache-Control "public, max-age=0";
        root /usr/home/randy/go/src/github.com/rwestlund/photos/build/default/;

        location /s/ {
            alias /usr/home/randy/go/src/github.com/rwestlund/photos/build/default/;
        }
        location /api/ {
            proxy_pass http://localhost:3000;
        }
        # Support refresh when client routes leak to server.
        location / {
            rewrite ^ /s/index.html;
        }
        location /service-worker.js {
            rewrite ^ /s/service-worker.js;
        }
    }
    ```

Optionally, use [paladin](https://github.com/rwestlund/paladin) to supervise
it.

# Testing
Run `govendor test --cover +local`.

# License
BSD-2-Clause
