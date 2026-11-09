package state

import (
	"context"
	"net/http"
	"sync"

	"github.com/hetznercloud/cli/internal/state/config"
)

// The SDK puts this marker on requests instead of a secret. Resolving it in the
// transport keeps client construction and configuration inspection side-effect free.
// Only replace this header when it is present: net/http removes Authorization on
// redirects to unrelated hosts, and we must not restore it in that case.
const credentialMarker = "hcloud-credential-helper" //nolint:gosec // Non-secret marker replaced by the transport.

type credentialResolver struct {
	once   sync.Once
	helper []string
	url    string
	token  string
	err    error
}

func (r *credentialResolver) resolve(ctx context.Context) (string, error) {
	r.once.Do(func() {
		r.token, r.err = config.ResolveCredential(ctx, r.helper, r.url)
	})
	return r.token, r.err
}

type credentialTransport struct {
	resolver *credentialResolver
	base     http.RoundTripper
}

func (t *credentialTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Header.Get("Authorization") != "Bearer "+credentialMarker {
		return t.base.RoundTrip(req)
	}
	token, err := t.resolver.resolve(req.Context())
	if err != nil {
		return nil, err
	}
	cloned := req.Clone(req.Context())
	cloned.Header.Set("Authorization", "Bearer "+token)
	return t.base.RoundTrip(cloned)
}
