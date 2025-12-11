package repository

import (
	"context"
	"encoding/json"
	"errors"
	"t-meeting-backend/internal/domain"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository interface {
	Create(ctx context.Context, e *domain.Event) error
	GetAll(ctx context.Context) ([]*domain.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	Update(ctx context.Context, id uuid.UUID, e *domain.Event) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// ин мемори реализация (для тестов)
type mapEventRepository struct {
	database map[uuid.UUID]*domain.Event
}

// Реализация поверх постгре
type PgxEventRepository struct {
	pool *pgxpool.Pool
}

func NewPgxEventRepository(pool *pgxpool.Pool) *PgxEventRepository {
	return &PgxEventRepository{pool: pool}
}

// pgxEventRepository — работа с БД

func (erep *PgxEventRepository) Create(ctx context.Context, e *domain.Event) error {
	id := uuid.New()
	e.ID = id
	if e.Status == "" {
		e.Status = "draft"
	}

	metaBytes, err := json.Marshal(e.Metadata)
	if err != nil {
		return err
	}

	contentBytes, err := json.Marshal(e.Content)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = erep.pool.Exec(ctx, `
		INSERT INTO events (id, name, metadata, content, status)
		VALUES ($1, $2, $3::jsonb, $4::jsonb, COALESCE($5, 'draft'))
	`, e.ID, e.Name, metaBytes, contentBytes, e.Status)
	return err
}

func (erep *PgxEventRepository) GetAll(ctx context.Context) ([]*domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := erep.pool.Query(ctx, `
		SELECT id, name, metadata, content, status, created_at, updated_at
		FROM events
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*domain.Event

	for rows.Next() {
		var (
			e            domain.Event
			metaBytes    []byte
			contentBytes []byte
		)

		if err := rows.Scan(
			&e.ID,
			&e.Name,
			&metaBytes,
			&contentBytes,
			&e.Status,
			&e.CreatedAt,
			&e.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(metaBytes, &e.Metadata); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(contentBytes, &e.Content); err != nil {
			return nil, err
		}

		res = append(res, &e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (erep *PgxEventRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var (
		e            domain.Event
		metaBytes    []byte
		contentBytes []byte
	)

	err := erep.pool.QueryRow(ctx, `
		SELECT id, name, metadata, content, status, created_at, updated_at
		FROM events
		WHERE id = $1
	`, id).Scan(
		&e.ID,
		&e.Name,
		&metaBytes,
		&contentBytes,
		&e.Status,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(metaBytes, &e.Metadata); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(contentBytes, &e.Content); err != nil {
		return nil, err
	}

	return &e, nil
}

func (erep *PgxEventRepository) Update(ctx context.Context, id uuid.UUID, e *domain.Event) error {
	metaBytes, err := json.Marshal(e.Metadata)
	if err != nil {
		return err
	}
	contentBytes, err := json.Marshal(e.Content)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cmdTag, err := erep.pool.Exec(ctx, `
		UPDATE events
		SET name = $2,
		    metadata = $3::jsonb,
		    content = $4::jsonb,
		    status = $5
		WHERE id = $1
	`, id, e.Name, metaBytes, contentBytes, e.Status)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("event not found")
	}
	return nil
}

func (erep *PgxEventRepository) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cmdTag, err := erep.pool.Exec(ctx, `DELETE FROM events WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("event not found")
	}
	return nil
}

//
// mapEventRepository — in-memory реализация

func (m *mapEventRepository) Create(_ context.Context, e *domain.Event) error {
	id := uuid.New()
	e.ID = id
	m.database[id] = e
	return nil
}

func (m *mapEventRepository) GetAll(_ context.Context) ([]*domain.Event, error) {
	res := make([]*domain.Event, 0, len(m.database))
	for _, v := range m.database {
		res = append(res, v)
	}
	return res, nil
}

func (m *mapEventRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.Event, error) {
	event, ok := m.database[id]
	if !ok {
		return nil, errors.New("мероприятие не найдено")
	}
	return event, nil
}

func (m *mapEventRepository) Update(_ context.Context, id uuid.UUID, e *domain.Event) error {
	if _, ok := m.database[id]; !ok {
		return errors.New("мероприятие не найдено")
	}
	e.ID = id
	m.database[id] = e
	return nil
}

func (m *mapEventRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(m.database, id)
	return nil
}
