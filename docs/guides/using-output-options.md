# Using output options

The CLI allows you to customize the output format of some commands using the `--output` flag. 
This guide shows you how to use this feature and use it in combination with other tools.

## JSON

`describe`, `list` and `create` commands support JSON output format. To use it, simply add the `--output json` flag to your command.

For example, to get the details of a location in JSON format, you can use the following command:

```bash
$ hcloud location describe fsn1 --output json
{
  "id": 1,
  "name": "fsn1",
  "description": "Falkenstein DC Park 1",
  "country": "DE",
  "city": "Falkenstein",
  "latitude": 50.47612,
  "longitude": 12.370071,
  "network_zone": "eu-central"
}
```

You can combine this with other tools to process the output. For example, you can use `jq` to filter the output:

```bash
$ hcloud location describe fsn1 --output json | jq '.name'
"fsn1"
```

```bash
$ hcloud location describe fsn1 --output json | jq '{id, name}'
{
  "id": 1,
  "name": "fsn1"
}
```

`list` commands return a list of objects in JSON format. For example, to get a list of all locations in JSON format, you can use the following command:

```bash
$ hcloud location list --output json
[
  {
    "id": 1,
    "name": "fsn1",
    "description": "Falkenstein DC Park 1",
    "country": "DE",
    "city": "Falkenstein",
    "latitude": 50.47612,
    "longitude": 12.370071,
    "network_zone": "eu-central"
  },
  {
    "id": 2,
    "name": "nbg1",
    "description": "Nuremberg DC Park 1",
    "country": "DE",
    "city": "Nuremberg",
    "latitude": 49.452102,
    "longitude": 11.076665,
    "network_zone": "eu-central"
  },
  ...
]
```

Once again, you can use `jq` to filter the output. Following example shows how to get the names of all locations of which the network zone is `eu-central`:

```bash
$ hcloud location list --output json | jq '[.[] | select(.network_zone == "eu-central") | .name]'    
[
  "fsn1",
  "nbg1",
  "hel1"
]
```

## YAML

`describe`, `list` and `create` commands support YAML output format as well.

```bash
$ hcloud location describe fsn1 --output yaml
id: 1
name: fsn1
description: Falkenstein DC Park 1
country: DE
city: Falkenstein
latitude: 50.47612
longitude: 12.370071
network_zone: eu-central
```

For YAML, you can use `yq` instead of `jq`.


```bash
$ hcloud location list --output yaml | yq '.[] | [{"id": .id, "name": .name}]'
- id: 1
  name: fsn1
- id: 2
  name: nbg1
- id: 3
  name: hel1
- id: 4
  name: ash
- id: 5
  name: hil
- id: 6
  name: sin
```

## Go Template Format

`describe` commands support the Go string template format as well. You can read up on the syntax in the 
[Go documentation](https://pkg.go.dev/text/template/). The data structures passed to the template are defined
by our API and can be found in [hcloud-go](https://pkg.go.dev/github.com/hetznercloud/hcloud-go/v2/hcloud/schema).

For example, you could obtain the number of cores of a server using the following command:

```bash
$ hcloud server describe my-server --output format='{{.ServerType.Cores}}'
2
```

## Table options

`list` commands support table options as well. These options allow you to customize the output format of the output table,
if not using JSON or YAML.

> [!NOTE]
> You can also combine both of the below options to use them at once: ``--output noheader --output columns=id,name,network_zone``

### noheader

This option removes the header from the output table.

```bash
$ hcloud location list --output noheader 
1   fsn1   Falkenstein DC Park 1   eu-central     DE   Falkenstein  
2   nbg1   Nuremberg DC Park 1     eu-central     DE   Nuremberg    
3   hel1   Helsinki DC Park 1      eu-central     FI   Helsinki     
4   ash    Ashburn, VA             us-east        US   Ashburn, VA  
5   hil    Hillsboro, OR           us-west        US   Hillsboro, OR
6   sin    Singapore               ap-southeast   SG   Singapore   
```

### columns

This option allows you to filter by columns.

```bash
$ hcloud location list --output columns=id,name,network_zone 
ID   NAME   NETWORK ZONE
1    fsn1   eu-central  
2    nbg1   eu-central  
3    hel1   eu-central  
4    ash    us-east     
5    hil    us-west     
6    sin    ap-southeast
```

Using the ``--help`` flag will show you a list of all available columns for this command. Note that these might include
more than the default columns.
