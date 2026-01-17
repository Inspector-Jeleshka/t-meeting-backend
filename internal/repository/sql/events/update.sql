UPDATE events
SET name = $2,
    metadata = $3::jsonb,
		    content = $4::jsonb,
		    status = $5
WHERE id = $1