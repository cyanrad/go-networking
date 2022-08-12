package HTTP

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func blockIndefinitely(w http.ResponseWriter, r *http.Request) {
	select {} // vlocks indefinitly
}

func TestTimeoutHTTP(t *testing.T) {
	// >> creating the test server
	testServer := httptest.NewServer(http.HandlerFunc(blockIndefinitely))

	// >> creating the hearbeat context
	//								   The parent context    wait till cancel
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// >> creating request with context
	// unlike the default one with no timeout
	//											the req Method  server to send to, the io read for response
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, testServer.URL, nil)
	if err != nil {
		t.Log(err)
	}

	// >> sending the request, and getting the response
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Log(err)
		} else {
			t.Log("Request Deadline Exceeded")
		}
		return
	}

	_ = resp.Body.Close()
}
