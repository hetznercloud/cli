#!/usr/bin/env python3
"""Read-only Git credential helper using a Passbolt resource ID as the URL path."""

import json
import subprocess
import sys
import uuid


def main():
    if sys.argv[1:] != ["get"]:
        return 0  # Retrieval only; never store or erase a secret.

    fields = {}
    for line in sys.stdin:
        line = line.rstrip("\r\n")
        if not line:
            break
        key, separator, value = line.partition("=")
        if not separator:
            raise ValueError("invalid credential request")
        fields[key] = value

    # Validate the lookup key before calling Passbolt. The URL is an identifier,
    # not the address of the Passbolt server (which Passbolt configures itself).
    resource_id = str(uuid.UUID(fields["path"]))
    result = subprocess.run(
        ["passbolt", "get", "resource", "--id", resource_id, "--json"],
        stdin=subprocess.DEVNULL,
        capture_output=True,
        text=True,
        check=True,
    )
    token = json.loads(result.stdout)["password"]
    if not isinstance(token, str) or not token.strip() or any(c in token for c in "\r\n\x00"):
        raise ValueError("invalid password")
    sys.stdout.write(f"password={token}\n\n")
    return 0


if __name__ == "__main__":
    try:
        sys.exit(main())
    except (OSError, ValueError, KeyError, TypeError, subprocess.SubprocessError):
        # Neither raw output nor exception details are safe to print here.
        print("Passbolt credential lookup failed", file=sys.stderr)
        sys.exit(1)
