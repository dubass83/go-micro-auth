-- SQL dump generated using DBML (dbml.dbdiagram.io)
-- Database: PostgreSQL
-- Generated at: 2024-09-02T18:48:51.328Z

CREATE TABLE "user" (
  "id" int PRIMARY KEY,
  "email" varchar UNIQUE NOT NULL,
  "first_name" varchar,
  "last_name" varchar,
  "password" varchar NOT NULL,
  "active" int NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
