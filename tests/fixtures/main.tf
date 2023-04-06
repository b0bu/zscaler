data "external" "zscaler_ranges" {
  program = ["zscaler"]
}

locals {
  // external provider only supports map[string]string as input, prefix_list must be split into []string for filtering
  zscaler_ranges = toset(split(" ", data.external.zscaler_ranges.result.prefix_list))
  // filter out ipv6 addresses
  zscaler_ipv4_ranges = [
    for cidr in local.zscaler_ranges : cidr
    if can(cidrnetmask(cidr))
  ]
}

output "zscaler_ipv4_ranges" {
  value = local.zscaler_ipv4_ranges
}

