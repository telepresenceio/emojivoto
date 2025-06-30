# Emoji.voto

A microservice application that allows users to vote for their favorite emoji,
and tracks votes received on a leaderboard. May the best emoji win.

The application is composed of the following 3 services:

* [emojivoto-web](emojivoto-web/): Web frontend and REST API
* [emojivoto-emoji-svc](emojivoto-emoji-svc/): gRPC API for finding and listing emoji
* [emojivoto-voting-svc](emojivoto-voting-svc/): gRPC API for voting and leaderboard

![Emojivoto Topology](assets/emojivoto-topology.png "Emojivoto Topology")

## Running

### In Minikube

Deploy the application to Minikube

1. Inject, Deploy, and Enjoy

    ```bash
    make kustomize | kubectl apply -f -
    ```

1. Use the app!

    ```bash
    minikube -n emojivoto service web-svc
    ```

### In docker-compose

It's also possible to run the app with docker-compose.

Build and run:

```bash
make deploy-to-docker-compose
```

The web app will be running on port 8080 of your docker host.

### Via URL

To deploy standalone to an existing cluster:

```bash
kubectl apply -k github.com/telepresenceio/emojivoto/kustomize/deployment
```

### Generating some traffic

The `VoteBot` service can generate some traffic for you. It votes on emoji
"randomly" as follows:

- It votes for :doughnut: 15% of the time.
- When not voting for :doughnut:, it picks an emoji at random

If you're running the app using the instructions above, the VoteBot will have
been deployed and will start sending traffic to the vote endpoint.

If you'd like to run the bot manually:

```bash
export WEB_HOST=localhost:8080 # replace with your web location
go run emojivoto-web/cmd/vote-bot/main.go
```

## Building

Building requires that you have protoc-gen-go v1.27.1 and
protoc-gen-go-grpc v1.1.0 on your path. These can be installed by running:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

Then you can set up proto files and build apps by running:

```bash
make build
```

## Prometheus Metrics

By default the voting service exposes Prometheus metrics about current vote count on port `8801`.

This can be disabled by unsetting the `PROM_PORT` environment variable.

## Local Development

### Emojivoto webapp

This app is written with React and bundled with webpack.
Use the following to run the emojivoto go services and develop on the frontend.

Start the voting service

```bash
GRPC_PORT=8081 go run emojivoto-voting-svc/cmd/server.go
```

[In a separate terminal window] Start the emoji service

```bash
GRPC_PORT=8082 go run emojivoto-emoji-svc/cmd/server.go
```

[In a separate terminal window] Bundle the frontend assets

```bash
cd emojivoto-web/webapp
yarn install
yarn webpack # one time asset-bundling OR
yarn webpack-dev-server --port 8083 # bundle/serve reloading assets
```

[In a separate terminal window] Start the web service

```bash
export WEB_PORT=8080
export VOTINGSVC_HOST=localhost:8081
export EMOJISVC_HOST=localhost:8082

# if you ran yarn webpack
export INDEX_BUNDLE=emojivoto-web/webapp/dist/index_bundle.js

# if you ran yarn webpack-dev-server
export WEBPACK_DEV_SERVER=http://localhost:8083

# start the webserver
go run emojivoto-web/cmd/server.go
```

[Optional] Start the vote bot for automatic traffic generation.

```bash
export WEB_HOST=localhost:8080
go run emojivoto-web/cmd/vote-bot/main.go
```

View emojivoto

```bash
open http://localhost:8080
```
