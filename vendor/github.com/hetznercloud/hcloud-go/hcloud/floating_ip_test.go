package hcloud

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestFloatingIPClientGetByID(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/floating_ips/1", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.FloatingIPGetResponse{
			FloatingIP: schema.FloatingIP{
				ID:   1,
				Type: "ipv4",
				IP:   "131.232.99.1",
			},
		})
	})

	ctx := context.Background()
	floatingIP, _, err := env.Client.FloatingIP.GetByID(ctx, 1)
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

func TestFloatingIPClientGetByIDNotFound(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/floating_ips/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(schema.ErrorResponse{
			Error: schema.Error{
				Code: string(ErrorCodeNotFound),
			},
		})
	})

	ctx := context.Background()
	floatingIP, _, err := env.Client.FloatingIP.GetByID(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if floatingIP != nil {
		t.Fatal("expected no Floating IP")
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
				{ID: 1, Type: "ipv4", IP: "131.232.99.1"},
				{ID: 2, Type: "ipv4", IP: "131.232.99.1"},
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
			FloatingIP: schema.FloatingIP{ID: 1, Type: "ipv4", IP: "131.232.99.1"},
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

func TestFloatingIPClientDelete(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/floating_ips/1", func(w http.ResponseWriter, r *http.Request) {
		return
	})

	var (
		ctx        = context.Background()
		floatingIP = &FloatingIP{ID: 1}
	)
	_, err := env.Client.FloatingIP.Delete(ctx, floatingIP)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFloatingIPClientUpdate(t *testing.T) {
	var (
		ctx        = context.Background()
		floatingIP = &FloatingIP{ID: 1}
	)

	t.Run("update description", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/floating_ips/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "PUT" {
				t.Error("expected PUT")
			}
			var reqBody schema.FloatingIPUpdateRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Description != "test" {
				t.Errorf("unexpected description: %v", reqBody.Description)
			}
			json.NewEncoder(w).Encode(schema.FloatingIPUpdateResponse{
				FloatingIP: schema.FloatingIP{
					ID: 1,
				},
			})
		})

		opts := FloatingIPUpdateOpts{
			Description: "test",
		}
		updatedFloatingIP, _, err := env.Client.FloatingIP.Update(ctx, floatingIP, opts)
		if err != nil {
			t.Fatal(err)
		}

		if updatedFloatingIP.ID != 1 {
			t.Errorf("unexpected Floating IP ID: %v", updatedFloatingIP.ID)
		}
	})

	t.Run("no updates", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/floating_ips/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "PUT" {
				t.Error("expected PUT")
			}
			var reqBody schema.FloatingIPUpdateRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Description != "" {
				t.Errorf("unexpected no description, but got: %v", reqBody.Description)
			}
			json.NewEncoder(w).Encode(schema.FloatingIPUpdateResponse{
				FloatingIP: schema.FloatingIP{
					ID: 1,
				},
			})
		})

		opts := FloatingIPUpdateOpts{}
		updatedFloatingIP, _, err := env.Client.FloatingIP.Update(ctx, floatingIP, opts)
		if err != nil {
			t.Fatal(err)
		}

		if updatedFloatingIP.ID != 1 {
			t.Errorf("unexpected Floating IP ID: %v", updatedFloatingIP.ID)
		}
	})
}

func TestFloatingIPClientAssign(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/floating_ips/1/actions/assign", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.FloatingIPActionAssignRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Server != 1 {
			t.Errorf("unexpected server ID: %d", reqBody.Server)
		}
		json.NewEncoder(w).Encode(schema.FloatingIPActionAssignResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	var (
		ctx        = context.Background()
		floatingIP = &FloatingIP{ID: 1}
		server     = &Server{ID: 1}
	)
	action, _, err := env.Client.FloatingIP.Assign(ctx, floatingIP, server)
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestFloatingIPClientUnassign(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/floating_ips/1/actions/unassign", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.FloatingIPActionAssignResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	var (
		ctx        = context.Background()
		floatingIP = &FloatingIP{ID: 1}
	)
	action, _, err := env.Client.FloatingIP.Unassign(ctx, floatingIP)
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestFloatingIPClientChangeDNSPtr(t *testing.T) {
	var (
		ctx        = context.Background()
		floatingIP = &FloatingIP{ID: 1}
	)

	t.Run("set", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/floating_ips/1/actions/change_dns_ptr", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.FloatingIPActionChangeDNSPtrRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.IP != "127.0.0.1" {
				t.Errorf("unexpected IP: %v", reqBody.IP)
			}
			if reqBody.DNSPtr == nil || *reqBody.DNSPtr != "example.com" {
				t.Errorf("unexpected DNS ptr: %v", reqBody.DNSPtr)
			}
			json.NewEncoder(w).Encode(schema.FloatingIPActionChangeDNSPtrResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		action, _, err := env.Client.FloatingIP.ChangeDNSPtr(ctx, floatingIP, "127.0.0.1", String("example.com"))
		if err != nil {
			t.Fatal(err)
		}
		if action.ID != 1 {
			t.Errorf("unexpected action ID: %d", action.ID)
		}
	})

	t.Run("reset", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/floating_ips/1/actions/change_dns_ptr", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.FloatingIPActionChangeDNSPtrRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.IP != "127.0.0.1" {
				t.Errorf("unexpected IP: %v", reqBody.IP)
			}
			if reqBody.DNSPtr != nil {
				t.Errorf("unexpected DNS ptr: %v", reqBody.DNSPtr)
			}
			json.NewEncoder(w).Encode(schema.FloatingIPActionChangeDNSPtrResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		action, _, err := env.Client.FloatingIP.ChangeDNSPtr(ctx, floatingIP, "127.0.0.1", nil)
		if err != nil {
			t.Fatal(err)
		}
		if action.ID != 1 {
			t.Errorf("unexpected action ID: %d", action.ID)
		}
	})
}

func TestFloatingIPClientChangeProtection(t *testing.T) {
	var (
		ctx        = context.Background()
		floatingIP = &FloatingIP{ID: 1}
	)

	t.Run("enable delete protection", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/floating_ips/1/actions/change_protection", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Error("expected POST")
			}
			var reqBody schema.FloatingIPActionChangeProtectionRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Delete == nil || *reqBody.Delete != true {
				t.Errorf("unexpected delete: %v", reqBody.Delete)
			}
			json.NewEncoder(w).Encode(schema.FloatingIPActionChangeProtectionResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		opts := FloatingIPChangeProtectionOpts{
			Delete: Bool(true),
		}
		action, _, err := env.Client.FloatingIP.ChangeProtection(ctx, floatingIP, opts)
		if err != nil {
			t.Fatal(err)
		}

		if action.ID != 1 {
			t.Errorf("unexpected action ID: %v", action.ID)
		}
	})
}
