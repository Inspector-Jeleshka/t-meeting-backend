INSERT INTO events (id, name, metadata, content, status)
VALUES ($1, $2, $3::jsonb, $4::jsonb, COALESCE($5, 'draft'))