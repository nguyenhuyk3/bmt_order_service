CREATE TYPE "order_statuses" AS ENUM (
  'created',
  'failed',
  'success'
);

CREATE TABLE "orders" (
  "id" serial PRIMARY KEY NOT NULL,
  "ordered_by" varchar(32) NOT NULL,
  "showtime_id" int NOT NULL,
  "show_date" date NOT NULL,
  "status" order_statuses NOT NULL DEFAULT 'created',
  "note" text NOT NULL DEFAULT '',
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "order_seats" (
  "id" serial PRIMARY KEY,
  "order_id" int NOT NULL,
  "seat_id" int NOT NULL
);

CREATE TABLE "order_fabs" (
  "id" serial PRIMARY KEY,
  "order_id" int NOT NULL,
  "fab_id" int NOT NULL,
  "quantity" int NOT NULL DEFAULT 1
);

CREATE TABLE "outboxes" (
  "id" uuid PRIMARY KEY NOT NULL DEFAULT (gen_random_uuid()),
  "aggregated_type" varchar(64) NOT NULL,
  "aggregated_id" int NOT NULL,
  "event_type" varchar(64) NOT NULL,
  "payload" jsonb NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE INDEX ON "orders" ("ordered_by");

CREATE INDEX ON "orders" ("showtime_id");

CREATE UNIQUE INDEX ON "order_seats" ("order_id", "seat_id");

CREATE INDEX ON "order_fabs" ("order_id", "fab_id");

CREATE INDEX ON "outboxes" ("aggregated_type", "aggregated_id");

CREATE PUBLICATION order_dbz_publication FOR TABLE outboxes;


