package models

type DomainEvent struct {
	ID            string `db:"id" json:"id"`
	EventType     string `db:"event_type" json:"event_type"`
	AggregateType string `db:"aggregate_type" json:"aggregate_type"`
	AggregateID   string `db:"aggregate_id" json:"aggregate_id"`
	Payload       string `db:"payload" json:"payload"`
	Status        string `db:"status" json:"status"`
	CreatedAt     string `db:"created_at" json:"created_at"`
}
