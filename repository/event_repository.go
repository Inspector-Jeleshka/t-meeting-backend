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

var dbpool *pgxpool.Pool

func init() {
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

	dbpool = pool
	log.Println("Connected to postgres (repository)")
}

type EventRepository interface {
	Create(c context.Context, e *domain.Event) error
	GetAll(c context.Context) ([]*domain.Event, error)
	GetByID(c context.Context, id uuid.UUID) (*domain.Event, error)
	Update(c context.Context, id uuid.UUID, e *domain.Event) error
	Delete(c context.Context, id uuid.UUID) error
}

type eventRepository struct {
	database map[uuid.UUID]*domain.Event
}

func NewEventRepository() EventRepository {
	db := make(map[uuid.UUID]*domain.Event)
	return &eventRepository{database: db}
}

func (er *eventRepository) Create(ctx context.Context, e *domain.Event) error {
	id := uuid.New()
	e.ID = id
	er.database[id] = e
	if dbpool == nil {
		return nil
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

	_, err = dbpool.Exec(ctx, `
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

func (er *eventRepository) GetAll(_ context.Context) ([]*domain.Event, error) {
	var res []*domain.Event
	for _, v := range er.database {
		res = append(res, v)
	}
	return res, nil
}

func (er *eventRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.Event, error) {
	event := er.database[id]
	if event == nil {
		return nil, errors.New("мероприятие не найдено")
	}
	return er.database[id], nil
}

func (er *eventRepository) Update(_ context.Context, id uuid.UUID, e *domain.Event) error {
	e.ID = id
	er.database[id] = e
	return nil
}

func (er *eventRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(er.database, id)
	return nil
}
