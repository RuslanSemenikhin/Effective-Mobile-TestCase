-- name: ListSubscriptions :many
SELECT
    ssub.user_uuid as "user_uuid",
    sser.name as "service_name",
    SUM(sser.price) as "service_price"
FROM services.subscription ssub
LEFT JOIN services.services sser on sser.uuid = ssub.service_uuid
WHERE
    ($1::DATE IS NULL OR ssub.start_date >= $1) AND
    ($2::DATE IS NULL OR ssub.stop_date <= $2) AND
    ($3::VARCHAR IS NULL OR ssub.user_uuid = $3) AND
    ($4::VARCHAR IS NULL OR sser.name = $4)
GROUP BY ssub.user_uuid, sser.name;

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