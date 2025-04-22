-- name: CreateCustomer :one
INSERT INTO customer (customer_name, phone, email)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetCustomerById :one
SELECT * FROM customer
WHERE id = $1
LIMIT 1;

-- name: CreateOrder :one
INSERT INTO "order" (customer_id, total_price)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateOrderStatus :exec
UPDATE "order"
SET order_status = $2
WHERE id = $1;

-- name: DeleteCustomer :exec
DELETE FROM customer 
WHERE id = $1;
