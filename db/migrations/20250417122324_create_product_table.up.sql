CREATE TABLE "product" (
  "id" VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
  "product_name" VARCHAR(36) NOT NULL,
  "description" VARCHAR(35) NOT NULL,
  "product_price"   NUMERIC(10,2) NOT NULL CHECK(product_price >= 0),
  "created_at" TIMESTAMP DEFAULT now()
);
