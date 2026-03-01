package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type NewEvent struct {
	Name     string         `json:"name"`
	Metadata EventMetadata  `json:"metadata"`
	Content  []ContentBlock `json:"content"`
	Status   EventStatus    `json:"status"`
}

type Event struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	Metadata  EventMetadata  `json:"metadata"`
	Content   []ContentBlock `json:"content"`
	Status    EventStatus    `json:"status"` // "draft"/"published" и тд
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
}

type EventStatus string

const (
	EventStatusDraft     EventStatus = "draft"
	EventStatusPublished EventStatus = "published"
	EventStatusArchived  EventStatus = "archived"
)

type EventMetadata struct {
	Datetime time.Time `json:"datetime"`         // "2025-11-23T22:00:00Z"
	Location string    `json:"location"`         // "Общежитие нгту 10, комната 1004-2"
	Reason   string    `json:"reason,omitempty"` // "мой день рождения"
}

type ContentBlock struct {
	Block   string          `json:"block"`   // "promo-text"/"map"/"timeline"
	Payload json.RawMessage `json:"payload"` // json, чтобы уходил в бд как есть. Если захотим распарсим в другое
}

type PromoTextPayload []string // promo-text когда нужен будет

type Point struct {
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Text string  `json:"text"`
}

type MapPayload struct {
	Background string  `json:"background"`
	Points     []Point `json:"points"`
}

type TimelineItem struct { //легендарный таймлайн
	Name string `json:"name"`
	Time string `json:"time"`
}

type TimelineItems []TimelineItem
