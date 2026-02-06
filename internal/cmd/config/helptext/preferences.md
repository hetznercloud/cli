| Option | Description | Type | Config key | Environment variable | Flag |
| --- | --- | --- | --- | --- | --- |
| debug | Enable debug output | boolean | debug | HCLOUD\_DEBUG | --debug |
| debug-file | File to write debug output to | string | debug\_file | HCLOUD\_DEBUG\_FILE | --debug-file |
| default-ssh-keys | Default SSH Keys for new Servers and Storage Boxes | string list | default\_ssh\_keys | HCLOUD\_DEFAULT\_SSH\_KEYS |  |
| endpoint | Hetzner Cloud API endpoint | string | endpoint | HCLOUD\_ENDPOINT | --endpoint |
| hetzner-endpoint | Hetzner API endpoint | string | hetzner\_endpoint | HETZNER\_ENDPOINT | --hetzner-endpoint |
| no-experimental-warnings | If true, experimental warnings are not shown | boolean | no\_experimental\_warnings | HCLOUD\_NO\_EXPERIMENTAL\_WARNINGS | --no-experimental-warnings |
| poll-interval | Interval at which to poll information, for example action progress | duration | poll\_interval | HCLOUD\_POLL\_INTERVAL | --poll-interval |
| quiet | If true, only print error messages | boolean | quiet | HCLOUD\_QUIET | --quiet |
| sort.\<resource\> | Default sorting for resource | string list | sort.\<resource\> | HCLOUD\_SORT\_\<RESOURCE\> |  |
