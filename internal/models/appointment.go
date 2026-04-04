package models

import "time"

type Appointment struct {
	ID                     string     `db:"id" json:"id"`
	PatientID              string     `db:"patient_id" json:"patient_id"`
	ClinicID               string     `db:"clinic_id" json:"clinic_id"`
	DoctorID               string     `db:"doctor_id" json:"doctor_id"`
	ServiceID              string     `db:"service_id" json:"service_id"`
	SlotID                 string     `db:"slot_id" json:"slot_id"`
	Status                 string     `db:"status" json:"status"`
	PaymentWindowExpiresAt *time.Time `db:"payment_window_expires_at" json:"payment_window_expires_at"`
	TotalAmount            float64    `db:"total_amount" json:"total_amount"`
	UserPayAmount          float64    `db:"user_pay_amount" json:"user_pay_amount"`
	InsuredAmount          float64    `db:"insured_amount" json:"insured_amount"`
	CreatedAt              string     `db:"created_at" json:"created_at"`
	UpdatedAt              string     `db:"updated_at" json:"updated_at"`
	StartTime              time.Time  `db:"start_time" json:"start_time"`
	BasePriceAtBooking     float64    `db:"base_price_at_booking" json:"base_price_at_booking"`
	SurchargeAtBooking     float64    `db:"surcharge_at_booking" json:"surcharge_at_booking"`
	TotalPriceAtBooking    float64    `db:"total_price_at_booking" json:"total_price_at_booking"`
	AppliedPolicySnapshot  string     `db:"applied_policy_snapshot" json:"applied_policy_snapshot"`
}

type AppointmentStateHistory struct {
	ID            string  `db:"id" json:"id"`
	AppointmentID string  `db:"appointment_id" json:"appointment_id"`
	FromState     *string `db:"from_state" json:"from_state"`
	ToState       string  `db:"to_state" json:"to_state"`
	ChangedBy     string  `db:"changed_by" json:"changed_by"`
	Reason        *string `db:"reason" json:"reason"`
	CreatedAt     string  `db:"created_at" json:"created_at"`
}
