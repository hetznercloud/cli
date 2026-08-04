# Architecture

`hcloud` is a Go command-line application built around Cobra and the
[`hcloud-go`](https://github.com/hetznercloud/hcloud-go) API client. The executable is intentionally thin: it
loads configuration, constructs process-scoped dependencies, builds the command tree, and executes one command with a
cancelable context.

## Runtime flow

```text
cmd/hcloud
  -> state/config: parse flags, environment and TOML
  -> state: construct API client, terminal, logger and owned resources
  -> internal/cli: compose the Cobra command tree
  -> internal/cmd/<resource>: validate and execute one command
  -> internal/hcapi2: typed API extensions and completion helpers
  -> hcloud-go: Hetzner Cloud API
```

`SIGINT` and `SIGTERM` cancel the context passed to `cobra.Command.ExecuteContext`. Command handlers and completion
callbacks pass `cmd.Context()` to API calls. Context is not stored in `state.State`; this follows the Go convention
that context belongs in call parameters. `State.Close` releases owned diagnostic files, and `main` joins cleanup errors
with execution errors before selecting an exit code.

## Components

| Component | Responsibility |
| --- | --- |
| `cmd/hcloud` | Process lifecycle, signals, error rendering, exit status and version linker input. |
| `internal/cli` | Root command composition and global options. It contains no resource behavior. |
| `internal/cmd/base` | Reusable create/list/describe/update/delete/label/protection/RDNS command behavior. |
| `internal/cmd/<resource>` | Resource-specific flags, validation, API calls, schemas and text presentation. |
| `internal/cmd/cmpl` | Cobra completion composition, context propagation and construction-error registration. |
| `internal/cmd/output` | Table, JSON, YAML and template output. Context-aware computed fields return errors. |
| `internal/cmd/util` | Small, domain-independent command helpers and docopt-compatible argument validation. |
| `internal/hcapi2` | Narrow extensions around `hcloud-go`, including names and label completion. |
| `internal/state` | Process-scoped API client, terminal, action polling, optional local diagnostics and cleanup. |
| `internal/state/config` | Immutable option catalog, precedence, contexts and atomic TOML persistence. |
| `internal/testutil` | Mock-backed command fixtures and process-global test-state guards. |
| `test/e2e` | Opt-in integration tests against an explicitly configured Hetzner Cloud project. |
| `scripts` | Documentation generation and release-support entry points. |

## Dependency and state boundaries

- Commands depend on the `state.State` interface, not on process globals.
- API access is through `hcapi2.Client`; generated GoMock implementations keep resource tests isolated from the
  network.
- Options are copied into each `config.Config`. Tests and embedders can add options with `config.WithOptions` without
  mutating a package-wide registry.
- Configuration input is explicit: `LoadDefault`, `LoadFile`, or `LoadReader`. Arguments, warning output and the
  default path are constructor options rather than reads of `os.Args` or direct process output.
- The only build-time global is `main.defaultConfigPath`, populated by GoReleaser. It is kept at the executable boundary.
- The table output layer supports pure fields and context-aware error-returning fields. Hidden API lookups therefore
  honor cancellation and cannot silently degrade into plausible output.

## Error model and observability

Expected failures return errors to Cobra. `main` formats the final error once and preserves child-process exit codes
through `util.ExitError`. Multi-resource action waits return typed failures so a failed action is associated with its
owning resource; interrupted or incomplete waits do not report resources as deleted.

Diagnostics are local and opt-in through `--debug` or the equivalent configuration. The state creates diagnostic
files with mode `0600`, never logs raw command arguments, emits structured JSON lifecycle records through `log/slog`,
and closes the file. HTTP debug data from `hcloud-go` shares the selected writer. There is no outbound telemetry or
user tracking.

## Code style

- Go formatting is enforced with `gofmt`, `goimports`, and `gci` through golangci-lint.
- Package names are short and lower-case; exported identifiers carry Go documentation where useful.
- Cobra handlers use `RunE`/`PreRunE` and return errors. Constructor-time registration errors are collected by
  `internal/cmd/registration` and make root construction fail.
- Resource command definitions are declarative package variables. This is deliberate: they are immutable command
  specifications, while per-execution state lives in Cobra commands, local values, or `State`.
- Generated files start with `zz_` and are checked for a clean diff in CI.
- Command `Use` strings follow the documented docopt-compatible conventions because runtime argument validation and
  generated reference documentation share that contract.

## Delivery architecture

CI runs formatting, linting, generation checks, unit tests, race tests, vulnerability analysis, and platform builds.
GitHub Actions are pinned to full commit SHAs, workflows have explicit permissions, concurrency cancellation and
timeouts, and releases generate checksums, SBOMs and GitHub artifact attestations. GoReleaser builds binaries,
packages and a digest-pinned Alpine container image. Release signing/notarization failures are fatal.

See [Development](development.md), [Security](../SECURITY.md), and the
[2026-08-04 audit](audits/2026-08-04.md) for operational detail and decision history.
