| Option | Description | Type | Config key | Environment variable | Flag |
| --- | --- | --- | --- | --- | --- |
| debug | Enable debug output | boolean | debug | HCLOUD_DEBUG | --debug |
| debug-file | File to write debug output to | string | debug_file | HCLOUD_DEBUG_FILE | --debug-file |
| default-ssh-keys | Default SSH Keys for new Servers and Storage Boxes | string list | default_ssh_keys | HCLOUD_DEFAULT_SSH_KEYS |  |
| endpoint | Hetzner Cloud API endpoint | string | endpoint | HCLOUD_ENDPOINT | --endpoint |
| hetzner-endpoint | Hetzner API endpoint | string | hetzner_endpoint | HETZNER_ENDPOINT | --hetzner-endpoint |
| no-experimental-warnings | If true, experimental warnings are not shown | boolean | no_experimental_warnings | HCLOUD_NO_EXPERIMENTAL_WARNINGS | --no-experimental-warnings |
| poll-interval | Interval at which to poll information, for example action progress | duration | poll_interval | HCLOUD_POLL_INTERVAL | --poll-interval |
| quiet | If true, only print error messages | boolean | quiet | HCLOUD_QUIET | --quiet |
| sort.<resource> | Default sorting for resource | string list | sort.<resource> | HCLOUD_SORT_<RESOURCE> |  |
