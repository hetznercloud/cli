# Changes

## v1.9.0

* Add `AllWithOpts()` to server, Floating IP, image, and SSH key client
* Expose labels of servers, Floating IPs, images, and SSH Keys

## v1.8.0

* Add `WithPollInterval()` option to `Client` which allows to specify the polling interval
  ([issue #92](https://github.com/hetznercloud/hcloud-go/issues/92))
* Add `CPUType` field to `ServerType` ([issue #91](https://github.com/hetznercloud/hcloud-go/pull/91))

## v1.7.0

* Add `Deprecated ` field to `Image` ([issue #88](https://github.com/hetznercloud/hcloud-go/issues/88))
* Add `StartAfterCreate` flag to `ServerCreateOpts` ([issue #87](https://github.com/hetznercloud/hcloud-go/issues/87))
* Fix enum types ([issue #89](https://github.com/hetznercloud/hcloud-go/issues/89))

## v1.6.0

* Add `ChangeProtection()` to server, Floating IP, and image client
* Expose protection of servers, Floating IPs, and images

## v1.5.0

* Add `GetByFingerprint()` to SSH key client

## v1.4.0

* Retry all calls that triggered the API ratelimit
* Slow down `WatchProgress()` in action client from 100ms polling interval to 500ms

## v1.3.1

* Make clients using the old error code for ratelimiting work as expected
  ([issue #73](https://github.com/hetznercloud/hcloud-go/issues/73))

## v1.3.0

* Support passing user data on server creation ([issue #70](https://github.com/hetznercloud/hcloud-go/issues/70))
* Fix leaking response body by not closing it ([issue #68](https://github.com/hetznercloud/hcloud-go/issues/68))

## v1.2.0

* Add `WatchProgress()` to action client
* Use correct error code for ratelimit error (deprecated
  `ErrorCodeLimitReached`, added `ErrorCodeRateLimitExceeded`)

## v1.1.0

* Add `Image` field to `Server`
