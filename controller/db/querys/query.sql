-- name: TotalPriceSubscriptions :one
SELECT
   SUM(sser.price) as "total_price"
FROM services.subscription ssub
LEFT JOIN services.services sser on sser.uuid = ssub.service_uuid
WHERE
    (ssub.start_date >= $1::DATE) AND
    (ssub.stop_date <= $2::DATE) AND
    CASE WHEN $3::VARCHAR LIKE '' THEN 1=1 ELSE ssub.user_uuid LIKE $3::VARCHAR END 
    AND
    CASE WHEN $4::VARCHAR LIKE '' THEN 1=1 ELSE sser.name LIKE $4::VARCHAR END;

-- name: ListSubscriptions :many
SELECT
    ssub.user_uuid as "user_uuid",
    sser.name as "service_name",
    sser.price as "service_price",
    ssub.start_date,
    ssub.stop_date
FROM services.subscription ssub
LEFT JOIN services.services sser on sser.uuid = ssub.service_uuid
WHERE
    CASE WHEN $1::VARCHAR LIKE '' THEN 1=1 ELSE ssub.user_uuid LIKE $1::VARCHAR END
    AND
    CASE WHEN $2::VARCHAR LIKE '' THEN 1=1 ELSE sser.name LIKE $2::VARCHAR END;

-- name: ExistingSubscriptionService :one
SELECT EXISTS(
    SELECT
        1
    FROM services.subscription sub
    LEFT JOIN services.services ser ON ser.uuid = sub.service_uuid
    WHERE
        sub.user_uuid = $1 AND
        ser.name = $2
);

-- name: ExistingService :one
SELECT
    ser.uuid
FROM services.services ser
WHERE ser.name = $1;

-- name: AddService :one
INSERT INTO services.services (uuid, name, price)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteService :exec
DELETE FROM services.services
WHERE uuid = $1;

-- name: AddSubscription :one
INSERT INTO services.subscription (user_uuid, service_uuid, start_date, stop_date)
VALUES ($1, $2, $3, $4)
RETURNING *;