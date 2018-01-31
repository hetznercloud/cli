# Changes

## v1.3.0

* Support passing user data on server creation ([issue #70](https://github.com/hetznercloud/hcloud-go/issues/70))
* Fix leaking response body by not closing it ([issue #68](https://github.com/hetznercloud/hcloud-go/issues/68))

## v1.2.0

* Add `WatchProgress()` to action client
* Use correct error code for ratelimit error (deprecated
  `ErrorCodeLimitReached`, added `ErrorCodeRateLimitExceeded`)

## v1.1.0

* Add `Image` field to `Server`
