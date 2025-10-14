storage "file" {
  path = "/mnt/vault/storage"
}

listener "tcp" {
  address = "0.0.0.0:8200"
  tls_disable = 1
}

disable_mlock = true

api_addr = "http://0.0.0.0:8200"
ui = true