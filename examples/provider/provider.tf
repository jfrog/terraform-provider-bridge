terraform {
  required_providers {
    bridge = {
      source  = "jfrog/bridge"
      version = "~> 1.0"
    }
  }
}

# Configure with explicit credentials
provider "bridge" {
  url          = "https://myinstance.jfrog.io"
  access_token = "my-admin-token"  # Or use JFROG_ACCESS_TOKEN env var
  insecure     = false             # Set to true to skip TLS verification
}

# Alternative: Configure with environment variables
# export JFROG_URL="https://myinstance.jfrog.io"
# export JFROG_ACCESS_TOKEN="my-admin-token"
#
# provider "bridge" {
#   # url and access_token will be read from environment variables
# }

# Alternative: Configure with OIDC
# provider "bridge" {
#   alias              = "oidc"
#   url                = "https://myinstance.jfrog.io"
#   oidc_provider_name = "my-oidc-provider"
# }
