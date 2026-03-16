SELECT
    id,
    email,
    password_hash,
    role
FROM users
WHERE id = $1;