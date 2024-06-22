# QUIC-pub-sub-app 

Messaging broker app which accepts publishes and subscribers by UDP ports (ports defined in `.env`) Using QUIC protocol.

## How to start
* Run command `go install`
* In root directory run command `go run *.go`

## Testing by running clients
In clients folder there is publisher and subscriber files which can be used for testing server.
To run the clients navigate to folder client/pub or sub and run command `go run *.go`

* Publisher will publish a message every 10 seconds.
* Subscriber will log message sent by server.

### TLS cert there just for local testing purposes.


### Original requirements for server
* Accepts QUIC connections on 2 ports.
    * Publisher port.
    * Subscriber port.
* The server notifies publishers if a subscriber has connected.
* If no subscribers are connected, the server must inform the publishers.
* The server sends any messages received from publishers to all connected subscribers.
