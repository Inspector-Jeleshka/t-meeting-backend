SELECT
    id,
    email,
    password_hash,
    role
FROM users
WHERE email = $1;
