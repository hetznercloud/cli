#!/usr/bin/env bats

@test "prints version" {
  run hcloud version
  test $status -eq 0
  test "$output" = "hcloud 1.0.0"
}
