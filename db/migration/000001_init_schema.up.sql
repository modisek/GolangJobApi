

CREATE TYPE "roles" AS ENUM (
  'admin',
  'user'
);

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "password" varchar NOT NULL,
  "role" roles DEFAULT 'user',
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "jobs" (
  "id" bigserial PRIMARY KEY,
  "title" varchar NOT NULL,
  "descrption" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "jobs_users" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigserial NOT NULL
);

ALTER TABLE "jobs_users" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
