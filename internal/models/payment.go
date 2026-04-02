package models

type Payment struct {
	ID              string  `db:"id" json:"id"`
	AppointmentID   string  `db:"appointment_id" json:"appointment_id"`
	Amount          float64 `db:"amount" json:"amount"`
	Currency        string  `db:"currency" json:"currency"`
	Gateway         string  `db:"gateway" json:"gateway"`
	TransactionCode string  `db:"transaction_code" json:"transaction_code"`
	Status          string  `db:"status" json:"status"`
	IdempotencyKey  string  `db:"idempotency_key" json:"idempotency_key"`
	CreatedAt       string  `db:"created_at" json:"created_at"`
	UpdatedAt       string  `db:"updated_at" json:"updated_at"`
}

type PaymentStateHistory struct {
	ID        string  `db:"id" json:"id"`
	PaymentID string  `db:"payment_id" json:"payment_id"`
	FromState *string `db:"from_state" json:"from_state"`
	ToState   string  `db:"to_state" json:"to_state"`
	CreatedAt string  `db:"created_at" json:"created_at"`
}

type PricingPolicySnapshot struct {
	ID                         string  `db:"id" json:"id"`
	AppointmentID              string  `db:"appointment_id" json:"appointment_id"`
	BasePrice                  float64 `db:"base_price" json:"base_price"`
	ExtraFee                   float64 `db:"extra_fee" json:"extra_fee"`
	TotalPrice                 float64 `db:"total_price" json:"total_price"`
	CancellationPolicySnapshot string  `db:"cancellation_policy_snapshot" json:"cancellation_policy_snapshot"`
	CreatedAt                  string  `db:"created_at" json:"created_at"`
}

type AppointmentInsuranceSnapshot struct {
	ID            string  `db:"id" json:"id"`
	AppointmentID string  `db:"appointment_id" json:"appointment_id"`
	PlanID        *string `db:"plan_id" json:"plan_id"`
	CoverageRate  float64 `db:"coverage_rate" json:"coverage_rate"`
	CoverageCap   float64 `db:"coverage_cap" json:"coverage_cap"`
	InsuredAmount float64 `db:"insured_amount" json:"insured_amount"`
	UserPayAmount float64 `db:"user_pay_amount" json:"user_pay_amount"`
	CreatedAt     string  `db:"created_at" json:"created_at"`
}

type PaymentAttempt struct {
	ID                     string `db:"id" json:"id"`
	PaymentID              string `db:"payment_id" json:"payment_id"`
	GatewayRequestPayload  string `db:"gateway_request_payload" json:"gateway_request_payload"`
	GatewayResponsePayload string `db:"gateway_response_payload" json:"gateway_response_payload"`
	Status                 string `db:"status" json:"status"`
	CreatedAt              string `db:"created_at" json:"created_at"`
}

type PaymentCallbackReceipt struct {
	ID              string `db:"id" json:"id"`
	CallbackID      string `db:"callback_id" json:"callback_id"`
	PaymentID       string `db:"payment_id" json:"payment_id"`
	Gateway         string `db:"gateway" json:"gateway"`
	TransactionCode string `db:"transaction_code" json:"transaction_code"`
	Status          string `db:"status" json:"status"`
	RawPayload      string `db:"raw_payload" json:"raw_payload"`
	Notes           string `db:"notes" json:"notes"`
	ReceivedAt      string `db:"received_at" json:"received_at"`
	ProcessedAt     string `db:"processed_at" json:"processed_at"`
}
