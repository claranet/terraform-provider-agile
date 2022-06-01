# Terraform Provider For Huawei Agile Controller
Terraform Provider for Huawei Agile Controller DCN

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- [Go](https://golang.org/doc/install) >= 1.17

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the `make build` command:

```sh
$ make build
```

## Using The Provider

<!-- https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin -->

If you are building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/cli/plugins/index.html) After placing it into your plugins directory, run `terraform init` to initialize it.

ex.
```hcl

terraform {
  required_providers {
    agile = {
      source = "claranet/agile"
    }
  }
}

# Configure provider with your huawei agile controller credentials.
provider "agile" {
  username       = "admin"
  password       = "password"
  api_url        = "https://<IP>:<PORT>"
  allow_insecure = true
}
```

## Documentation

Full, comprehensive documentation is available on the Terraform website:

https://registry.terraform.io/providers/claranet/agile/latest/docs

##  Developing The Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate-docs`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

### Using Provider Local

```shell
export TERRAFORM_AGILE_VERSION=0.0.1
make build-local
terraform init -plugin-dir=${HOME}/.terraform.d/plugin-cache
```