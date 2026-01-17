SELECT id, name, metadata, content, status, created_at, updated_at
FROM events
WHERE id = $1