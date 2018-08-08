package hcloud

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestDatacenterClient(t *testing.T) {
	t.Run("GetByID", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/datacenters/1", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(schema.DatacenterGetResponse{
				Datacenter: schema.Datacenter{
					ID: 1,
				},
			})
		})

		ctx := context.Background()
		datacenter, _, err := env.Client.Datacenter.GetByID(ctx, 1)
		if err != nil {
			t.Fatal(err)
		}
		if datacenter == nil {
			t.Fatal("no datacenter")
		}
		if datacenter.ID != 1 {
			t.Errorf("unexpected datacenter ID: %v", datacenter.ID)
		}

		t.Run("via Get", func(t *testing.T) {
			datacenter, _, err := env.Client.Datacenter.Get(ctx, "1")
			if err != nil {
				t.Fatal(err)
			}
			if datacenter == nil {
				t.Fatal("no datacenter")
			}
			if datacenter.ID != 1 {
				t.Errorf("unexpected datacenter ID: %v", datacenter.ID)
			}
		})
	})

	t.Run("GetByID (not found)", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/datacenters/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(schema.ErrorResponse{
				Error: schema.Error{
					Code: ErrorCodeNotFound,
				},
			})
		})

		ctx := context.Background()
		datacenter, _, err := env.Client.Datacenter.GetByID(ctx, 1)
		if err != nil {
			t.Fatal(err)
		}
		if datacenter != nil {
			t.Fatal("expected no datacenter")
		}
	})

	t.Run("GetByName", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/datacenters", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.RawQuery != "name=fsn1-dc8" {
				t.Fatal("missing name query")
			}
			json.NewEncoder(w).Encode(schema.DatacenterListResponse{
				Datacenters: []schema.Datacenter{
					{
						ID: 1,
					},
				},
			})
		})

		ctx := context.Background()
		datacenter, _, err := env.Client.Datacenter.GetByName(ctx, "fsn1-dc8")
		if err != nil {
			t.Fatal(err)
		}
		if datacenter == nil {
			t.Fatal("no datacenter")
		}
		if datacenter.ID != 1 {
			t.Errorf("unexpected datacenter ID: %v", datacenter.ID)
		}

		t.Run("via Get", func(t *testing.T) {
			datacenter, _, err := env.Client.Datacenter.Get(ctx, "fsn1-dc8")
			if err != nil {
				t.Fatal(err)
			}
			if datacenter == nil {
				t.Fatal("no datacenter")
			}
			if datacenter.ID != 1 {
				t.Errorf("unexpected datacenter ID: %v", datacenter.ID)
			}
		})
	})

	t.Run("GetByName (not found)", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/datacenters", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.RawQuery != "name=fsn1-dc8" {
				t.Fatal("missing name query")
			}
			json.NewEncoder(w).Encode(schema.DatacenterListResponse{
				Datacenters: []schema.Datacenter{},
			})
		})

		ctx := context.Background()
		datacenter, _, err := env.Client.Datacenter.GetByName(ctx, "fsn1-dc8")
		if err != nil {
			t.Fatal(err)
		}
		if datacenter != nil {
			t.Fatal("unexpected datacenter")
		}
	})

	t.Run("List", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/datacenters", func(w http.ResponseWriter, r *http.Request) {
			if page := r.URL.Query().Get("page"); page != "2" {
				t.Errorf("expected page 2; got %q", page)
			}
			if perPage := r.URL.Query().Get("per_page"); perPage != "50" {
				t.Errorf("expected per_page 50; got %q", perPage)
			}
			json.NewEncoder(w).Encode(schema.DatacenterListResponse{
				Datacenters: []schema.Datacenter{
					{ID: 1},
					{ID: 2},
				},
			})
		})

		opts := DatacenterListOpts{}
		opts.Page = 2
		opts.PerPage = 50

		ctx := context.Background()
		datacenters, _, err := env.Client.Datacenter.List(ctx, opts)
		if err != nil {
			t.Fatal(err)
		}
		if len(datacenters) != 2 {
			t.Fatal("expected 2 datacenters")
		}
	})

	t.Run("All", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/datacenters", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(struct {
				Datacenters []schema.Datacenter `json:"datacenters"`
				Meta        schema.Meta         `json:"meta"`
			}{
				Datacenters: []schema.Datacenter{
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
		datacenters, err := env.Client.Datacenter.All(ctx)
		if err != nil {
			t.Fatalf("Datacenter.List failed: %s", err)
		}
		if len(datacenters) != 3 {
			t.Fatalf("expected 3 datacenters; got %d", len(datacenters))
		}
		if datacenters[0].ID != 1 || datacenters[1].ID != 2 || datacenters[2].ID != 3 {
			t.Errorf("unexpected datacenters")
		}
	})
}
