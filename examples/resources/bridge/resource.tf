resource "bridge_bridge" "example" {
  # Required bridge identifier and pairing token (pairing token is needed only on create)
  bridge_id     = "demo"
  pairing_token = "<pairing token from bridge server>"

  remote {
    url      = "https://your_SaaS_JPD.org"
    insecure = false

    proxy {
      enabled               = true
      cache_expiration_secs = 3600
      key                   = "platform"
      scheme_override       = ""
    }
  }

  local {
    url = "https://your_self_managed_JPD:8082"
    anonymous_endpoints = [
      ".*/system/(ping|readiness|liveness)",
    ]
  }

  min_tunnels = 2
  max_tunnels = 10

  target_usage {
    low  = 1
    high = 5
  }

  jobs {
    tunnel_creation {
      interval_minutes = 1
    }
    tunnel_closing {
      cron_expr                = "0 0 * * *"
      allow_close_used_tunnels = true
    }
  }
}
