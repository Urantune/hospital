package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"hospital/internal/models"
	"hospital/internal/repository"
)

type CreateAppointmentPaymentInput struct {
	AppointmentID              string
	Amount                     float64
	Currency                   string
	Gateway                    string
	TransactionCode            string
	IdempotencyKey             string
	BasePrice                  float64
	ExtraFee                   float64
	CancellationPolicySnapshot json.RawMessage
	PlanID                     *string
	CoverageRate               float64
	CoverageCap                float64
}

type HandlePaymentCallbackInput struct {
	CallbackID      string          `json:"callback_id"`
	PaymentID       string          `json:"payment_id"`
	Gateway         string          `json:"gateway"`
	TransactionCode string          `json:"transaction_code"`
	Status          string          `json:"status"`
	Payload         json.RawMessage `json:"payload"`
}

type AppointmentFinancialSnapshot struct {
	Payment   *models.Payment                      `json:"payment"`
	Pricing   *models.PricingPolicySnapshot        `json:"pricing"`
	Insurance *models.AppointmentInsuranceSnapshot `json:"insurance"`
}

var paymentTransitions = map[string]map[string]bool{
	"initiated": {
		"success": true,
		"failed":  true,
	},
	"success": {
		"partial_refunded": true,
		"refunded":         true,
	},
	"partial_refunded": {
		"refunded": true,
	},
	"failed":   {},
	"refunded": {},
}

func CreateAppointmentPayment(input CreateAppointmentPaymentInput, changedBy, ipAddress string) (*models.Payment, error) {
	if input.AppointmentID == "" || input.Amount <= 0 || input.Currency == "" || input.Gateway == "" {
		return nil, errors.New("appointment_id, amount, currency and gateway are required")
	}

	appointment, err := repository.GetAppointmentByID(input.AppointmentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("appointment not found")
		}
		return nil, err
	}

	payment := &models.Payment{
		AppointmentID:   input.AppointmentID,
		Amount:          input.Amount,
		Currency:        input.Currency,
		Gateway:         input.Gateway,
		TransactionCode: input.TransactionCode,
		Status:          "initiated",
		IdempotencyKey:  input.IdempotencyKey,
	}

	if err := repository.CreatePayment(payment); err != nil {
		return nil, err
	}

	if err := repository.CreatePaymentStateHistory(&models.PaymentStateHistory{
		PaymentID: payment.ID,
		FromState: nil,
		ToState:   payment.Status,
	}); err != nil {
		return nil, err
	}

	cancellationPolicy := input.CancellationPolicySnapshot
	if len(cancellationPolicy) == 0 {
		cancellationPolicy = json.RawMessage(`{}`)
	}

	if err := repository.UpsertPricingPolicySnapshot(&models.PricingPolicySnapshot{
		AppointmentID:              input.AppointmentID,
		BasePrice:                  input.BasePrice,
		ExtraFee:                   input.ExtraFee,
		TotalPrice:                 appointment.TotalAmount,
		CancellationPolicySnapshot: string(cancellationPolicy),
	}); err != nil {
		return nil, err
	}

	if err := repository.UpsertAppointmentInsuranceSnapshot(&models.AppointmentInsuranceSnapshot{
		AppointmentID: input.AppointmentID,
		PlanID:        input.PlanID,
		CoverageRate:  input.CoverageRate,
		CoverageCap:   input.CoverageCap,
		InsuredAmount: appointment.InsuredAmount,
		UserPayAmount: appointment.UserPayAmount,
	}); err != nil {
		return nil, err
	}

	if _, err := TransitionAppointmentStatus(input.AppointmentID, "PENDING_PAYMENT", "Payment initiated", changedBy); err != nil {
		return nil, err
	}

	_ = emitDomainEvent("payment.initiated", "payments", payment.ID, map[string]any{
		"appointment_id": input.AppointmentID,
		"amount":         input.Amount,
		"gateway":        input.Gateway,
	})
	_ = emitDomainEvent("appointment.pending_payment", "appointments", input.AppointmentID, map[string]any{
		"payment_id": payment.ID,
	})

	_ = repository.CreateAuditLog(&models.AuditLog{
		UserID:      changedBy,
		ClinicID:    appointment.ClinicID,
		Action:      "payment.initiated",
		Resource:    "payments",
		ResourceID:  payment.ID,
		Description: "Payment initiated with financial and insurance snapshots",
		IPAddress:   ipAddress,
	})

	_ = repository.CreateAuditLog(&models.AuditLog{
		UserID:      changedBy,
		ClinicID:    appointment.ClinicID,
		Action:      "insurance.snapshot.created",
		Resource:    "appointments",
		ResourceID:  input.AppointmentID,
		Description: "Financial snapshot per appointment created",
		IPAddress:   ipAddress,
	})

	return payment, nil
}

