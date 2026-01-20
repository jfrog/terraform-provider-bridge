# Terraform Provider Bridge - Project Summary

## Overview

Terraform provider for JFrog Bridge API v1, implementing full bridge lifecycle management (create/update/delete).

## Features Implemented

### Resources

1. **bridge**
   - Manages a JFrog Bridge lifecycle (define/update/delete)
   - Uses:
     - `POST /bridge-client/api/v1/bridges` - Create new bridge
     - `PATCH /bridge-client/api/v1/bridges/{id}` - Update bridge configuration
     - `DELETE /bridge-client/api/v1/bridges/{id}` - Delete bridge
   - Supports import functionality

## Project Structure

```
terraform-provider-bridge/
├── main.go                           # Provider entry point
├── go.mod                            # Go module definition
├── go.sum                            # Go module checksums
├── GNUmakefile                       # Build and test automation
├── LICENSE                           # Apache 2.0 license
├── README.md                         # User documentation
├── CHANGELOG.md                      # Version history
├── CODEOWNERS                        # Code ownership
├── CONTRIBUTIONS.md                  # Contribution guidelines
├── sample.tf                         # Sample Terraform configuration
├── terraform-registry-manifest.json  # Terraform registry metadata
├── .goreleaser.yml                   # Release configuration
├── .gitignore                        # Git ignore patterns
├── pkg/bridge/
│   ├── provider.go                   # Provider implementation
│   └── resource_bridge.go            # Resource: bridge lifecycle
├── docs/
│   ├── index.md                      # Provider documentation
│   └── resources/
│       └── bridge.md                 # Bridge resource documentation
├── examples/
│   ├── provider/
│   │   └── provider.tf               # Provider configuration examples
│   └── resources/
│       └── bridge/
│           └── resource.tf           # Bridge resource example
├── templates/
│   └── index.md.tmpl                 # Documentation template
└── tools/
    └── tools.go                      # Build tools
```

## Provider Configuration

The provider supports multiple authentication methods (preferred first):

1. **Access Token** - Via `access_token` attribute or `JFROG_ACCESS_TOKEN` environment variable
2. **OIDC** - Via `oidc_provider_name` attribute
3. **Terraform Cloud Workload Identity** - Automatic when running in TFC

Example configuration:

```terraform
provider "bridge" {
  url          = "https://myinstance.jfrog.io"
  access_token = "my-admin-token"  # Or use JFROG_ACCESS_TOKEN env var
  insecure     = false             # Skip TLS verification (use with caution)
}
```

## API Endpoints Implemented

### Define a New Bridge
- **Method:** POST
- **Endpoint:** `/bridge-client/api/v1/bridges`
- **Description:** Creates a new bridge between Bridge Server and Bridge Client
- **Request Body:** `{ "bridge_id": "", "remote": "<url>", "local": "<url>", "pairing_token": "" }`
- **Response Codes:** 201 (Success), 400/401/403 (Error)

### Modify Bridge Configuration
- **Method:** PATCH
- **Endpoint:** `/bridge-client/api/v1/bridges/{id}`
- **Description:** Updates bridge configuration and thresholds
- **Response Codes:** 204 (Success), 400/401/403 (Error)

### Delete a Bridge
- **Method:** DELETE
- **Endpoint:** `/bridge-client/api/v1/bridges/{id}`
- **Description:** Removes a bridge and terminates its tunnels
- **Response Codes:** 204 (Success), 400/401/403 (Error)

## Building the Provider

```bash
# Initialize dependencies
go mod tidy

# Build the provider
make build

# Install locally for testing
make install

# Run tests
make test

# Run acceptance tests
make acceptance
```

## Usage Examples

### Minimal Bridge Resource

```terraform
resource "bridge" "demo" {
  bridge_id     = "demo"
  pairing_token = "<pairing token from bridge server>"  # Required on create

  remote = {
    url = "https://your_SaaS_JPD.org"
  }

  local = {
    url = "https://your_self_managed_JPD:8082"
  }
}
```

### Full Bridge Resource with All Options

```terraform
resource "bridge" "demo" {
  bridge_id     = "demo"
  pairing_token = "<pairing token from bridge server>"

  remote = {
    url      = "https://your_SaaS_JPD.org"
    insecure = false
    proxy = {
      enabled               = true
      cache_expiration_secs = 3600
      key                   = "platform"
      scheme_override       = ""
    }
  }

  local = {
    url = "https://your_self_managed_JPD:8082"
    anonymous_endpoints = [
      ".*/system/(ping|readiness|liveness)",
    ]
  }

  min_tunnels = 2
  max_tunnels = 10

  target_usage = {
    low  = 1
    high = 5
  }

  jobs = {
    tunnel_creation = {
      interval_minutes = 1
    }
    tunnel_closing = {
      cron_expr                = "0 0 * * *"
      allow_close_used_tunnels = true
    }
  }
}
```

### Import Existing Bridge

```bash
terraform import bridge.demo my-bridge-id
```

## Resource Schema

### Required Arguments

| Attribute | Type | Description |
|-----------|------|-------------|
| `bridge_id` | String | Unique identifier of the bridge. Changing forces replacement. |
| `pairing_token` | String (Sensitive) | Pairing token from bridge server. Required on create. |
| `remote` | Object | Remote (bridge server) configuration. |
| `remote.url` | String | URL of the bridge server (remote JPD). |
| `local` | Object | Local (bridge client) configuration. |
| `local.url` | String | URL of the bridge client (local JPD). |

### Optional Arguments

| Attribute | Type | Description |
|-----------|------|-------------|
| `remote.insecure` | Bool | Allow insecure TLS to remote. |
| `remote.proxy` | Object | Proxy configuration. |
| `local.anonymous_endpoints` | List(String) | Anonymous endpoints allowed. |
| `min_tunnels` | Int64 | Minimum tunnels. |
| `max_tunnels` | Int64 | Maximum tunnels. |
| `target_usage` | Object | Target usage thresholds. |
| `jobs` | Object | Job configuration for tunnel creation/closing. |

### Read-Only Attributes

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | String | Internal Terraform resource ID. |
| `created_at` | String | Timestamp when bridge was created. |

## Security Considerations

1. **Admin Privileges Required:** All operations require an Access Token with Admin privileges
2. **Sensitive Data:** Pairing tokens and access tokens are marked as sensitive in Terraform
3. **TLS Verification:** Use `insecure = true` only for testing with self-signed certificates
4. **Pairing Token:** One-time use; stored in state after creation for consistency

## Development Notes

- Built with Terraform Plugin Framework (v1.16.1)
- Uses JFrog shared library (v1.30.6) for common functionality
- Follows patterns established by other JFrog Terraform providers
- Compatible with Go 1.24.0+
- Supports Terraform 1.0+

## Version

Initial version: 1.0.0

## Dependencies

Key dependencies:
- github.com/hashicorp/terraform-plugin-framework v1.16.1
- github.com/jfrog/terraform-provider-shared v1.30.6
- github.com/go-resty/resty/v2 v2.16.5
- github.com/hashicorp/terraform-plugin-docs v0.20.1

## Next Steps

1. Add unit tests for resource
2. Add acceptance tests
3. Set up CI/CD pipeline
4. Prepare for Terraform Registry publication
5. Add data sources for bridge debug snapshots (optional)

## License

Apache 2.0 - Copyright (c) 2025 JFrog Ltd
