package HTTP

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type User struct {
	First string
	Last  string
}

// >> returns a function for handeling POST requests

func handlePostUser(t *testing.T) func(http.ResponseWriter, *http.Request) {
	// the returned function
	return func(w http.ResponseWriter, r *http.Request) {
		// closing the request's body after the code is done
		defer func(r io.ReadCloser) {
			_, _ = io.Copy(ioutil.Discard, r) // clearing the read channel
			_ = r.Close()                     // closing it
		}(r.Body)

		// Checking if we got the correct http request method
		if r.Method != http.MethodPost {
			http.Error(w, "", http.StatusMethodNotAllowed) // (405)
			return
		}

		// creating the user
		var u User
		// decoding the json and storing it in the new user
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			t.Error(err)
			http.Error(w, "Decode Faild", http.StatusBadRequest) // (400)
			return
		}

		if r.Header.Get("content-type") != "application/json" {
			http.Error(w, "Incorrect/Missing Content Type", http.StatusBadRequest) // (400)
			t.Log(r.Header.Get("content-type"))
			t.Fatal("Incorrect/Missing Content Type")
			return
		}

		// if all went well the response is accepted (202)
		w.WriteHeader(http.StatusAccepted)

		t.Logf("First: %s, Last: %s", u.First, u.Last)
	}
}

func TestPostUser(t *testing.T) {
	// 							takes a function with the same args as
	//                          handlePostUser returned function
	ts := httptest.NewServer(http.HandlerFunc(handlePostUser(t)))
	defer ts.Close()

	// will not error out, but will return 405(method not allowed)
	// since the handler func we created only handles POST
	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	// if the statusCode is not what we expect
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d; actual status %d",
			http.StatusMethodNotAllowed, resp.StatusCode)
	}

	buf := new(bytes.Buffer)                   // creating the buffer that holds the json data
	u := User{First: "Adam", Last: "Woodbeck"} // creating the user to POST
	err = json.NewEncoder(buf).Encode(&u)      // writing the json into buf
	if err != nil {
		t.Fatal(err)
	}

	// making the POST request
	//                 the server, content type(of body) the json data
	resp, err = http.Post(ts.URL, "application/json", buf)
	if err != nil {
		t.Fatal(err)
	}

	// if we don't get accepted status code (202)
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("expected status %d; actual status %d",
			http.StatusAccepted, resp.StatusCode)
	}

	_ = resp.Body.Close() // closing the resp body
}

func TestMultipartPost(t *testing.T) {
	// creating a byte buffer as the request body
	reqBody := new(bytes.Buffer)
	w := multipart.NewWriter(reqBody) // writer that wraps the buffer

	// instead of using a stuct i suppose
	// constructing a map and looping over it
	for k, v := range map[string]string{
		"data":        time.Now().Format(time.RFC3339),
		"description": "Attached files",
	} {
		// >> Write Form field into its own part
		// seperate each form field into its own part
		// remember the parts seperated by a string thing
		// with the header stuff idk
		err := w.WriteField(k, v)
		if err != nil {
			t.Fatal(err)
		}
	}

	// looping over files to send
	for i, file := range []string{
		"./files/testing.jpg",
	} { // forloop start
		// createing a new form-data header
		// with the provided field name and file name.
		filePart, err := w.CreateFormFile(fmt.Sprintf("file%d", i+1),
			filepath.Base(file))
		if err != nil {
			t.Fatal(err)
		}

		// >> opening/ Creating the file
		var f *os.File
		if _, err := os.Stat(file); err == nil { // if exists
			f, err = os.Open(file)
			if err != nil {
				t.Fatal(err)
			}
		} else if errors.Is(err, os.ErrNotExist) {
			f, err = os.Create(file)
			if err != nil {
				t.Fatal(err)
			}
		}

		// copying the filePart into the file
		_, err = io.Copy(filePart, f)
		_ = f.Close()
		if err != nil {
			t.Fatal(err)
		}
	}

	err := w.Close()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// creating the request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://httpbin.org/post", reqBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType()) // setting the content-type

	// sending the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d; actual status %d",
			http.StatusOK, resp.StatusCode)
	}

	t.Logf("\n%s", b) //printing the respose body
}
