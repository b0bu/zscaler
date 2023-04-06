# zscaler_ip_prefixes

```bash
go test -v 
```

```bash
go get github.com/b0bu/zscaler
```

implementing this for tf where `local.zscaler_ipv4_ranges` can be used to update zscaler prefixes.

```hcl
data "external" "zscaler_ranges" {
  program = ["zscaler"]
}

locals {
  zscaler_ranges = toset(split(" ", data.external.zscaler_ranges.result.prefix_list))
  zscaler_ipv4_ranges = [
    for cidr in local.zscaler_ranges : cidr
    if can(cidrnetmask(cidr))
  ]
}
```