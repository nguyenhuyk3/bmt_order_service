CREATE TYPE "order_statuses" AS ENUM (
  'pending',
  'canceled',
  'failed',
  'success'
);

CREATE TABLE "orders" (
  "id" serial PRIMARY KEY NOT NULL,
  "ordered_by" varchar(32) NOT NULL,
  "showtime_id" int NOT NULL,
  "show_date" date NOT NULL,
  "status" order_statuses NOT NULL DEFAULT 'pending',
  "note" text NOT NULL DEFAULT '',
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "order_seats" (
  "id" serial PRIMARY KEY,
  "order_id" int NOT NULL,
  "seat_id" int NOT NULL,
  "price" int NOT NULL
);

CREATE TABLE "order_fabs" (
  "id" serial PRIMARY KEY,
  "order_id" int NOT NULL,
  "fab_id" int NOT NULL,
  "quantity" int NOT NULL DEFAULT 1,
  "price" int NOT NULL
);

CREATE INDEX ON "orders" ("ordered_by");

CREATE INDEX ON "orders" ("showtime_id");

CREATE UNIQUE INDEX ON "order_seats" ("order_id", "seat_id");

CREATE INDEX ON "order_fabs" ("order_id", "fab_id");

ALTER TABLE "order_seats" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "order_fabs" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");
