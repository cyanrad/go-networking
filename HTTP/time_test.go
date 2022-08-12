package HTTP

import (
	"net/http"
	"testing"
	"time"
)

func iferr(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestTime(t *testing.T) {
	// >> sending request and getting the response
	resp, err := http.Get("https://www.time.gov/") // No body in this since we don't
	// want to render/ transfer resources
	iferr(err, t)

	// The default HTTP client's Transport may not
	// reuse HTTP/1.x "keep-alive" TCP connections if the Body is
	// not read to completion and closed.
	_ = resp.Body.Close()

	t.Log(resp.Status, resp.ContentLength)

	now := time.Now().Round(time.Second) // getting the system time
	date := resp.Header.Get("Date")      // getting the date from the time.gov response
	if date == "" {                      // if we get no date in the header
		t.Fatal("Big issue here :<\nno Date header received from time.gov")
	}

	// parsing the recived date into the RFC1123 standard
	// so we can subtract it from the system  time
	dt, err := time.Parse(time.RFC1123, date)
	iferr(err, t)

	// logging the difference between them
	t.Logf("time.gov: %s (skew %s)", dt, now.Sub(dt))
}
