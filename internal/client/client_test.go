package client

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func testHTTPClient(fn func(*http.Request) (*http.Response, error)) *http.Client {
	return &http.Client{Transport: roundTripFunc(fn)}
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

func TestRequestAddsBearerToken(t *testing.T) {
	httpClient := testHTTPClient(func(r *http.Request) (*http.Response, error) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Fatalf("Authorization = %q", r.Header.Get("Authorization"))
		}
		return jsonResponse(http.StatusOK, `{"ok":true}`), nil
	})

	_, err := New("https://api.example", "test-token", httpClient).Do(context.Background(), http.MethodGet, "/team", nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRequestEncodesJSONBody(t *testing.T) {
	httpClient := testHTTPClient(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("Content-Type = %q", r.Header.Get("Content-Type"))
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != `{"name":"demo"}` {
			t.Fatalf("body = %s", body)
		}
		return jsonResponse(http.StatusOK, `{"id":1}`), nil
	})

	_, err := New("https://api.example", "token", httpClient).Do(context.Background(), http.MethodPost, "/team", []byte(`{"name":"demo"}`))
	if err != nil {
		t.Fatal(err)
	}
}

func TestRequestReturnsHTTPErrorEnvelope(t *testing.T) {
	httpClient := testHTTPClient(func(r *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusBadRequest, `{"message":"bad"}`), nil
	})

	_, err := New("https://api.example", "token", httpClient).Do(context.Background(), http.MethodGet, "/team", nil)
	var httpErr HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("error = %T %v, want HTTPError", err, err)
	}
	if httpErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("StatusCode = %d", httpErr.StatusCode)
	}
}

func TestRequestTrimsBaseURLSlash(t *testing.T) {
	httpClient := testHTTPClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != "https://api.example/team" {
			t.Fatalf("url = %q", r.URL.String())
		}
		return jsonResponse(http.StatusOK, `{"ok":true}`), nil
	})

	_, err := New("https://api.example/", "token", httpClient).Do(context.Background(), http.MethodGet, "/team", nil)
	if err != nil {
		t.Fatal(err)
	}
}
