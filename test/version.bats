#!/usr/bin/env bats

@test "prints version provided on compile" {
  run hcloud version
  test $status -eq 0
  test "$output" != "hcloud was not build properly"
}
