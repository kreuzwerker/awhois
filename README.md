# awhois

`awhois` determines if a given IP address (v4 or v6) is part of Amazon Web Services (AWS).

## Usage

`awhois :ip-address`

This will either

* return an exit code of `0` and a JSON that indicates to which service(s) this IP belongs
* return an exit code of `1` and a JSON with an empty list of matches

## Examples

### Using IPv6

```
awhois 2a05:d050:8080:0:0:0:0:0

{
  "IP": "2a05:d050:8080::",
  "Matches": [
    "2a05:d050:8080::/44 (AMAZON eu-west-1)"
  ]
}
```

### Using IPv4

```
awhois 54.76.43.209

{
  "IP": "54.76.43.209",
  "Matches": [
    "54.76.0.0/15 (AMAZON eu-west-1)",
    "54.76.0.0/15 (EC2 eu-west-1)"
  ]
}
```

### Using loopback

```
awhois 127.0.0.1

{
  "IP": "127.0.0.1",
  "Matches": []
}
```
