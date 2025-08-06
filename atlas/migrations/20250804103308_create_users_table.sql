-- Create "users" table
CREATE TABLE "users" (
  "id" bigserial NOT NULL,
  "public_id" uuid NOT NULL,
  "email" character varying(255) NOT NULL,
  "password_hash" text NOT NULL,
  "name" character varying(255) NOT NULL,
  "role" character varying(50) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  "deleted_at" timestamp NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "users_email_unique" UNIQUE ("email"),
  CONSTRAINT "users_public_id_unique" UNIQUE ("public_id")
);
