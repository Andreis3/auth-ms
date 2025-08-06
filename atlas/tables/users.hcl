schema "public" {}

table "users" {
  schema = schema.public
  column "id" {
    type     = bigserial
    null     = false
  }
  column "public_id" {
    type     = uuid
    null     = false
  }
  column "email" {
    type     = varchar(255)
    null     = false
  }
  column "password_hash" {
    type     = text
    null     = false
  }
  column "name" {
    type     = varchar(255)
    null     = false
  }
  column "role" {
    type     = varchar(50)
    null     = false
  }
  column "created_at" {
    type     = timestamp
    default  = sql("now()")
    null     = false
  }
  column "updated_at" {
    type     = timestamp
    default  = sql("now()")
    null     = false
  }
  column "deleted_at" {
    type = timestamp
    null = true
  }

  primary_key {
    columns = [column.id]
  }

  unique "users_public_id_unique" {
    columns = [column.public_id]
  }

  unique "users_email_unique" {
    columns = [column.email]
  }


}
