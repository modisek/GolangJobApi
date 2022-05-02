

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "email" varchar NOT NULL,
  "pass" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "jobs" (
  "id" bigserial PRIMARY KEY,
  "title" varchar NOT NULL,
  "description" varchar NOT NULL,
  "email" varchar NOT NULL,
  "type" varchar NOT NULL,
  "category" varchar NOT NULL,
  "location" varchar NOT NULL,
  "expires" timestamptz NOT NULL DEFAULT 'now()',
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "user_applied_job" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigserial NOT NULL
);
CREATE TABLE "user_created_job" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigserial NOT NULL
);

