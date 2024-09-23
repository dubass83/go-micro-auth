CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "email" varchar UNIQUE NOT NULL,
  "first_name" varchar,
  "last_name" varchar,
  "password" varchar NOT NULL,
  "active" int NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);