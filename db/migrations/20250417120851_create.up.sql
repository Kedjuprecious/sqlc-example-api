CREATE TABLE "order" (
  "id" VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
  "customer_id" VARCHAR(36) NOT NULL,
  "order_status" TEXT NOT NULL DEFAULT 'Pending',
  "total_price" VARCHAR(10) NOT NULL,
  "created_at" TIMESTAMP DEFAULT now(),
  CONSTRAINT fk_customer FOREIGN KEY ("customer_id") REFERENCES "customer" ("id") ON DELETE CASCADE
);
