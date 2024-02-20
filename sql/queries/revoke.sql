-- name: RevokeToken :exec
INSERT INTO revocation(token, revoked_at) VALUES($1, $2);

-- name: VerifyRefresh :one
SELECT EXISTS (
    SELECT 1 FROM revocation WHERE token ==$1
) AS value_exists;
