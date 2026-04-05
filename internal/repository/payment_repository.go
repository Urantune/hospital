package repository

import (
	"hospital/internal/config"
	"hospital/internal/models"
)

func CreatePayment(payment *models.Payment) error {
	query := `
	INSERT INTO payments
	(appointment_id, amount, currency, gateway, transaction_code, status, idempotency_key, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	RETURNING id, created_at, updated_at
	`

	return config.DB.QueryRow(
		query,
		payment.AppointmentID,
		payment.Amount,
		payment.Currency,
		payment.Gateway,
		payment.TransactionCode,
		payment.Status,
		payment.IdempotencyKey,
	).Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt)
}

func GetPaymentByID(id string) (*models.Payment, error) {
	var payment models.Payment

	query := `
	SELECT id, appointment_id, amount, currency, gateway, transaction_code, status, idempotency_key, created_at, updated_at
	FROM payments
	WHERE id = $1
	`

	err := config.DB.Get(&payment, query, id)
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func GetLatestPaymentByAppointmentID(appointmentID string) (*models.Payment, error) {
	var payment models.Payment

	query := `
	SELECT id, appointment_id, amount, currency, gateway, transaction_code, status, idempotency_key, created_at, updated_at
	FROM payments
	WHERE appointment_id = $1
	ORDER BY created_at DESC
	LIMIT 1
	`

	err := config.DB.Get(&payment, query, appointmentID)
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func UpdatePaymentStatus(paymentID, status, transactionCode string) error {
	query := `
	UPDATE payments
	SET status = $1,
	    transaction_code = COALESCE(NULLIF($2, ''), transaction_code),
	    updated_at = NOW()
	WHERE id = $3
	`

	_, err := config.DB.Exec(query, status, transactionCode, paymentID)
	return err
}

func CreatePaymentStateHistory(item *models.PaymentStateHistory) error {
	query := `
	INSERT INTO payment_state_history
	(payment_id, from_state, to_state, created_at)
	VALUES ($1, $2, $3, NOW())
	RETURNING id, created_at
	`

	return config.DB.QueryRow(query, item.PaymentID, item.FromState, item.ToState).Scan(&item.ID, &item.CreatedAt)
}

func ListPaymentStateHistory(paymentID string) ([]models.PaymentStateHistory, error) {
	var items []models.PaymentStateHistory

	query := `
	SELECT id, payment_id, from_state, to_state, created_at
	FROM payment_state_history
	WHERE payment_id = $1
	ORDER BY created_at ASC
	`

	err := config.DB.Select(&items, query, paymentID)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func CreatePaymentAttempt(attempt *models.PaymentAttempt) error {
	query := `
	INSERT INTO payment_attempts
	(payment_id, gateway_request_payload, gateway_response_payload, status, created_at)
	VALUES ($1, $2::jsonb, $3::jsonb, $4, NOW())
	RETURNING id, created_at
	`

	return config.DB.QueryRow(
		query,
		attempt.PaymentID,
		attempt.GatewayRequestPayload,
		attempt.GatewayResponsePayload,
		attempt.Status,
	).Scan(&attempt.ID, &attempt.CreatedAt)
}

func UpsertPricingPolicySnapshot(snapshot *models.PricingPolicySnapshot) error {
	query := `
	INSERT INTO pricing_policy_snapshot
	(appointment_id, base_price, extra_fee, total_price, cancellation_policy_snapshot, created_at)
	VALUES ($1, $2, $3, $4, $5::jsonb, NOW())
	ON CONFLICT (appointment_id)
	DO UPDATE SET
		base_price = EXCLUDED.base_price,
		extra_fee = EXCLUDED.extra_fee,
		total_price = EXCLUDED.total_price,
		cancellation_policy_snapshot = EXCLUDED.cancellation_policy_snapshot
	RETURNING id, created_at
	`

	return config.DB.QueryRow(
		query,
		snapshot.AppointmentID,
		snapshot.BasePrice,
		snapshot.ExtraFee,
		snapshot.TotalPrice,
		snapshot.CancellationPolicySnapshot,
	).Scan(&snapshot.ID, &snapshot.CreatedAt)
}

func GetPricingPolicySnapshotByAppointmentID(appointmentID string) (*models.PricingPolicySnapshot, error) {
	var snapshot models.PricingPolicySnapshot

	query := `
	SELECT id, appointment_id, base_price, extra_fee, total_price, cancellation_policy_snapshot::text AS cancellation_policy_snapshot, created_at
	FROM pricing_policy_snapshot
	WHERE appointment_id = $1
	`

	err := config.DB.Get(&snapshot, query, appointmentID)
	if err != nil {
		return nil, err
	}

	return &snapshot, nil
}

func UpsertAppointmentInsuranceSnapshot(snapshot *models.AppointmentInsuranceSnapshot) error {
	query := `
	INSERT INTO appointment_insurance_snapshot
	(appointment_id, plan_id, coverage_rate, coverage_cap, insured_amount, user_pay_amount, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, NOW())
	ON CONFLICT (appointment_id)
	DO UPDATE SET
		plan_id = EXCLUDED.plan_id,
		coverage_rate = EXCLUDED.coverage_rate,
		coverage_cap = EXCLUDED.coverage_cap,
		insured_amount = EXCLUDED.insured_amount,
		user_pay_amount = EXCLUDED.user_pay_amount
	RETURNING id, created_at
	`

	return config.DB.QueryRow(
		query,
		snapshot.AppointmentID,
		snapshot.PlanID,
		snapshot.CoverageRate,
		snapshot.CoverageCap,
		snapshot.InsuredAmount,
		snapshot.UserPayAmount,
	).Scan(&snapshot.ID, &snapshot.CreatedAt)
}

func GetAppointmentInsuranceSnapshotByAppointmentID(appointmentID string) (*models.AppointmentInsuranceSnapshot, error) {
	var snapshot models.AppointmentInsuranceSnapshot

	query := `
	SELECT id, appointment_id, plan_id, coverage_rate, coverage_cap, insured_amount, user_pay_amount, created_at
	FROM appointment_insurance_snapshot
	WHERE appointment_id = $1
	`

	err := config.DB.Get(&snapshot, query, appointmentID)
	if err != nil {
		return nil, err
	}

	return &snapshot, nil
}

func CreatePaymentCallbackReceipt(receipt *models.PaymentCallbackReceipt) error {
	query := `
	INSERT INTO payment_callback_receipts
	(callback_id, payment_id, gateway, transaction_code, status, raw_payload, notes, received_at, processed_at)
	VALUES ($1, $2, $3, $4, $5, $6::jsonb, $7, NOW(), NOW())
	RETURNING id, received_at, processed_at
	`

	return config.DB.QueryRow(
		query,
		receipt.CallbackID,
		receipt.PaymentID,
		receipt.Gateway,
		receipt.TransactionCode,
		receipt.Status,
		receipt.RawPayload,
		receipt.Notes,
	).Scan(&receipt.ID, &receipt.ReceivedAt, &receipt.ProcessedAt)
}

func GetPaymentCallbackReceiptByCallbackID(callbackID string) (*models.PaymentCallbackReceipt, error) {
	var receipt models.PaymentCallbackReceipt

	query := `
	SELECT id, callback_id, payment_id, gateway, transaction_code, status, raw_payload::text AS raw_payload, notes, received_at, processed_at
	FROM payment_callback_receipts
	WHERE callback_id = $1
	`

	err := config.DB.Get(&receipt, query, callbackID)
	if err != nil {
		return nil, err
	}

	return &receipt, nil
}

func PaymentStatus(paymentID string, newStatus string) error {
	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}

	var oldStatus string
	err = tx.QueryRow(`
		SELECT status FROM payments WHERE id=$1
	`, paymentID).Scan(&oldStatus)

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		UPDATE payments
		SET status=$1, updated_at=now()
		WHERE id=$2
	`, newStatus, paymentID)

	if err != nil {
		tx.Rollback()
		return err
	}

	// insert history
	_, err = tx.Exec(`
		INSERT INTO payment_state_history
		(payment_id, from_state, to_state, created_at)
		VALUES ($1,$2,$3,now())
	`, paymentID, oldStatus, newStatus)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
