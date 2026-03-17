package models

type CMSChangeEvent struct {
	ID           string `db:"id" json:"id"`
	EventID      string `db:"event_id" json:"event_id"`
	Source       string `db:"source" json:"source"`
	EntityType   string `db:"entity_type" json:"entity_type"`
	EntityID     string `db:"entity_id" json:"entity_id"`
	Action       string `db:"action" json:"action"`
	Payload      string `db:"payload" json:"payload"`
	Status       string `db:"status" json:"status"`
	ErrorMessage string `db:"error_message" json:"error_message"`
	ProcessedBy  string `db:"processed_by" json:"processed_by"`
	ProcessedAt  string `db:"processed_at" json:"processed_at"`
	CreatedAt    string `db:"created_at" json:"created_at"`
}
