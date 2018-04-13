package hcloud

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestServerClientGetByID(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerGetResponse{
			Server: schema.Server{
				ID: 1,
			},
		})
	})
	ctx := context.Background()

	server, _, err := env.Client.Server.GetByID(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if server == nil {
		t.Fatal("no server")
	}
	if server.ID != 1 {
		t.Errorf("unexpected server ID: %v", server.ID)
	}

	t.Run("called via Get", func(t *testing.T) {
		server, _, err := env.Client.Server.Get(ctx, "1")
		if err != nil {
			t.Fatal(err)
		}
		if server == nil {
			t.Fatal("no server")
		}
		if server.ID != 1 {
			t.Errorf("unexpected server ID: %v", server.ID)
		}
	})
}

func TestServerClientGetByIDNotFound(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(schema.ErrorResponse{
			Error: schema.Error{
				Code: ErrorCodeNotFound,
			},
		})
	})

	ctx := context.Background()
	server, _, err := env.Client.Server.GetByID(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if server != nil {
		t.Fatal("expected no server")
	}
}

func TestServerClientGetByName(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "name=myserver" {
			t.Fatal("missing name query")
		}
		json.NewEncoder(w).Encode(schema.ServerListResponse{
			Servers: []schema.Server{
				{
					ID:   1,
					Name: "myserver",
				},
			},
		})
	})
	ctx := context.Background()

	server, _, err := env.Client.Server.GetByName(ctx, "myserver")
	if err != nil {
		t.Fatal(err)
	}
	if server == nil {
		t.Fatal("no server")
	}
	if server.ID != 1 {
		t.Errorf("unexpected server ID: %v", server.ID)
	}

	t.Run("via Get", func(t *testing.T) {
		server, _, err := env.Client.Server.Get(ctx, "myserver")
		if err != nil {
			t.Fatal(err)
		}
		if server == nil {
			t.Fatal("no server")
		}
		if server.ID != 1 {
			t.Errorf("unexpected server ID: %v", server.ID)
		}
	})
}

func TestServerClientGetByNameNotFound(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "name=myserver" {
			t.Fatal("missing name query")
		}
		json.NewEncoder(w).Encode(schema.ServerListResponse{
			Servers: []schema.Server{},
		})
	})

	ctx := context.Background()
	server, _, err := env.Client.Server.GetByName(ctx, "myserver")
	if err != nil {
		t.Fatal(err)
	}
	if server != nil {
		t.Fatal("unexpected server")
	}
}

func TestServersList(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		if page := r.URL.Query().Get("page"); page != "2" {
			t.Errorf("expected page 2; got %q", page)
		}
		if perPage := r.URL.Query().Get("per_page"); perPage != "50" {
			t.Errorf("expected per_page 50; got %q", perPage)
		}
		json.NewEncoder(w).Encode(schema.ServerListResponse{
			Servers: []schema.Server{
				{ID: 1},
				{ID: 2},
			},
		})
	})

	opts := ServerListOpts{}
	opts.Page = 2
	opts.PerPage = 50

	ctx := context.Background()
	servers, _, err := env.Client.Server.List(ctx, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(servers) != 2 {
		t.Fatal("expected 2 servers")
	}
}

