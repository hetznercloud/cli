package hcloud

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestActionClientGetByID(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/actions/1", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ActionGetResponse{
			Action: schema.Action{
				ID:       1,
				Status:   "running",
				Command:  "create_server",
				Progress: 50,
				Started:  time.Date(2017, 12, 4, 14, 31, 1, 0, time.UTC),
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.Action.GetByID(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if action == nil {
		t.Fatal("no action")
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %v", action.ID)
	}
}

func TestActionClientGetByIDNotFound(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/actions/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(schema.ErrorResponse{
			Error: schema.Error{
				Code: string(ErrorCodeNotFound),
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.Action.GetByID(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if action != nil {
		t.Fatal("expected no action")
	}
}

func TestActionClientList(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/actions", func(w http.ResponseWriter, r *http.Request) {
		if page := r.URL.Query().Get("page"); page != "2" {
			t.Errorf("expected page 2; got %q", page)
		}
		if perPage := r.URL.Query().Get("per_page"); perPage != "50" {
			t.Errorf("expected per_page 50; got %q", perPage)
		}
		json.NewEncoder(w).Encode(schema.ActionListResponse{
			Actions: []schema.Action{
				{ID: 1},
				{ID: 2},
			},
		})
	})

	opts := ActionListOpts{}
	opts.Page = 2
	opts.PerPage = 50

	ctx := context.Background()
	actions, _, err := env.Client.Action.List(ctx, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(actions) != 2 {
		t.Fatal("expected 2 actions")
	}
}

func TestActionClientAll(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/actions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			Actions []schema.Action `json:"actions"`
			Meta    schema.Meta     `json:"meta"`
		}{
			Actions: []schema.Action{
				{ID: 1},
				{ID: 2},
				{ID: 3},
			},
			Meta: schema.Meta{
				Pagination: &schema.MetaPagination{
					Page:         1,
					LastPage:     1,
					PerPage:      3,
					TotalEntries: 3,
				},
			},
		})
	})

	ctx := context.Background()
	actions, err := env.Client.Action.All(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(actions) != 3 {
		t.Fatalf("expected 3 actions; got %d", len(actions))
	}
	if actions[0].ID != 1 || actions[1].ID != 2 || actions[2].ID != 3 {
		t.Errorf("unexpected actions")
	}
}

func TestActionClientWatchProgress(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	callCount := 0
	env.Mux.HandleFunc("/actions/1", func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		switch callCount {
		case 1:
			json.NewEncoder(w).Encode(schema.ActionGetResponse{
				Action: schema.Action{
					ID:       1,
					Status:   "running",
					Progress: 50,
				},
			})
		case 2:
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(schema.ErrorResponse{
				Error: schema.Error{
					Code:    string(ErrorCodeRateLimitExceeded),
					Message: "ratelimited",
				},
			})
			return
		case 3:
			json.NewEncoder(w).Encode(schema.ActionGetResponse{
				Action: schema.Action{
					ID:       1,
					Status:   "error",
					Progress: 100,
					Error: &schema.ActionError{
						Code:    "action_failed",
						Message: "action failed",
					},
				},
			})
		default:
			t.Errorf("unexpected number of calls to the test server: %v", callCount)
		}
	})
	action := &Action{
		ID:       1,
		Status:   ActionStatusRunning,
		Progress: 0,
	}

	ctx := context.Background()
	progressCh, errCh := env.Client.Action.WatchProgress(ctx, action)
	var (
		progressUpdates []int
		err             error
	)

loop:
	for {
		select {
		case progress := <-progressCh:
			progressUpdates = append(progressUpdates, progress)
		case err = <-errCh:
			break loop
		}
	}

	if err == nil {
		t.Fatal("expected an error")
	}
	if e, ok := err.(ActionError); !ok || e.Code != "action_failed" {
		t.Fatalf("expected hcloud.Error, but got: %#v", err)
	}
	if len(progressUpdates) != 1 || progressUpdates[0] != 50 {
		t.Fatalf("unexpected progress updates: %v", progressUpdates)
	}
}
