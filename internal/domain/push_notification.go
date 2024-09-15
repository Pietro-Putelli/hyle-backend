package domain

import "github.com/google/uuid"

type PushNotificationAlert struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type PusNotificationAps struct {
	Alert PushNotificationAlert `json:"alert"`
	Badge int                   `json:"badge"`
}

type PushNotificationPayloadData struct {
	BookID uuid.UUID `json:"bookId"`
}

type PushNotificationPayload struct {
	Aps  PusNotificationAps          `json:"aps"`
	Data PushNotificationPayloadData `json:"data"`
}
