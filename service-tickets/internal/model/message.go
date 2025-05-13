package model

import "time"

// Message represents a message in a ticket
type Message struct {
	ID          string    `json:"id"`
	TicketID    string    `json:"ticketId"`
	FromAddress string    `json:"fromAddress"`
	ToAddress   string    `json:"toAddress"`
	Subject     string    `json:"subject"`
	Body        string    `json:"body"`
	Channel     string    `json:"channel"`
	CreatedAt   time.Time `json:"createdAt"`
}
