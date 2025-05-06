export * from "./client";
export * from "./activities/process-data";
export * from "./workflows";
export * from "./worker";
import { createNewTicket, addMessageToTicket, classifyTicket, assignTicket, closeTicket, findTicketBySender, markAsResolved } from "./activities/ticket";

export const ticketActivities = {
  findTicketBySender,
  createNewTicket,
  addMessageToTicket,
  classifyTicket,
  assignTicket,
  closeTicket,
  markAsResolved
}; 