package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	Metadata  EventMetadata  `json:"metadata"`
	Content   []ContentBlock `json:"content"`
	Status    string         `json:"status"` // "draft"/"published" и тд
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
}

type EventMetadata struct {
	Date     string `json:"date"`             // "2025-11-23"
	Time     string `json:"time"`             // "22:00"
	Location string `json:"location"`         // "Общежитие нгту 10, комната 1004-2"
	Reason   string `json:"reason,omitempty"` // "мой день рождения"
}

type ContentBlock struct {
	Block   string          `json:"block"`   // "promo-text"/"map"/"timeline"
	Payload json.RawMessage `json:"payload"` // json, чтобы уходил в бд как есть. Если захотим распарсим в другое
}

type PromoTextPayload []string // promo-text когда нужен будет

type MapPayload struct {
	Lon   float64 `json:"lon"` // map примерно так должен описываться, тоже потом пригодится
	Lat   float64 `json:"lat"`
	Title string  `json:"title"`
	Icon  string  `json:"icon"`
}

type TimelineItem struct { //легендарный таймлайн
	Name string `json:"name"`
	Time string `json:"time"`
}

type TimelineItems []TimelineItem
