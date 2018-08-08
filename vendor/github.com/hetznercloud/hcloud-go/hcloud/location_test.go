package hcloud

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestLocationClient(t *testing.T) {
	t.Run("GetByID", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/locations/1", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(schema.LocationGetResponse{
				Location: schema.Location{
					ID: 1,
				},
			})
		})

		ctx := context.Background()
		location, _, err := env.Client.Location.GetByID(ctx, 1)
		if err != nil {
			t.Fatal(err)
		}
		if location == nil {
			t.Fatal("no location")
		}
		if location.ID != 1 {
			t.Errorf("unexpected location ID: %v", location.ID)
		}

		t.Run("via Get", func(t *testing.T) {
			location, _, err := env.Client.Location.Get(ctx, "1")
			if err != nil {
				t.Fatal(err)
			}
			if location == nil {
				t.Fatal("no location")
			}
			if location.ID != 1 {
				t.Errorf("unexpected location ID: %v", location.ID)
			}
		})
	})

	t.Run("GetByID (not found)", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/locations/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(schema.ErrorResponse{
				Error: schema.Error{
					Code: string(ErrorCodeNotFound),
				},
			})
		})

		ctx := context.Background()
		location, _, err := env.Client.Location.GetByID(ctx, 1)
		if err != nil {
			t.Fatal(err)
		}
		if location != nil {
			t.Fatal("expected no location")
		}
	})

	t.Run("GetByName", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/locations", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.RawQuery != "name=fsn1-dc8" {
				t.Fatal("missing name query")
			}
			json.NewEncoder(w).Encode(schema.LocationListResponse{
				Locations: []schema.Location{
					{
						ID: 1,
					},
				},
			})
		})

		ctx := context.Background()
		location, _, err := env.Client.Location.GetByName(ctx, "fsn1-dc8")
		if err != nil {
			t.Fatal(err)
		}
		if location == nil {
			t.Fatal("no location")
		}
		if location.ID != 1 {
			t.Errorf("unexpected location ID: %v", location.ID)
		}

		t.Run("via Get", func(t *testing.T) {
			location, _, err := env.Client.Location.Get(ctx, "fsn1-dc8")
			if err != nil {
				t.Fatal(err)
			}
			if location == nil {
				t.Fatal("no location")
			}
			if location.ID != 1 {
				t.Errorf("unexpected location ID: %v", location.ID)
			}
		})
	})

	t.Run("GetByName (not found)", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/locations", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.RawQuery != "name=fsn1-dc8" {
				t.Fatal("missing name query")
			}
			json.NewEncoder(w).Encode(schema.LocationListResponse{
				Locations: []schema.Location{},
			})
		})

		ctx := context.Background()
		location, _, err := env.Client.Location.GetByName(ctx, "fsn1-dc8")
		if err != nil {
			t.Fatal(err)
		}
		if location != nil {
			t.Fatal("unexpected location")
		}
	})

	t.Run("List", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/locations", func(w http.ResponseWriter, r *http.Request) {
			if page := r.URL.Query().Get("page"); page != "2" {
				t.Errorf("expected page 2; got %q", page)
			}
			if perPage := r.URL.Query().Get("per_page"); perPage != "50" {
				t.Errorf("expected per_page 50; got %q", perPage)
			}
			json.NewEncoder(w).Encode(schema.LocationListResponse{
				Locations: []schema.Location{
					{ID: 1},
					{ID: 2},
				},
			})
		})

		opts := LocationListOpts{}
		opts.Page = 2
		opts.PerPage = 50

		ctx := context.Background()
		locations, _, err := env.Client.Location.List(ctx, opts)
		if err != nil {
			t.Fatal(err)
		}
		if len(locations) != 2 {
			t.Fatal("expected 2 locations")
		}
	})

	t.Run("All", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/locations", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(struct {
				Locations []schema.Location `json:"locations"`
				Meta      schema.Meta       `json:"meta"`
			}{
				Locations: []schema.Location{
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
		locations, err := env.Client.Location.All(ctx)
		if err != nil {
			t.Fatalf("Location.List failed: %s", err)
		}
		if len(locations) != 3 {
			t.Fatalf("expected 3 locations; got %d", len(locations))
		}
		if locations[0].ID != 1 || locations[1].ID != 2 || locations[2].ID != 3 {
			t.Errorf("unexpected locations")
		}
	})
}
