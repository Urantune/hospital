package service

import "hospital/internal/repository"

func RunReconciliation() error {

	payments, err := repository.GetAllTodayPayments()
	if err != nil {
		return err
	}
	mismatchCount := 0
	for _, p := range payments {
		gatewayStatus := mockGatewayCheck(p.TransactionCode)
		if gatewayStatus != p.Status {
			mismatchCount++
			err := repository.PaymentStatus(p.ID, gatewayStatus)
			if err != nil {
				return err
			}
		}
	}
	return repository.InsertReconciliationReport("mock_gateway", mismatchCount)
}

func mockGatewayCheck(txCode string) string {

	return "success"
}
