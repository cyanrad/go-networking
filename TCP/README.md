# TCP tests
Includes several tests involving the tcp4 client and server nodes.
Each file testing a TCP feature provided by the `net` package.

Mainly used to better understand the functionality of TCP.

<br/>

`tcp`
- Creating tcp4 server
- Creating tcp4 client, and sending data to server

`ctx_timeout`
- Creating a client that immidiately times out, to test the error functionality
- Creating a client with a deadline, and testing context expiry
- Creating multiple clients all with the same context, then canceling the context

`deadline`
- Using the deadline functionality of a server to control the timeouts 

`pinger`
- Creating and testing a pinging function, to ping some server
(example should be used in main.go)

`ICMP`
- A pinger for measuring time to establish a connection
(used as either a CLI tool, or a normal function)

`proxy` (work in progress)
- Creating a proxy server 
(using `io.Copy()` between `io.Reader` and `io.Writer`)

`chat` (work in progress)
- Creating a TCP chat client and server 
(more detail in file README)