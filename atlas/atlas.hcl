env "local" {
  url = "postgres://admin:admin@localhost:5432/auth_main?sslmode=disable&search_path=public"
  src = "file://tables"
  dev = "docker://postgres/17.4/dev?search_path=public"

  migration {
    dir    = "file://migrations"
    format = "atlas"
  }
}
