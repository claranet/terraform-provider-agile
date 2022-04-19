provider "agile" {
  username = var.username # optionally use AGILE_USERNAME env var
  password = var.password # optionally use AGILE_PASSWORD env var
  api_url  = var.api_url  # optionally use AGILE_API env var

  # you may need to allow insecure TLS communications unless you have configured
  # certificates for your controller
  allow_insecure = var.insecure # optionally use AGILE_INSECURE env var
}
