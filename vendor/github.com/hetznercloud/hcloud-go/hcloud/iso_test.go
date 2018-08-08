package hcloud

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestISOIsDeprecated(t *testing.T) {
	t.Run("not deprecated", func(t *testing.T) {
		iso := &ISO{}
		if iso.IsDeprecated() {
			t.Errorf("unexpected value for IsDeprecated: %v", iso.IsDeprecated())
		}
	})

	t.Run("deprecated", func(t *testing.T) {
		iso := &ISO{
			Deprecated: time.Now(),
		}
		if !iso.IsDeprecated() {
			t.Errorf("unexpected value for IsDeprecated: %v", iso.IsDeprecated())
		}
	})
}

func TestISOClient(t *testing.T) {
	t.Run("GetByID", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/isos/1", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(schema.ISOGetResponse{
				ISO: schema.ISO{
					ID: 1,
				},
			})
		})

		ctx := context.Background()
		iso, _, err := env.Client.ISO.GetByID(ctx, 1)
		if err != nil {
			t.Fatal(err)
		}
		if iso == nil {
			t.Fatal("no iso")
		}
		if iso.ID != 1 {
			t.Errorf("unexpected iso ID: %v", iso.ID)
		}

		t.Run("via Get", func(t *testing.T) {
			iso, _, err := env.Client.ISO.Get(ctx, "1")
			if err != nil {
				t.Fatal(err)
			}
			if iso == nil {
				t.Fatal("no iso")
			}
			if iso.ID != 1 {
				t.Errorf("unexpected iso ID: %v", iso.ID)
			}
		})
	})

	t.Run("GetByID (not found)", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/isos/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(schema.ErrorResponse{
				Error: schema.Error{
					Code: string(ErrorCodeNotFound),
				},
			})
		})

		ctx := context.Background()
		iso, _, err := env.Client.ISO.GetByID(ctx, 1)
		if err != nil {
			t.Fatal(err)
		}
		if iso != nil {
			t.Fatal("expected no iso")
		}
	})

	t.Run("GetByName", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/isos", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.RawQuery != "name=debian-9" {
				t.Fatal("missing name query")
			}
			json.NewEncoder(w).Encode(schema.ISOListResponse{
				ISOs: []schema.ISO{
					{
						ID: 1,
					},
				},
			})
		})

		ctx := context.Background()
		iso, _, err := env.Client.ISO.GetByName(ctx, "debian-9")
		if err != nil {
			t.Fatal(err)
		}
		if iso == nil {
			t.Fatal("no iso")
		}
		if iso.ID != 1 {
			t.Errorf("unexpected iso ID: %v", iso.ID)
		}

		t.Run("via Get", func(t *testing.T) {
			iso, _, err := env.Client.ISO.Get(ctx, "debian-9")
			if err != nil {
				t.Fatal(err)
			}
			if iso == nil {
				t.Fatal("no iso")
			}
			if iso.ID != 1 {
				t.Errorf("unexpected iso ID: %v", iso.ID)
			}
		})
	})

	t.Run("GetByName (not found)", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/isos", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.RawQuery != "name=debian-9" {
				t.Fatal("missing name query")
			}
			json.NewEncoder(w).Encode(schema.ISOListResponse{
				ISOs: []schema.ISO{},
			})
		})

		ctx := context.Background()
		iso, _, err := env.Client.ISO.GetByName(ctx, "debian-9")
		if err != nil {
			t.Fatal(err)
		}
		if iso != nil {
			t.Fatal("unexpected iso")
		}
	})

	t.Run("List", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/isos", func(w http.ResponseWriter, r *http.Request) {
			if page := r.URL.Query().Get("page"); page != "2" {
				t.Errorf("expected page 2; got %q", page)
			}
			if perPage := r.URL.Query().Get("per_page"); perPage != "50" {
				t.Errorf("expected per_page 50; got %q", perPage)
			}
			json.NewEncoder(w).Encode(schema.ISOListResponse{
				ISOs: []schema.ISO{
					{ID: 1},
					{ID: 2},
				},
			})
		})

		opts := ISOListOpts{}
		opts.Page = 2
		opts.PerPage = 50

		ctx := context.Background()
		isos, _, err := env.Client.ISO.List(ctx, opts)
		if err != nil {
			t.Fatal(err)
		}
		if len(isos) != 2 {
			t.Fatal("expected 2 isos")
		}
	})

	t.Run("All", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/isos", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(struct {
				ISOs []schema.ISO `json:"isos"`
				Meta schema.Meta  `json:"meta"`
			}{
				ISOs: []schema.ISO{
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
		isos, err := env.Client.ISO.All(ctx)
		if err != nil {
			t.Fatalf("ISO.List failed: %s", err)
		}
		if len(isos) != 3 {
			t.Fatalf("expected 3 isos; got %d", len(isos))
		}
		if isos[0].ID != 1 || isos[1].ID != 2 || isos[2].ID != 3 {
			t.Errorf("unexpected isos")
		}
	})
}
