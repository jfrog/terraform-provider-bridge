## 1.0.0 (January 22, 2026)

FEATURES:

* **New Resource:** `bridge` - Manages JFrog Bridge connections between Bridge Server (SaaS JPD) and Bridge Client (Self-managed JPD)
  * Create new bridges with pairing token
  * Update bridge configuration (remote, local, tunnels, jobs)
  * Delete bridges
  * Import existing bridges

NOTES:

* Initial release of the JFrog Bridge provider
* Supports Bridge Client API v1 endpoints for bridge lifecycle management
* Provider supports `insecure` option for TLS certificate verification
* Authentication via Access Token, OIDC, or Terraform Cloud Workload Identity
