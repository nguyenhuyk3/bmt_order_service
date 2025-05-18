ALTER TABLE "order_seats" DROP CONSTRAINT order_seats_order_id_fkey;
ALTER TABLE "order_fabs" DROP CONSTRAINT order_fabs_order_id_fkey;

DROP TABLE IF EXISTS "order_fabs";
DROP TABLE IF EXISTS "order_seats";
DROP TABLE IF EXISTS "orders";

DROP TYPE IF EXISTS "order_statuses";
