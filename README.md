# Terraform Provider for Cloudflare Tunnels (Training Project)

[![Go](https://github.com/soulnai/cloudflare_tf_provider/actions/workflows/go.yml/badge.svg)](https://github.com/soulnai/cloudflare_tf_provider/actions/workflows/go.yml)

> **DISCLAIMER**: This is a **training project** created for educational purposes to learn how to build Terraform providers using the Terraform Plugin Framework. It is **NOT** intended for production use. For production, please use the official [Cloudflare Terraform Provider](https://registry.terraform.io/providers/cloudflare/cloudflare/latest).

## Overview

This provider allows you to manage Cloudflare Tunnels via Terraform. It implements basic CRUD operations for Tunnels and includes a Data Source to read tunnel information.

## Features

- **Resources**:
  - `cloudflare-tunnel`: Create, Read, Update, and Delete Cloudflare Tunnels.
- **Data Sources**:
  - `cloudflare-tunnel`: Retrieve information about an existing tunnel.

## Prerequisites

- [Go](https://go.dev/) 1.18+
- [Terraform](https://www.terraform.io/) 1.0+
- A Cloudflare Account ID and API Token.

## Building and Installing

Since this is a local training provider, you need to build it and configure Terraform to use the local binary.

1.  **Build the Provider**:
    ```bash
    go build -o terraform-provider-cloudflare-tunnel.exe
    ```

2.  **Configure Development Overrides**:
    Create or edit your `.terraformrc` (Windows: `%APPDATA%\terraform.rc`) to point to your local build directory:

    ```hcl
    provider_installation {
      dev_overrides {
        "registry.terraform.io/cloudflare/cloudflare-tunnel" = "/path/to/your/project/cloudflare_tf_provider"
      }
      direct {}
    }
    ```

## Usage Example

```hcl
terraform {
  required_providers {
    cloudflare-tunnel = {
      source  = "registry.terraform.io/cloudflare/cloudflare-tunnel"
      version = "1.0.0"
    }
  }
}

provider "cloudflare-tunnel" {
  api_token  = "YOUR_API_TOKEN"
  account_id = "YOUR_ACCOUNT_ID"
  base_url   = "https://api.cloudflare.com/client/v4"
}

resource "cloudflare-tunnel" "example" {
  name         = "my-training-tunnel"
  tunnel_token = "AQIDBAUGBwgBAgMEBQYHCAECAwQFBgcIAQIDBAUGBwg=" # 32 bytes base64
}

data "cloudflare-tunnel" "lookup" {
  id = cloudflare-tunnel.example.id
}

output "tunnel_id" {
  value = data.cloudflare-tunnel.lookup.id
}
```
