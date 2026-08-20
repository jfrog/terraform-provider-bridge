## 1.0.1 (August 20, 2026)

SECURITY:

* **CVE-2026-39821:** Upgrade Go to 1.25.13 and `golang.org/x/net` to v0.58.0 to reject ASCII-only Punycode-encoded labels.
* **CVE-2026-56865:** Verify transparency log tiles against their parents to prevent checksum database bypass.
* **CVE-2026-56864:** Ignore unrelated unauthenticated hashes returned by checksum database lookups.
* **CVE-2026-33818:** Enforce a maximum recursion depth when decoding ASN.1 data.
* **CVE-2026-46600:** Prevent panics while parsing malformed SVCB and HTTPS DNS records.
* **CVE-2026-56862:** Limit accepted post-handshake TLS messages to prevent resource exhaustion.
* **CVE-2026-56859:** Guard XML decoding recursion depth to prevent stack exhaustion.
* **CVE-2026-56860:** Eliminate quadratic complexity when resolving relative URL paths.
* **CVE-2026-56858:** Correct JavaScript regular expression context tracking in HTML templates.
* **CVE-2026-56853:** Apply `ReadHeaderTimeout` while detecting unencrypted HTTP/2 connections.
* **CVE-2026-25680:** Bound CPU consumption when parsing malicious HTML input.
* **CVE-2026-42506:** Correct namespaced element handling in foreign HTML content to prevent XSS.
* **CVE-2026-42502:** Correct foreign-content element parsing and rendering to prevent XSS.
* **CVE-2026-25681:** Correct character-reference handling in HTML DOCTYPE nodes to prevent XSS.
* **CVE-2026-27136:** Preserve duplicate HTML attributes consistently to prevent mutation XSS.
* **CVE-2026-46595:** Enforce source-address restrictions for SSH `VerifiedPublicKeyCallback` authentication.
* **CVE-2026-42508:** Enforce `@revoked` status for SSH certificate authority signature keys.
* **CVE-2026-39834:** Prevent infinite loops when writing SSH channel payloads larger than 4 GB.
* **CVE-2026-39833:** Reject unsupported SSH agent key constraints instead of silently dropping them.
* **CVE-2026-39832:** Preserve SSH agent constraints when forwarding keys through memory keyrings.
* **CVE-2026-39831:** Require user presence for FIDO/U2F SSH security-key signatures.
* **CVE-2026-39830:** Handle unexpected SSH responses without deadlocking clients.
* **CVE-2026-39829:** Bound SSH RSA and DSA parameter sizes to prevent denial of service.
* **CVE-2026-46597:** Prevent AES-GCM SSH packet length underflow and panic.
* **CVE-2026-39828:** Enforce SSH certificate restrictions when callbacks return nil permissions.
* **CVE-2026-39827:** Release rejected SSH channels to prevent memory-exhaustion denial of service.
* **CVE-2026-39835:** Prevent SSH server panics during host-key checking and authentication.
* **CVE-2026-46598:** Reject pathological SSH agent inputs that could panic clients.
* **CVE-2025-47914:** Validate SSH agent identity-request message sizes to prevent panics.
* **CVE-2025-58181:** Bound SSH GSSAPI mechanism counts to prevent unbounded memory consumption.
* **CVE-2026-1229:** Upgrade `github.com/cloudflare/circl` to v1.6.5 to correct secp384r1 `CombinedMult` calculations.

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
