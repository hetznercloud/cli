package hcloud

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestPricingClientGet(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/pricing", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.PricingGetResponse{
			Pricing: schema.Pricing{
				Currency: "EUR",
			},
		})
	})
	ctx := context.Background()

	pricing, _, err := env.Client.Pricing.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if pricing.Image.PerGBMonth.Currency != "EUR" {
		t.Errorf("unexpected currency: %v", pricing.Image.PerGBMonth.Currency)
	}
}
