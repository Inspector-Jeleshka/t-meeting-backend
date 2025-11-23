package repository

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"t-meeting-backend/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository interface {
	Create(c context.Context, e *domain.Event) error
	GetAll(c context.Context) ([]*domain.Event, error)
	GetByID(c context.Context, id uuid.UUID) (*domain.Event, error)
	Update(c context.Context, id uuid.UUID, e *domain.Event) error
	Delete(c context.Context, id uuid.UUID) error
}

type mapEventRepository struct {
	database map[uuid.UUID]*domain.Event
}

type pgxEventRepository struct {
	database *pgxpool.Pool
}

func NewMapEventRepository() EventRepository {
	db := make(map[uuid.UUID]*domain.Event)
	return &mapEventRepository{database: db}
}

func NewEventRepository() EventRepository {
	dsn := "postgres://postgres:coolpassword@localhost:5433/tmeeting?sslmode=disable"

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("pgx parse config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		log.Fatalf("pgx new pool: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("pgx ping: %v", err)
	}

	log.Println("Connected to postgres (repository)")

	return &pgxEventRepository{database: pool}
}

func (erep *pgxEventRepository) Create(ctx context.Context, e *domain.Event) error {
	id := uuid.New()
	e.ID = id

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

	_, err = erep.database.Exec(ctx, `
    INSERT INTO events (id, name, metadata, content, status)
    VALUES ($1, $2, $3::jsonb, $4::jsonb, COALESCE($5, 'draft'))
`,
		e.ID,
		e.Name,
		metaBytes,
		contentBytes,
		e.Status,
	)
	return err
}

func (erep *pgxEventRepository) GetAll(ctx context.Context) ([]*domain.Event, error) {
	return nil, nil
}

func (erep *pgxEventRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	return nil, nil
}

func (erep *pgxEventRepository) Update(ctx context.Context, id uuid.UUID, event *domain.Event) error {
	return nil
}

func (erep *pgxEventRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (erep *mapEventRepository) Create(_ context.Context, event *domain.Event) error {
	id := uuid.New()
	event.ID = id
	erep.database[id] = event
	return nil
}

func (erep *mapEventRepository) GetAll(_ context.Context) ([]*domain.Event, error) {
	var res []*domain.Event
	for _, v := range erep.database {
		res = append(res, v)
	}
	return res, nil
}

func (erep *mapEventRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.Event, error) {
	event := erep.database[id]
	if event == nil {
		return nil, errors.New("мероприятие не найдено")
	}
	return erep.database[id], nil
}

func (erep *mapEventRepository) Update(_ context.Context, id uuid.UUID, event *domain.Event) error {
	event.ID = id
	erep.database[id] = event
	return nil
}

func (erep *mapEventRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(erep.database, id)
	return nil
}
