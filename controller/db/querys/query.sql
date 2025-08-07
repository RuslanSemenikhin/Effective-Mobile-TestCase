-- name: GetAllSubscriptions :many
SELECT
    ssub.user_uuid as "user_uuid",
    sser.name as "service_name",
    SUM(sser.price) as "service_price"
FROM services.subscription ssub
LEFT JOIN services.services sser on sser.uuid = ssub.service_uuid
WHERE
    ssub.start_date >= $1 AND
    ssub.stop_date <= $2 AND
    ($3::VARCHAR IS NULL OR ssub.user_uuid = $3) AND
    ($4::VARCHAR IS NULL OR sser.name = $4)
GROUP BY ssub.user_uuid, sser.name;