func TestServersAll(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			Servers []schema.Server `json:"servers"`
			Meta    schema.Meta     `json:"meta"`
		}{
			Servers: []schema.Server{
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
	servers, err := env.Client.Server.All(ctx)
	if err != nil {
		t.Fatalf("Servers.List failed: %s", err)
	}
	if len(servers) != 3 {
		t.Fatalf("expected 3 servers; got %d", len(servers))
	}
	if servers[0].ID != 1 || servers[1].ID != 2 || servers[2].ID != 3 {
		t.Errorf("unexpected servers")
	}
}

func TestServersCreateWithSSHKeys(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if len(reqBody.SSHKeys) != 2 || reqBody.SSHKeys[0] != 1 || reqBody.SSHKeys[1] != 2 {
			t.Errorf("unexpected SSH keys: %v", reqBody.SSHKeys)
		}
		json.NewEncoder(w).Encode(schema.ServerCreateResponse{
			Server: schema.Server{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:       "test",
		ServerType: &ServerType{ID: 1},
		Image:      &Image{ID: 2},
		SSHKeys: []*SSHKey{
			{ID: 1},
			{ID: 2},
		},
	})
	if err != nil {
		t.Fatalf("Server.Create failed: %s", err)
	}
	if result.Server == nil {
		t.Fatal("no server")
	}
	if result.Server.ID != 1 {
		t.Errorf("unexpected server ID: %v", result.Server.ID)
	}
	if result.RootPassword != "" {
		t.Errorf("expected no root password, got: %v", result.RootPassword)
	}
}

func TestServersCreateWithoutSSHKeys(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if len(reqBody.SSHKeys) != 0 {
			t.Errorf("expected no SSH keys, but got %v", reqBody.SSHKeys)
		}
		json.NewEncoder(w).Encode(schema.ServerCreateResponse{
			Server: schema.Server{
				ID: 1,
			},
			RootPassword: String("test"),
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:       "test",
		ServerType: &ServerType{ID: 1},
		Image:      &Image{ID: 2},
	})
	if err != nil {
		t.Fatalf("Server.Create failed: %s", err)
	}
	if result.Server == nil {
		t.Fatal("no server")
	}
	if result.Server.ID != 1 {
		t.Errorf("unexpected server ID: %v", result.Server.ID)
	}
	if result.RootPassword != "test" {
		t.Errorf("unexpected root password: %v", result.RootPassword)
	}
}

func TestServersCreateWithDatacenterID(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Datacenter != "1" {
			t.Errorf("unexpected datacenter: %v", reqBody.Datacenter)
		}
		json.NewEncoder(w).Encode(schema.ServerCreateResponse{
			Server: schema.Server{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:       "test",
		ServerType: &ServerType{ID: 1},
		Image:      &Image{ID: 2},
		Datacenter: &Datacenter{ID: 1},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Server == nil {
		t.Fatal("no server")
	}
}

func TestServersCreateWithDatacenterName(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Datacenter != "dc1" {
			t.Errorf("unexpected datacenter: %v", reqBody.Datacenter)
		}
		json.NewEncoder(w).Encode(schema.ServerCreateResponse{
			Server: schema.Server{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:       "test",
		ServerType: &ServerType{ID: 1},
		Image:      &Image{ID: 2},
		Datacenter: &Datacenter{Name: "dc1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Server == nil {
		t.Fatal("no server")
	}
}

func TestServersCreateWithLocationID(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Location != "1" {
			t.Errorf("unexpected location: %v", reqBody.Location)
		}
		json.NewEncoder(w).Encode(schema.ServerCreateResponse{
			Server: schema.Server{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:       "test",
		ServerType: &ServerType{ID: 1},
		Image:      &Image{ID: 2},
		Location:   &Location{ID: 1},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Server == nil {
		t.Fatal("no server")
	}
}

func TestServersCreateWithLocationName(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Location != "loc1" {
			t.Errorf("unexpected location: %v", reqBody.Location)
		}
		json.NewEncoder(w).Encode(schema.ServerCreateResponse{
			Server: schema.Server{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:       "test",
		ServerType: &ServerType{ID: 1},
		Image:      &Image{ID: 2},
		Location:   &Location{Name: "loc1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Server == nil {
		t.Fatal("no server")
	}
}

func TestServersCreateWithUserData(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.UserData != "---user data---" {
			t.Errorf("unexpected userdata: %v", reqBody.UserData)
		}
		json.NewEncoder(w).Encode(schema.ServerCreateResponse{
			Server: schema.Server{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:       "test",
		ServerType: &ServerType{ID: 1},
		Image:      &Image{ID: 2},
		UserData:   "---user data---",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Server == nil {
		t.Fatal("no server")
	}
}

func TestServersDelete(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1", func(w http.ResponseWriter, r *http.Request) {
		return
	})

	var (
		ctx    = context.Background()
		server = &Server{ID: 1}
	)
	_, err := env.Client.Server.Delete(ctx, server)
	if err != nil {
		t.Fatalf("Server.Delete failed: %s", err)
	}
}

func TestServerClientUpdate(t *testing.T) {
	var (
		ctx    = context.Background()
		server = &Server{ID: 1}
	)

	t.Run("update name", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/servers/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "PUT" {
				t.Error("expected PUT")
			}
			var reqBody schema.ServerUpdateRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Name != "test" {
				t.Errorf("unexpected name: %v", reqBody.Name)
			}
			json.NewEncoder(w).Encode(schema.ServerUpdateResponse{
				Server: schema.Server{
					ID: 1,
				},
			})
		})

		opts := ServerUpdateOpts{
			Name: "test",
		}
		updatedServer, _, err := env.Client.Server.Update(ctx, server, opts)
		if err != nil {
			t.Fatal(err)
		}

		if updatedServer.ID != 1 {
			t.Errorf("unexpected server ID: %v", updatedServer.ID)
		}
	})

	t.Run("no updates", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/servers/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "PUT" {
				t.Error("expected PUT")
			}
			var reqBody schema.ServerUpdateRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Name != "" {
				t.Errorf("unexpected no name, but got: %v", reqBody.Name)
			}
			json.NewEncoder(w).Encode(schema.ServerUpdateResponse{
				Server: schema.Server{
					ID: 1,
				},
			})
		})

		opts := ServerUpdateOpts{}
		updatedServer, _, err := env.Client.Server.Update(ctx, server, opts)
		if err != nil {
			t.Fatal(err)
		}

		if updatedServer.ID != 1 {
			t.Errorf("unexpected server ID: %v", updatedServer.ID)
		}
	})
}

func TestServerClientPoweron(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/poweron", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionPoweronResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.Server.Poweron(ctx, &Server{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestServerClientReboot(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/reboot", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionRebootResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.Server.Reboot(ctx, &Server{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestServerClientReset(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/reset", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionResetResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.Server.Reset(ctx, &Server{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestServerClientShutdown(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/shutdown", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionShutdownResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.Server.Shutdown(ctx, &Server{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestServerClientPoweroff(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/poweroff", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionPoweroffResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.Server.Poweroff(ctx, &Server{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestServerClientResetPassword(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/reset_password", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionResetPasswordResponse{
			Action: schema.Action{
				ID: 1,
			},
			RootPassword: "secret",
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.ResetPassword(ctx, &Server{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if result.Action.ID != 1 {
		t.Errorf("unexpected action ID: %d", result.Action.ID)
	}
	if result.RootPassword != "secret" {
		t.Errorf("unexpected root password: %v", result.RootPassword)
	}
}

func TestServerClientCreateImageNoOptions(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/create_image", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionCreateImageResponse{
			Action: schema.Action{
				ID: 1,
			},
			Image: schema.Image{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.CreateImage(ctx, &Server{ID: 1}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if result.Action.ID != 1 {
		t.Errorf("unexpected action ID: %d", result.Action.ID)
	}
	if result.Image.ID != 1 {
		t.Errorf("unexpected image ID: %d", result.Image.ID)
	}
}

func TestServerClientCreateImageWithOptions(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/create_image", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerActionCreateImageRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Type == nil || *reqBody.Type != "backup" {
			t.Errorf("unexpected type: %v", reqBody.Type)
		}
		if reqBody.Description == nil || *reqBody.Description != "my backup" {
			t.Errorf("unexpected description: %v", reqBody.Description)
		}
		json.NewEncoder(w).Encode(schema.ServerActionCreateImageResponse{
			Action: schema.Action{
				ID: 1,
			},
			Image: schema.Image{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	opts := &ServerCreateImageOpts{
		Type:        ImageTypeBackup,
		Description: String("my backup"),
	}
	result, _, err := env.Client.Server.CreateImage(ctx, &Server{ID: 1}, opts)
	if err != nil {
		t.Fatal(err)
	}
	if result.Action.ID != 1 {
		t.Errorf("unexpected action ID: %d", result.Action.ID)
	}
	if result.Image.ID != 1 {
		t.Errorf("unexpected image ID: %d", result.Image.ID)
	}
}

func TestServerClientEnableRescue(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/enable_rescue", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerActionEnableRescueRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Type == nil || *reqBody.Type != "linux64" {
			t.Errorf("unexpected type: %v", reqBody.Type)
		}
		if len(reqBody.SSHKeys) != 2 || reqBody.SSHKeys[0] != 1 || reqBody.SSHKeys[1] != 2 {
			t.Errorf("unexpected SSH keys: %v", reqBody.SSHKeys)
		}
		json.NewEncoder(w).Encode(schema.ServerActionEnableRescueResponse{
			Action: schema.Action{
				ID: 1,
			},
			RootPassword: "test",
		})
	})

	ctx := context.Background()
	opts := ServerEnableRescueOpts{
		Type: ServerRescueTypeLinux64,
		SSHKeys: []*SSHKey{
			{ID: 1},
			{ID: 2},
		},
	}
	result, _, err := env.Client.Server.EnableRescue(ctx, &Server{ID: 1}, opts)
	if err != nil {
		t.Fatal(err)
	}
	if result.Action.ID != 1 {
		t.Errorf("unexpected action ID: %d", result.Action.ID)
	}
	if result.RootPassword != "test" {
		t.Errorf("unexpected root password: %s", result.RootPassword)
	}
}

func TestServerClientDisableRescue(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/disable_rescue", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionDisableRescueResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.Server.DisableRescue(ctx, &Server{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestServerClientRebuild(t *testing.T) {
	var (
		ctx    = context.Background()
		server = &Server{ID: 1}
	)

	t.Run("with image ID", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/servers/1/actions/rebuild", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.ServerActionRebuildRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if id, ok := reqBody.Image.(float64); !ok || id != 1 {
				t.Errorf("unexpected image ID: %v", reqBody.Image)
			}
			json.NewEncoder(w).Encode(schema.ServerActionRebuildResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		opts := ServerRebuildOpts{
			Image: &Image{ID: 1},
		}
		action, _, err := env.Client.Server.Rebuild(ctx, server, opts)
		if err != nil {
			t.Fatal(err)
		}
		if action.ID != 1 {
			t.Errorf("unexpected action ID: %d", action.ID)
		}
	})

	t.Run("with image name", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/servers/1/actions/rebuild", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.ServerActionRebuildRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if name, ok := reqBody.Image.(string); !ok || name != "debian-9" {
				t.Errorf("unexpected image name: %v", reqBody.Image)
			}
			json.NewEncoder(w).Encode(schema.ServerActionRebuildResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		opts := ServerRebuildOpts{
			Image: &Image{Name: "debian-9"},
		}
		action, _, err := env.Client.Server.Rebuild(ctx, server, opts)
		if err != nil {
			t.Fatal(err)
		}
		if action.ID != 1 {
			t.Errorf("unexpected action ID: %d", action.ID)
		}
	})
}

func TestServerClientAttachISO(t *testing.T) {
	var (
		ctx    = context.Background()
		server = &Server{ID: 1}
	)

	t.Run("with ISO ID", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/servers/1/actions/attach_iso", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.ServerActionAttachISORequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if id, ok := reqBody.ISO.(float64); !ok || id != 1 {
				t.Errorf("unexpected ISO ID: %v", reqBody.ISO)
			}
			json.NewEncoder(w).Encode(schema.ServerActionAttachISOResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		iso := &ISO{ID: 1}
		action, _, err := env.Client.Server.AttachISO(ctx, server, iso)
		if err != nil {
			t.Fatal(err)
		}
		if action.ID != 1 {
			t.Errorf("unexpected action ID: %d", action.ID)
		}
	})

	t.Run("with ISO name", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/servers/1/actions/attach_iso", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.ServerActionAttachISORequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if name, ok := reqBody.ISO.(string); !ok || name != "debian.iso" {
				t.Errorf("unexpected ISO name: %v", reqBody.ISO)
			}
			json.NewEncoder(w).Encode(schema.ServerActionAttachISOResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		iso := &ISO{Name: "debian.iso"}
		action, _, err := env.Client.Server.AttachISO(ctx, server, iso)
		if err != nil {
			t.Fatal(err)
		}
		if action.ID != 1 {
			t.Errorf("unexpected action ID: %d", action.ID)
		}
	})
}

func TestServerClientDetachISO(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	var (
		ctx    = context.Background()
		server = &Server{ID: 1}
	)

	env.Mux.HandleFunc("/servers/1/actions/detach_iso", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionDetachISOResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	action, _, err := env.Client.Server.DetachISO(ctx, server)
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestServerClientEnableBackup(t *testing.T) {
	var (
		ctx    = context.Background()
		server = &Server{ID: 1}
	)

	t.Run("with a backup window", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/servers/1/actions/enable_backup", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.ServerActionEnableBackupRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.BackupWindow == nil || *reqBody.BackupWindow != "9-17" {
				t.Errorf("unexpected backup window: %v", reqBody.BackupWindow)
			}
			json.NewEncoder(w).Encode(schema.ServerActionEnableBackupResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		action, _, err := env.Client.Server.EnableBackup(ctx, server, "9-17")
		if err != nil {
			t.Fatal(err)
		}
		if action.ID != 1 {
			t.Errorf("unexpected action ID: %d", action.ID)
		}
	})

	t.Run("without a backup window", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/servers/1/actions/enable_backup", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.ServerActionEnableBackupRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.BackupWindow != nil {
				t.Errorf("unexpected backup window: %v", reqBody.BackupWindow)
			}
			json.NewEncoder(w).Encode(schema.ServerActionEnableBackupResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		action, _, err := env.Client.Server.EnableBackup(ctx, server, "")
		if err != nil {
			t.Fatal(err)
		}
		if action.ID != 1 {
			t.Errorf("unexpected action ID: %d", action.ID)
		}
	})
}

func TestServerClientDisableBackup(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	var (
		ctx    = context.Background()
		server = &Server{ID: 1}
	)

	env.Mux.HandleFunc("/servers/1/actions/disable_backup", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionDisableBackupResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	action, _, err := env.Client.Server.DisableBackup(ctx, server)
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestServerClientChangeType(t *testing.T) {
	var (
		ctx    = context.Background()
		server = &Server{ID: 1}
	)

	t.Run("with server type ID", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/servers/1/actions/change_type", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.ServerActionChangeTypeRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if id, ok := reqBody.ServerType.(float64); !ok || id != 1 {
				t.Errorf("unexpected server type ID: %v", reqBody.ServerType)
			}
			if !reqBody.UpgradeDisk {
				t.Error("expected to upgrade disk")
			}
			json.NewEncoder(w).Encode(schema.ServerActionChangeTypeResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		opts := ServerChangeTypeOpts{
			ServerType:  &ServerType{ID: 1},
			UpgradeDisk: true,
		}
		action, _, err := env.Client.Server.ChangeType(ctx, server, opts)
		if err != nil {
			t.Fatal(err)
		}
		if action.ID != 1 {
			t.Errorf("unexpected action ID: %d", action.ID)
		}
	})

	t.Run("with server type name", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/servers/1/actions/change_type", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.ServerActionChangeTypeRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if name, ok := reqBody.ServerType.(string); !ok || name != "type" {
				t.Errorf("unexpected server type name: %v", reqBody.ServerType)
			}
			if !reqBody.UpgradeDisk {
				t.Error("expected to upgrade disk")
			}
			json.NewEncoder(w).Encode(schema.ServerActionChangeTypeResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		opts := ServerChangeTypeOpts{
			ServerType:  &ServerType{Name: "type"},
			UpgradeDisk: true,
		}
		action, _, err := env.Client.Server.ChangeType(ctx, server, opts)
		if err != nil {
			t.Fatal(err)
		}
		if action.ID != 1 {
			t.Errorf("unexpected action ID: %d", action.ID)
		}
	})
}

func TestServerClientChangeDNSPtr(t *testing.T) {
	var (
		ctx    = context.Background()
		server = &Server{ID: 1}
	)

	t.Run("set", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/servers/1/actions/change_dns_ptr", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.ServerActionChangeDNSPtrRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.IP != "127.0.0.1" {
				t.Errorf("unexpected IP: %v", reqBody.IP)
			}
			if reqBody.DNSPtr == nil || *reqBody.DNSPtr != "example.com" {
				t.Errorf("unexpected DNS ptr: %v", reqBody.DNSPtr)
			}
			json.NewEncoder(w).Encode(schema.ServerActionChangeDNSPtrResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		action, _, err := env.Client.Server.ChangeDNSPtr(ctx, server, "127.0.0.1", String("example.com"))
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

		env.Mux.HandleFunc("/servers/1/actions/change_dns_ptr", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.ServerActionChangeDNSPtrRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.IP != "127.0.0.1" {
				t.Errorf("unexpected IP: %v", reqBody.IP)
			}
			if reqBody.DNSPtr != nil {
				t.Errorf("unexpected DNS ptr: %v", reqBody.DNSPtr)
			}
			json.NewEncoder(w).Encode(schema.ServerActionChangeDNSPtrResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		action, _, err := env.Client.Server.ChangeDNSPtr(ctx, server, "127.0.0.1", nil)
		if err != nil {
			t.Fatal(err)
		}
		if action.ID != 1 {
			t.Errorf("unexpected action ID: %d", action.ID)
		}
	})
}

func TestServerClientChangeProtection(t *testing.T) {
	var (
		ctx    = context.Background()
		server = &Server{ID: 1}
	)

	t.Run("enable delete and rebuild protection", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/servers/1/actions/change_protection", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Error("expected POST")
			}
			var reqBody schema.ServerActionChangeProtectionRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Delete == nil || *reqBody.Delete != true {
				t.Errorf("unexpected delete: %v", reqBody.Delete)
			}
			if reqBody.Rebuild == nil || *reqBody.Rebuild != true {
				t.Errorf("unexpected rebuild: %v", reqBody.Rebuild)
			}
			json.NewEncoder(w).Encode(schema.ImageActionChangeProtectionResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		opts := ServerChangeProtectionOpts{
			Delete:  Bool(true),
			Rebuild: Bool(true),
		}
		action, _, err := env.Client.Server.ChangeProtection(ctx, server, opts)
		if err != nil {
			t.Fatal(err)
		}

		if action.ID != 1 {
			t.Errorf("unexpected action ID: %v", action.ID)
		}
	})
}
