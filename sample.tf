terraform {
  required_providers {
    bridge = {
      source  = "jfrog/bridge"
      version = "0.0.1"
    }
  }
}

provider "bridge" {
  url          = "https://client_URL:8082"
  access_token = ""
  insecure     = true  # Skip TLS certificate verification (use with caution)
}

# Minimal bridge resource example
resource "bridge" "resource_name" {
  bridge_id     = "resource_name"
  pairing_token = "<paring_token>" # required on create

  remote = {
    url = "https://remote_URL"
    insecure = true
    # proxy = {
    #   enabled               = true
    #   cache_expiration_secs = 3600
    #   key                   = "platform"
    #   scheme_override       = ""
    # }
  }

  local = {
    url = "https://:8082"
    anonymous_endpoints = [
      ".*/system/(ping|readiness|liveness)",
   ]
  }
}