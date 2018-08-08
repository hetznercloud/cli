package hcloud

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestServerTypeClient(t *testing.T) {
	t.Run("GetByID", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/server_types/1", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(schema.ServerTypeGetResponse{
				ServerType: schema.ServerType{
					ID: 1,
				},
			})
		})

		ctx := context.Background()
		serverType, _, err := env.Client.ServerType.GetByID(ctx, 1)
		if err != nil {
			t.Fatal(err)
		}
		if serverType == nil {
			t.Fatal("no server type")
		}
		if serverType.ID != 1 {
			t.Errorf("unexpected server type ID: %v", serverType.ID)
		}

		t.Run("via Get", func(t *testing.T) {
			serverType, _, err := env.Client.ServerType.Get(ctx, "1")
			if err != nil {
				t.Fatal(err)
			}
			if serverType == nil {
				t.Fatal("no server type")
			}
			if serverType.ID != 1 {
				t.Errorf("unexpected server type ID: %v", serverType.ID)
			}
		})
	})

	t.Run("GetByID (not found)", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/server_types/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(schema.ErrorResponse{
				Error: schema.Error{
					Code: string(ErrorCodeNotFound),
				},
			})
		})

		ctx := context.Background()
		serverType, _, err := env.Client.ServerType.GetByID(ctx, 1)
		if err != nil {
			t.Fatal(err)
		}
		if serverType != nil {
			t.Fatal("expected no server type")
		}
	})

	t.Run("GetByName", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/server_types", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.RawQuery != "name=cx10" {
				t.Fatal("missing name query")
			}
			json.NewEncoder(w).Encode(schema.ServerTypeListResponse{
				ServerTypes: []schema.ServerType{
					{
						ID: 1,
					},
				},
			})
		})

		ctx := context.Background()
		serverType, _, err := env.Client.ServerType.GetByName(ctx, "cx10")
		if err != nil {
			t.Fatal(err)
		}
		if serverType == nil {
			t.Fatal("no server type")
		}
		if serverType.ID != 1 {
			t.Errorf("unexpected server type ID: %v", serverType.ID)
		}

		t.Run("via Get", func(t *testing.T) {
			serverType, _, err := env.Client.ServerType.Get(ctx, "cx10")
			if err != nil {
				t.Fatal(err)
			}
			if serverType == nil {
				t.Fatal("no serverType")
			}
			if serverType.ID != 1 {
				t.Errorf("unexpected server type ID: %v", serverType.ID)
			}
		})
	})

	t.Run("GetByName (not found)", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/server_types", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.RawQuery != "name=cx10" {
				t.Fatal("missing name query")
			}
			json.NewEncoder(w).Encode(schema.ServerTypeListResponse{
				ServerTypes: []schema.ServerType{},
			})
		})

		ctx := context.Background()
		serverType, _, err := env.Client.ServerType.GetByName(ctx, "cx10")
		if err != nil {
			t.Fatal(err)
		}
		if serverType != nil {
			t.Fatal("unexpected server type")
		}
	})

	t.Run("List", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/server_types", func(w http.ResponseWriter, r *http.Request) {
			if page := r.URL.Query().Get("page"); page != "2" {
				t.Errorf("expected page 2; got %q", page)
			}
			if perPage := r.URL.Query().Get("per_page"); perPage != "50" {
				t.Errorf("expected per_page 50; got %q", perPage)
			}
			json.NewEncoder(w).Encode(schema.ServerTypeListResponse{
				ServerTypes: []schema.ServerType{
					{ID: 1},
					{ID: 2},
				},
			})
		})

		opts := ServerTypeListOpts{}
		opts.Page = 2
		opts.PerPage = 50

		ctx := context.Background()
		serverTypes, _, err := env.Client.ServerType.List(ctx, opts)
		if err != nil {
			t.Fatal(err)
		}
		if len(serverTypes) != 2 {
			t.Fatal("expected 2 server types")
		}
	})

	t.Run("All", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/server_types", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(struct {
				ServerTypes []schema.ServerType `json:"server_types"`
				Meta        schema.Meta         `json:"meta"`
			}{
				ServerTypes: []schema.ServerType{
					{ID: 1},
					{ID: 2},
					{ID: 3},
				},
				Meta: schema.Meta{
					Pagination: &schema.MetaPagination{
						Page:         1,
						LastPage:     3,
						PerPage:      3,
						TotalEntries: 3,
					},
				},
			})
		})

		ctx := context.Background()
		serverTypes, err := env.Client.ServerType.All(ctx)
		if err != nil {
			t.Fatalf("ServerTypes.List failed: %s", err)
		}
		if len(serverTypes) != 3 {
			t.Fatalf("expected 3 server types; got %d", len(serverTypes))
		}
		if serverTypes[0].ID != 1 || serverTypes[1].ID != 2 || serverTypes[2].ID != 3 {
			t.Errorf("unexpected server types")
		}
	})
}
