package activity

import (
	"github.com/aimustaev/service-workflow/internal/generated/proto"
)

// Client represents a ticket service client
type Activity struct {
	ticketClient proto.TicketServiceClient
}

func NewActivity(ticketClient proto.TicketServiceClient) *Activity {
	return &Activity{
		ticketClient: ticketClient,
	}
}
