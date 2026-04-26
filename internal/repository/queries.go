package repository

const (
	qEventCreate = `INSERT INTO events (id, name, metadata, content, status)
VALUES ($1, $2, $3::jsonb, $4::jsonb, COALESCE($5, 'draft'))`
	qEventGetAll = `SELECT id, name, metadata, content, status, created_at, updated_at
FROM events
ORDER BY created_at DESC`
	qEventGetByID = `SELECT id, name, metadata, content, status, created_at, updated_at
FROM events
WHERE id = $1`
	qEventUpdate = `UPDATE events
SET name = $2,
	metadata = $3::jsonb,
	content = $4::jsonb,
	status = $5
WHERE id = $1`
	qEventDelete = `DELETE FROM events WHERE id = $1`
	qUserCreate  = `INSERT INTO users (id, email, password_hash, role)
VALUES ($1, $2, $3, $4)`
	qUserGetByEmail = `SELECT id, email, password_hash, role
FROM users
WHERE email = $1`
	qUserGetByID = `SELECT id, email, password_hash, role
FROM users
WHERE id = $1`
)
