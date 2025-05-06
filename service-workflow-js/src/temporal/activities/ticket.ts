import { logger } from "../../utils";

interface Message {
  id: string;
  text: string;
  timestamp: string;
  [key: string]: any;
}

interface Ticket {
  id: string;
  status: 'open' | 'in_progress' | 'closed' | 'resolved';
  customerId: string;
  channel: string;
  messages: Message[];
  previousTicketId?: string;
  assignedTo?: string;
}

interface TicketParams {
  customerId: string;
  channel: string;
  initialMessage: Message;
  previousTicketId?: string;
}

// Имитация базы данных
const ticketsDB: Record<string, Ticket> = {};

// Активность для обработки данных
export async function findTicketBySender(sender: string): Promise<Ticket | null> {
  // Ищем последний открытый тикет для этого отправителя
  const tickets = Object.values(ticketsDB).filter(
    (t) => t.customerId === sender && t.status !== "closed"
  );
  logger.info("[Activity: findTicketBySender] Обработка данных:", tickets);
  return tickets.length > 0 ? tickets[0] : null;
}

export async function createNewTicket(params: TicketParams): Promise<Ticket> {
  const newTicket: Ticket = {
    id: `ticket-${Date.now()}`,
    status: "open",
    customerId: params.customerId,
    channel: params.channel,
    messages: [params.initialMessage],
    previousTicketId: params.previousTicketId,
  };

  ticketsDB[newTicket.id] = newTicket;
  logger.info("[Activity: createNewTicket] Обработка данных:", newTicket);
  return newTicket;
}

export async function addMessageToTicket(ticketId: string, message: Message): Promise<Ticket> {
  const ticket = ticketsDB[ticketId];
  if (!ticket) {
    throw new Error("Ticket not found");
  }

  ticket.messages.push(message);
  logger.info(
    "[Activity: addMessageToTicket] Обработка данных:",
    ticketId,
    message,
    ticket.messages
  );

  return ticket;
}

export async function classifyTicket(ticketId: string): Promise<string> {
  // Здесь может быть вызов ML модели или правила для классификации
  return 'general'; // пример простой классификации
}

export async function assignTicket(ticketId: string, classification: string): Promise<Ticket> {
  const ticket = ticketsDB[ticketId];
  if (!ticket) throw new Error('Ticket not found');
  
  // Простая логика назначения - в реальности может быть сложнее
  const agent = classification === 'priority' ? 'agent-priority' : 'agent-general';
  
  ticket.assignedTo = agent;
  ticket.status = 'in_progress';

  return ticket;
}

export async function closeTicket(ticketId: string): Promise<Ticket> {
  const ticket = ticketsDB[ticketId];
  if (!ticket) throw new Error('Ticket not found');
  
  ticket.status = 'closed';

  return ticket;
}

export async function markAsResolved(ticketId: string, resolution: string): Promise<Ticket> {
  console.log('markAsResolved', ticketId, resolution);
  const ticket = ticketsDB[ticketId];
  if (!ticket) throw new Error('Ticket not found');
  
  ticket.status = 'resolved';

  return ticket;
} 