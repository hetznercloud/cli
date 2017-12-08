package hcloud

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestFloatingIPClientGet(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/floating_ips/1", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.FloatingIPGetResponse{
			FloatingIP: schema.FloatingIP{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	floatingIP, _, err := env.Client.FloatingIP.Get(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if floatingIP == nil {
		t.Fatal("no Floating IP")
	}
	if floatingIP.ID != 1 {
		t.Errorf("unexpected ID: %v", floatingIP.ID)
	}
}

func TestFloatingIPClientList(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/floating_ips", func(w http.ResponseWriter, r *http.Request) {
		if page := r.URL.Query().Get("page"); page != "2" {
			t.Errorf("expected page 2; got %q", page)
		}
		if perPage := r.URL.Query().Get("per_page"); perPage != "50" {
			t.Errorf("expected per_page 50; got %q", perPage)
		}
		json.NewEncoder(w).Encode(schema.FloatingIPListResponse{
			FloatingIPs: []schema.FloatingIP{
				{ID: 1},
				{ID: 2},
			},
		})
	})

	opts := FloatingIPListOpts{}
	opts.Page = 2
	opts.PerPage = 50

	ctx := context.Background()
	floatingIPs, _, err := env.Client.FloatingIP.List(ctx, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(floatingIPs) != 2 {
		t.Fatal("expected 2 Floating IPs")
	}
}

func TestFloatingIPClientCreate(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/floating_ips", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("expected POST")
		}
		json.NewEncoder(w).Encode(schema.FloatingIPCreateResponse{
			FloatingIP: schema.FloatingIP{
				ID: 1,
			},
			Action: &schema.Action{
				ID: 1,
			},
		})
	})

	opts := FloatingIPCreateOpts{
		Type:         FloatingIPTypeIPv4,
		Description:  String("test"),
		HomeLocation: &Location{Name: "test"},
		Server:       &Server{ID: 1},
	}

	ctx := context.Background()
	result, _, err := env.Client.FloatingIP.Create(ctx, opts)
	if err != nil {
		t.Fatal(err)
	}

	if result.FloatingIP.ID != 1 {
		t.Errorf("unexpected Floating IP ID: %d", result.FloatingIP.ID)
	}
	if result.Action.ID != 1 {
		t.Errorf("unexpected action ID: %d", result.Action.ID)
	}
}