func HandlePaymentCallback(input HandlePaymentCallbackInput, changedBy, ipAddress string) (string, *models.Payment, error) {
	if input.CallbackID == "" || input.PaymentID == "" || input.Status == "" {
		return "", nil, errors.New("callback_id, payment_id and status are required")
	}

	if _, err := repository.GetPaymentCallbackReceiptByCallbackID(input.CallbackID); err == nil {
		payment, paymentErr := repository.GetPaymentByID(input.PaymentID)
		return "duplicate", payment, paymentErr
	}

	payment, err := repository.GetPaymentByID(input.PaymentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil, errors.New("payment not found")
		}
		return "", nil, err
	}

	callbackPayload := string(input.Payload)
	if callbackPayload == "" {
		callbackPayload = `{}`
	}

	resultStatus := "processed"
	if isDelayedOrDuplicatePaymentCallback(payment.Status, input.Status) {
		resultStatus = "ignored_delayed"
		_ = repository.CreatePaymentCallbackReceipt(&models.PaymentCallbackReceipt{
			CallbackID:      input.CallbackID,
			PaymentID:       payment.ID,
			Gateway:         input.Gateway,
			TransactionCode: input.TransactionCode,
			Status:          resultStatus,
			RawPayload:      callbackPayload,
			Notes:           "Callback ignored because current payment state is already newer or terminal",
		})
		return resultStatus, payment, nil
	}

	fromState := payment.Status
	if err := repository.UpdatePaymentStatus(payment.ID, input.Status, input.TransactionCode); err != nil {
		return "", nil, err
	}

	if err := repository.CreatePaymentStateHistory(&models.PaymentStateHistory{
		PaymentID: payment.ID,
		FromState: &fromState,
		ToState:   input.Status,
	}); err != nil {
		return "", nil, err
	}

	_ = repository.CreatePaymentAttempt(&models.PaymentAttempt{
		PaymentID:              payment.ID,
		GatewayRequestPayload:  `{"source":"gateway_callback"}`,
		GatewayResponsePayload: callbackPayload,
		Status:                 input.Status,
	})

	validStatuses := map[string]bool{
		"initiated":        true,
		"success":          true,
		"failed":           true,
		"refunded":         true,
		"partial_refunded": true,
	}

	if !validStatuses[input.Status] {
		return "", nil, errors.New("invalid payment status")
	}
	if !canTransitionPayment(payment.Status, input.Status) {
		return "", nil, errors.New("invalid payment state transition")
	}

	_ = repository.CreatePaymentCallbackReceipt(&models.PaymentCallbackReceipt{
		CallbackID:      input.CallbackID,
		PaymentID:       payment.ID,
		Gateway:         input.Gateway,
		TransactionCode: input.TransactionCode,
		Status:          resultStatus,
		RawPayload:      callbackPayload,
		Notes:           "Callback applied",
	})
	appointment, err := repository.GetAppointmentByID(payment.AppointmentID)
	if err == nil {
		switch input.Status {
		case "success":
			_, _ = TransitionAppointmentStatus(payment.AppointmentID, "CONFIRMED", "Payment callback success", changedBy)
			_ = emitDomainEvent("payment.succeeded", "payments", payment.ID, map[string]any{
				"appointment_id":   payment.AppointmentID,
				"transaction_code": input.TransactionCode,
			})
			_ = emitDomainEvent("appointment.confirmed", "appointments", payment.AppointmentID, map[string]any{
				"payment_id": payment.ID,
			})
		case "failed":
			_ = emitDomainEvent("payment.failed", "payments", payment.ID, map[string]any{
				"appointment_id":   payment.AppointmentID,
				"transaction_code": input.TransactionCode,
			})
		case "partial_refunded", "refunded":
			_ = emitDomainEvent("payment.refunded", "payments", payment.ID, map[string]any{
				"appointment_id": payment.AppointmentID,
				"refund_status":  input.Status,
			})
		}

		_ = repository.CreateAuditLog(&models.AuditLog{
			UserID:      changedBy,
			ClinicID:    appointment.ClinicID,
			Action:      "payment.callback." + input.Status,
			Resource:    "payments",
			ResourceID:  payment.ID,
			Description: "Payment callback processed with appointment lifecycle enforcement",
			IPAddress:   ipAddress,
		})

		_ = repository.CreateAuditLog(&models.AuditLog{
			UserID:      changedBy,
			ClinicID:    appointment.ClinicID,
			Action:      "insurance.audit",
			Resource:    "appointments",
			ResourceID:  appointment.ID,
			Description: "Insurance and payment callback audit recorded",
			IPAddress:   ipAddress,
		})
	}

	updatedPayment, err := repository.GetPaymentByID(payment.ID)
	if err != nil {
		return "", nil, err
	}

	return resultStatus, updatedPayment, nil
}

func GetAppointmentFinancialSnapshot(appointmentID string) (*AppointmentFinancialSnapshot, error) {
	payment, _ := repository.GetLatestPaymentByAppointmentID(appointmentID)
	pricing, _ := repository.GetPricingPolicySnapshotByAppointmentID(appointmentID)
	insurance, _ := repository.GetAppointmentInsuranceSnapshotByAppointmentID(appointmentID)

	if payment == nil && pricing == nil && insurance == nil {
		return nil, errors.New("financial snapshot not found")
	}

	return &AppointmentFinancialSnapshot{
		Payment:   payment,
		Pricing:   pricing,
		Insurance: insurance,
	}, nil
}

func ListPaymentHistory(paymentID string) ([]models.PaymentStateHistory, error) {
	return repository.ListPaymentStateHistory(paymentID)
}

func ListAppointmentDomainEvents(appointmentID string) ([]models.DomainEvent, error) {
	return repository.ListDomainEventsByAggregateID(appointmentID)
}

func canTransitionPayment(fromState, toState string) bool {
	if fromState == toState {
		return true
	}

	allowedNext, ok := paymentTransitions[fromState]
	if !ok {
		return false
	}

	return allowedNext[toState]
}

func isDelayedOrDuplicatePaymentCallback(currentStatus, incomingStatus string) bool {
	if currentStatus == incomingStatus {
		return true
	}

	if !canTransitionPayment(currentStatus, incomingStatus) {
		switch currentStatus {
		case "success", "partial_refunded", "refunded":
			return true
		}
	}

	return false
}

func emitDomainEvent(eventType, aggregateType, aggregateID string, payload any) error {
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return repository.CreateDomainEvent(&models.DomainEvent{
		EventType:     eventType,
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		Payload:       string(rawPayload),
		Status:        "pending",
	})
}
