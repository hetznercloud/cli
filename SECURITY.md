# Security policy

## Reporting a vulnerability

Do not open a public issue for a suspected vulnerability. Use GitHub's private vulnerability reporting for
`hetznercloud/cli` when available, or contact Hetzner through the
[Support Center](https://console.hetzner.com/support). Include affected versions, impact, reproduction steps, and any
suggested remediation. Do not access data or projects that you do not own or have explicit permission to test.

## Supported versions

Security fixes are delivered in current releases. Users should upgrade to the latest published version before
reporting an issue that may already have been fixed.

## Security properties

- API tokens are read from the selected configuration/environment and are never intentionally emitted by structured
  diagnostics.
- Configuration and file-backed diagnostics are written with owner-only permissions (`0600`). Configuration updates
  use a same-directory temporary file, synchronization, and atomic rename.
- Debugging is local and opt-in. It can contain HTTP request/response details from `hcloud-go`; protect and remove
  debug files as sensitive operational data.
- CI checks reachable Go vulnerabilities with `govulncheck`. Release artifacts include checksums, SBOMs, and GitHub
  artifact attestations; automation dependencies and the container base image are pinned.

Verify downloaded artifacts using the release checksums and provenance supplied with the release. Building from a
source checkout should use the versions pinned in `mise.toml`.
