# Development

## Prerequisites and setup

1. Install [Git](https://git-scm.com/), [mise](https://mise.jdx.dev/), and a platform C toolchain (required by the Go
   race detector).
2. Clone the repository and enter it.
3. Run `mise trust` after reviewing `mise.toml`, then `mise install`.
4. Run `go mod download`.
5. Verify the checkout with `go test ./...` and `golangci-lint run`.

`mise.toml` pins the supported Go release and repository tools. Do not install an unversioned `latest` tool in CI or
release automation. The module uses Go's toolchain mechanism, so commands may download the declared toolchain when it
is not already present.

No cloud credentials are required for builds or unit tests. E2E tests are separate and must only be run against an
explicitly selected disposable project; see `test/e2e` before enabling them.

## Common commands

```sh
# Build and run locally
go build -o ./dist/hcloud ./cmd/hcloud
./dist/hcloud version

# Unit and race tests
go test ./...
go test -race ./...

# Static and security checks
gofmt -l .
go vet ./...
golangci-lint run
govulncheck ./...

# Regenerate mocks, help tables, and the command manual
go generate ./...
git diff --exit-code

# Validate release configuration and make a non-publishing snapshot
goreleaser check
goreleaser release --snapshot --clean --skip=publish
```

Run `go mod tidy` after dependency changes and verify that it leaves no unexpected diff. Run `go fix ./...` when
adopting a new Go release, review every rewrite, then run the complete suite.

## Adding or changing commands

Prefer the generic command definitions in `internal/cmd/base` for standard resource operations. Resource packages
should contain only domain-specific flags, validation, API mapping, output fields and schemas. Pass
`cmd.Context()` to every call that can block.

Registration APIs return errors. Use the checked helpers in `internal/cmd/cmpl` and `internal/cmd/util` for completion,
required flags and filename flags; root construction must fail when the command graph is invalid. Use context-aware
error-returning output fields when rendering needs an API lookup.

Pflag getters return errors even when a flag was registered locally with the same type. Existing command code treats
those same-scope declarations as construction invariants. New code should prefer binding flags to a typed per-command
options structure, or return getter errors from `RunE`, when values cross a package or abstraction boundary.

## Tests

Command tests use `testutil.NewFixture`, generated `hcapi2` mocks, and Cobra I/O streams. Tests must not use the real
network, user configuration, stdin/stdout/stderr, or mutate process globals without a cleanup guard. Use
`testutil.SetTimezone` for timezone-sensitive output.

When an API client interface changes, regenerate its mock using the pinned tool:

```sh
go run go.uber.org/mock/mockgen@v0.6.0 \
  -package mock \
  -destination internal/hcapi2/mock/zz_hcapi_mock.go \
  github.com/hetznercloud/cli/internal/hcapi2 \
  ActionClient,CertificateClient,DatacenterClient,ImageClient,ISOClient,FirewallClient,FloatingIPClient,PrimaryIPClient,LocationClient,LoadBalancerClient,LoadBalancerTypeClient,NetworkClient,ServerClient,ServerTypeClient,SSHKeyClient,VolumeClient,PlacementGroupClient,RDNSClient,PricingClient,StorageBoxClient,StorageBoxTypeClient,ZoneClient
```

## Pull requests and releases

Keep generated files in the same commit as their source. Do not commit credentials, debug logs, build output, or local
analysis directories. PRs should state behavior changes, breaking CLI changes, tests run, and security or release
impact. A release is accepted only when checksums, SBOMs, signatures/notarization where configured, and artifact
attestations succeed.
