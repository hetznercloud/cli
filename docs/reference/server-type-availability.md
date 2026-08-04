# Server type support and availability

The Server Type API exposes two different concepts for each location:

- **Supported location**: the server type can be used in that location, subject to deprecation and the complete create
  request.
- **Availability signal**: the API's current capacity observation. It is advisory, can change at any time, and is not
  a reservation or a guarantee that a create request will succeed.

`hcloud server-type list` therefore uses these columns:

| Column | Meaning |
| --- | --- |
| `location` | Supported locations that have not passed their deprecation deadline. |
| `location_availability_signal` | Supported locations whose current advisory signal is available. |
| `location_recommended` | Locations recommended by the API. |

The previous `location_available` name overstated what the field could prove and has been removed as a deliberate
breaking change.

Availability is evaluated for the complete server-create request. Server type and location are only part of that
request: image architecture and product support, primary IP assignment, public-network settings, private Networks,
Volumes, Firewalls, Placement Groups, SSH keys, user data and other options can independently make a request invalid
or unsatisfiable. The API's create response is authoritative.

For planning, enumerate supported locations from `location`, filter impossible combinations using the relevant
resource metadata, and treat the availability signal only as a prioritization hint. Do not interpret one false signal
as proof that every configuration is impossible, and do not interpret one true signal as capacity reserved for a
later request. Automation must handle a create failure, refresh metadata, and select another valid combination under
its own policy. The CLI does not create probe resources because probes are billable, race with real creation, and can
themselves alter capacity.

Hetzner documents the field in the [API changelog](https://docs.hetzner.cloud/changelog) and the
[Cloud API reference](https://docs.hetzner.cloud/reference/cloud).
