package service

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"hospital/internal/repository"

	"github.com/google/uuid"
)

const MaxRetries = 3

func QueueNotification(notifType, recipient, content, referenceID string) error {

	exists, err := repository.CheckNotificationExists(referenceID, notifType)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("[Notification] Skipped duplicate notification: %s for reference: %s\n", notifType, referenceID)
		return nil
	}

	notif := &repository.Notification{
		ID:          uuid.New().String(),
		Type:        notifType,
		Recipient:   recipient,
		Content:     content,
		ReferenceID: referenceID,
	}

	log.Printf("[Notification] Queued notification for reference: %s\n", referenceID)
	return repository.CreateNotification(notif)
}

func StartNotificationWorker() {
	log.Println("[Worker] Starting Notification Background Worker...")

	go func() {
		for {

			notifs, err := repository.GetPendingNotifications(MaxRetries)
			if err == nil && len(notifs) > 0 {
				for _, n := range notifs {

					errSend := mockSendExternalAPI(n)

					if errSend != nil {

						n.RetryCount++
						log.Printf("[Worker] Failed to send %s (Retry %d/%d): %v\n", n.ReferenceID, n.RetryCount, MaxRetries, errSend)
						repository.UpdateNotificationStatus(n.ID, "failed", n.RetryCount)
					} else {

						log.Printf("[Worker] Successfully sent: %s\n", n.ReferenceID)
						repository.UpdateNotificationStatus(n.ID, "sent", n.RetryCount)
					}
				}
			}

			time.Sleep(5 * time.Second)
		}
	}()
}

func mockSendExternalAPI(n repository.Notification) error {

	if rand.Float32() < 0.3 {
		return fmt.Errorf("timeout connecting to external provider")
	}
	return nil
}
