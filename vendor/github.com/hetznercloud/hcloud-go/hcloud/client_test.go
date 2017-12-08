package hcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

type testEnv struct {
	Server *httptest.Server
	Mux    *http.ServeMux
	Client *Client
}

func (env *testEnv) Teardown() {
	env.Server.Close()
	env.Server = nil
	env.Mux = nil
	env.Client = nil
}

func newTestEnv() testEnv {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	client := NewClient(
		WithEndpoint(server.URL),
		WithToken("token"),
		WithBackoffFunc(func(_ int) time.Duration { return 0 }),
	)
	return testEnv{
		Server: server,
		Mux:    mux,
		Client: client,
	}
}

func TestClientEndpointTrailingSlashesRemoved(t *testing.T) {
	client := NewClient(WithEndpoint("http://api/v1.0/////"))
	if strings.HasSuffix(client.endpoint, "/") {
		t.Fatalf("endpoint has trailing slashes: %q", client.endpoint)
	}
}

func TestClientError(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, `{
			"error": {
				"code": "service_error",
				"message": "An error occured"
			}
		}`)
	})

	ctx := context.Background()
	req, err := env.Client.NewRequest(ctx, "GET", "/error", nil)
	if err != nil {
		t.Fatalf("error creating request: %s", err)
	}

	_, err = env.Client.Do(req, nil)
	if _, ok := err.(Error); !ok {
		t.Fatalf("unexpected error of type %T: %v", err, err)
	}

	apiError := err.(Error)

	if apiError.Code != "service_error" {
		t.Errorf("unexpected error code: %q", apiError.Code)
	}
	if apiError.Message != "An error occured" {
		t.Errorf("unexpected error message: %q", apiError.Message)
	}
}

func TestClientMeta(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("RateLimit-Limit", "1000")
		w.Header().Set("RateLimit-Remaining", "999")
		w.Header().Set("RateLimit-Reset", "1511954577")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"foo": "bar",
			"meta": {
				"pagination": {
					"page": 1
				}
			}
		}`)
	})

	ctx := context.Background()
	req, err := env.Client.NewRequest(ctx, "GET", "/", nil)
	if err != nil {
		t.Fatalf("error creating request: %s", err)
	}

	response, err := env.Client.Do(req, nil)
	if err != nil {
		t.Fatalf("request failed: %s", err)
	}

	if response.Meta.Ratelimit.Limit != 1000 {
		t.Errorf("unexpected ratelimit limit: %d", response.Meta.Ratelimit.Limit)
	}
	if response.Meta.Ratelimit.Remaining != 999 {
		t.Errorf("unexpected ratelimit remaining: %d", response.Meta.Ratelimit.Remaining)
	}
	if !response.Meta.Ratelimit.Reset.Equal(time.Unix(1511954577, 0)) {
		t.Errorf("unexpected ratelimit reset: %v", response.Meta.Ratelimit.Reset)
	}

	if response.Meta.Pagination.Page != 1 {
		t.Error("missing pagination")
	}
}

func TestClientMetaNonJSON(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "foo")
	})

	ctx := context.Background()
	req, err := env.Client.NewRequest(ctx, "GET", "/", nil)
	if err != nil {
		t.Fatalf("error creating request: %s", err)
	}

	response, err := env.Client.Do(req, nil)
	if err != nil {
		t.Fatalf("request failed: %s", err)
	}

	if response.Meta.Pagination != nil {
		t.Fatal("pagination should not be present")
	}
}

func TestPaginationFromSchema(t *testing.T) {
	data := []byte(`{
		"page": 2,
		"per_page": 25,
		"previous_page": 1,
		"next_page": 3,
		"last_page": 13,
		"total_entries": 322
	}`)

	var s schema.MetaPagination
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	p := PaginationFromSchema(s)

	if p.Page != 2 {
		t.Errorf("unexpected page: %v", p.Page)
	}
	if p.PerPage != 25 {
		t.Errorf("unexpected per page: %v", p.PerPage)
	}
	if p.PreviousPage != 1 {
		t.Errorf("unexpected previous page: %v", p.PreviousPage)
	}
	if p.NextPage != 3 {
		t.Errorf("unexpected next page: %d", p.NextPage)
	}
	if p.LastPage != 13 {
		t.Errorf("unexpected last page: %d", p.LastPage)
	}
	if p.TotalEntries != 322 {
		t.Errorf("unexpected total entries: %d", p.TotalEntries)
	}
}
