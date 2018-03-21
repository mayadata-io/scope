## How to run this mock server :

> First, run the weavescope UI :

- Go to scope/client: `cd scope/client`
- You need at least Node.js 6.9.0
- Get Yarn: `npm install -g yarn`
- Setup: `yarn install`
- Develop: `yarn start`

This will start a webpack-dev-server that serves the UI.

> Now ,run the mock-server
  ## Running mock-server (go-server)
- Go to scope/mock-server: `cd scope/mock-server`
- Download and install gin-gonc: `go get github.com/gin-gonic/gin`
- Download and install gorilla: `go get github.com/gorilla/websocket`
- Start mock-server : `go run mock-server.go` and then open `http://localhost:4042/`
