package request

import (
	"io"
	"log"
	"net"
	"os"
)

const request string = `GET /ass HTTP/1.1
Host: example.com
Connection: keep-alive
Upgrade-Insecure-Requests: 1
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
Accept-Encoding: gzip, deflate
Accept-Language: en-US,en;q=0.9`

// Monitor embeds a log.Logger meant for logging network traffic.
type Monitor struct {
	*log.Logger
}

// Write implements the io.Writer interface.
// since MultiWriter and TeeReader require io.Writer
func (m *Monitor) Write(p []byte) (int, error) {
	return len(p), m.Output(2, string(p))
}

func Send_request() {
	monitor := &Monitor{Logger: log.New(os.Stdout, "monitor: ", 0)}

	client, err := net.Dial("tcp", "example.com:http")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	//w := io.MultiWriter(client, monitor)
	_, err = client.Write([]byte(request)) // echo the message
	if err != nil && err != io.EOF {
		monitor.Println(err)
		return
	}

	monitor.Println("done")
	response := make([]byte, 50000)
	_, err = client.Read(response)
	if err != nil && err != io.EOF {
		monitor.Println(err)
		return
	}
	log.Print(string(response))
}
