package repository

import (
	"hospital/internal/config"
)

// Khai báo cấu trúc Notification
type Notification struct {
	ID          string `db:"id" json:"id"`
	Type        string `db:"type" json:"type"`                 // EMAIL, SMS
	Recipient   string `db:"recipient" json:"recipient"`       // Địa chỉ nhận
	Content     string `db:"content" json:"content"`           // Nội dung
	ReferenceID string `db:"reference_id" json:"reference_id"` // Dùng để CHỐNG TRÙNG (VD: PaymentID)
	Status      string `db:"status" json:"status"`             // pending, sent, failed
	RetryCount  int    `db:"retry_count" json:"retry_count"`   // Đếm số lần gửi lại
}

// 1. Kiểm tra chống trùng lặp (De-duplication)
func CheckNotificationExists(referenceID, notifType string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM notifications WHERE reference_id = $1 AND type = $2`
	err := config.DB.Get(&count, query, referenceID, notifType)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// 2. Tạo thông báo mới (vào hàng đợi)
func CreateNotification(n *Notification) error {
	query := `
	INSERT INTO notifications (id, type, recipient, content, reference_id, status, retry_count)
	VALUES ($1, $2, $3, $4, $5, 'pending', 0)
	`
	_, err := config.DB.Exec(query, n.ID, n.Type, n.Recipient, n.Content, n.ReferenceID)
	return err
}

// 3. Lấy các thông báo đang chờ hoặc bị lỗi (để Retry)
func GetPendingNotifications(maxRetries int) ([]Notification, error) {
	var notifs []Notification
	// Lấy những thằng 'pending' HOẶC 'failed' nhưng số lần thử lại vẫn còn < maxRetries
	query := `
	SELECT id, type, recipient, content, reference_id, status, retry_count 
	FROM notifications 
	WHERE status = 'pending' OR (status = 'failed' AND retry_count < $1)
	LIMIT 50
	`
	err := config.DB.Select(&notifs, query, maxRetries)
	return notifs, err
}

// 4. Cập nhật trạng thái sau khi gửi
func UpdateNotificationStatus(id, status string, retryCount int) error {
	query := `UPDATE notifications SET status = $1, retry_count = $2 WHERE id = $3`
	_, err := config.DB.Exec(query, status, retryCount, id)
	return err
}
