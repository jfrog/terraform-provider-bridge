# Terraform Provider for JFrog Bridge

## Quick Start

Create a new Terraform file with the `bridge` provider to manage JFrog Bridges:

### HCL Example

```terraform
# Required for Terraform 1.0 and later
terraform {
  required_providers {
    bridge = {
      source  = "jfrog/bridge"
      version = "1.0.0"
    }
  }
}

provider "bridge" {
  url          = "https://myinstance.jfrog.io"
  access_token = var.jfrog_access_token  # Or use JFROG_ACCESS_TOKEN env var
  insecure     = false                    # Set to true to skip TLS verification
}

# Create a JFrog Bridge
resource "bridge" "example" {
  bridge_id     = "my-bridge"
  pairing_token = var.pairing_token  # Required on create, generated on bridge server

  remote = {
    url      = "https://saas-jpd.jfrog.io"
    insecure = false
  }

  local = {
    url = "https://self-managed-jpd:8082"
  }
}
```

Initialize Terraform:
```sh
$ terraform init
```

Plan (or Apply):
```sh
$ terraform plan
```

Detailed documentation of the resource and attributes will be available on [Terraform Registry](https://registry.terraform.io/providers/jfrog/bridge/latest/docs).

## Resources

### bridge

Manages a JFrog Bridge connection between a Bridge Server (SaaS JPD) and Bridge Client (Self-managed JPD).

#### Required Arguments

- `bridge_id` - (Required) Unique identifier of the bridge. Changing this forces a new resource.
- `pairing_token` - (Required on create) Pairing token generated on the bridge server.
- `remote` - (Required) Remote (bridge server) configuration block:
  - `url` - (Required) URL of the bridge server (remote JPD).
  - `insecure` - (Optional) Allow insecure TLS when connecting to the remote.
  - `proxy` - (Optional) Proxy configuration block.
- `local` - (Required) Local (bridge client) configuration block:
  - `url` - (Required) URL of the bridge client (local JPD).
  - `anonymous_endpoints` - (Optional) List of anonymous endpoints allowed through the bridge.

#### Optional Arguments

- `min_tunnels` - (Optional) Minimum tunnels.
- `max_tunnels` - (Optional) Maximum tunnels.
- `target_usage` - (Optional) Target usage thresholds block:
  - `low` - (Optional) Low threshold.
  - `high` - (Optional) High threshold.
- `jobs` - (Optional) Job configuration block:
  - `tunnel_creation` - (Optional) Tunnel creation settings.
  - `tunnel_closing` - (Optional) Tunnel closing settings.

#### Read-Only Attributes

- `id` - Internal Terraform resource ID.
- `created_at` - Timestamp when the bridge was created.

#### Import

Existing bridges can be imported:

```sh
terraform import bridge.example my-bridge-id
```

## Requirements

- Terraform 1.0+
- JFrog Platform with Bridge API access
- Access Token with Admin privileges
- Pairing token from bridge server (for creating new bridges)

## Authentication

The provider supports the following authentication methods (order of precedence):

1. **Access Token**: Set via `access_token` attribute or `JFROG_ACCESS_TOKEN` environment variable
2. **OIDC**: Configure via `oidc_provider_name` attribute
3. **Terraform Cloud Workload Identity**: Automatic when running in Terraform Cloud

### Provider Configuration

```terraform
provider "bridge" {
  url          = "https://myinstance.jfrog.io"  # Or JFROG_URL env var
  access_token = "my-admin-token"               # Or JFROG_ACCESS_TOKEN env var
  insecure     = false                          # Skip TLS verification (use with caution)
}
```

## API Endpoints

This provider uses the following JFrog Bridge API endpoints:

- `POST /bridge-client/api/v1/bridges` - Create a new bridge
- `PATCH /bridge-client/api/v1/bridges/{id}` - Update bridge configuration
- `DELETE /bridge-client/api/v1/bridges/{id}` - Delete a bridge

## Versioning

In general, this project follows [semver](https://semver.org/) as closely as we can for tagging releases of the package. We've adopted the following versioning policy:

* We increment the **major version** with any incompatible change to functionality, including changes to the exported Go API surface or behavior of the API.
* We increment the **minor version** with any backwards-compatible changes to functionality.
* We increment the **patch version** with any backwards-compatible bug fixes.

## License

Copyright (c) 2025 JFrog.

Apache 2.0 licensed, see [LICENSE][LICENSE] file.

[LICENSE]: ./LICENSE
