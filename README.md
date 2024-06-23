# QUIC-pub-sub-app 

Messaging broker app which accepts publishes and subscribers by UDP ports (ports defined in `.env`) Using QUIC protocol.

## How to start
* Run command `go install`
* In root directory run command `go run *.go`

## Testing by running clients
In clients folder there is publisher and subscriber files which can be used for testing server.(Please note clients created for testing purposes only, so they are just as good as)
To run the clients navigate to folder client/pub or sub and run command `go run *.go`

* Publisher will publish a message every 10 seconds.
* Subscriber will log message sent by server.

## Build docker image
* Run command `docker-compose up` This will start docker container with ports by default 4444 for pub and 5555 for sub or alternately you can build image manually:
    * Build command `docker build -t quic-pub-sub-app . `
    * Run command `docker run -d -p 4444:4444/udp -p 5555:5555/udp --name quic-pub-sub-app-container quic-pub-sub-app `

### TLS cert there just for local testing purposes.


### Original requirements for server
* Accepts QUIC connections on 2 ports.
    * Publisher port.
    * Subscriber port.
* The server notifies publishers if a subscriber has connected.
* If no subscribers are connected, the server must inform the publishers.
* The server sends any messages received from publishers to all connected subscribers.
