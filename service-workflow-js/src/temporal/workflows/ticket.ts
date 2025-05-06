import { proxyActivities } from '@temporalio/workflow';
import { defineSignal, setHandler, sleep, condition } from '@temporalio/workflow';

export const resolveSignal = defineSignal<[string]>('resolve');
export const reassignSignal = defineSignal<[string]>('reassign');

interface Message {
  id: string;
  sender: string;
  channel: string;
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

const { 
  findTicketBySender,
  createNewTicket,
  addMessageToTicket,
  classifyTicket,
  assignTicket,
  closeTicket,
  markAsResolved
} = proxyActivities({
  startToCloseTimeout: '1 minute',
});

// Основной workflow обработки сообщения
export async function processMessageWorkflow(message: Message): Promise<void> {
  // 1. Пытаемся найти открытый тикет по отправителю
  const existingTicket = await findTicketBySender(message.sender);
  
  if (existingTicket) {
    if (existingTicket.status !== 'closed') {
      // 2. Если тикет есть и не закрыт - добавляем сообщение
      await addMessageToTicket(existingTicket.id, message);
    } else {
      // 3. Если тикет закрыт - создаем новый с ссылкой на старый
      const newTicket = await createNewTicket({
        customerId: message.sender,
        channel: message.channel,
        previousTicketId: existingTicket.id,
        initialMessage: message
      });
      
      // Запускаем процесс классификации
      await classifyAndAssignTicket(newTicket.id);
    }
  } else {
    // 4. Если тикета нет - создаем новый
    const newTicket = await createNewTicket({
      customerId: message.sender,
      channel: message.channel,
      initialMessage: message
    });
    
    // Запускаем процесс классификации
    await classifyAndAssignTicket(newTicket.id);
  }
}

// Вложенный workflow для классификации и назначения тикета
async function classifyAndAssignTicket(ticketId: string): Promise<void> {
  let status: string = 'open';
  let resolutionResult: string | undefined;

  try {
    // Классифицируем тикет
    const classification = await classifyTicket(ticketId);
    status = 'classified';

    // Назначаем тикет агенту
    await assignTicket(ticketId, classification);
    status = 'assigned';

    // Здесь можно добавить логику ожидания решения агента
    // Например, подписаться на событие "тикет решен"
    setHandler(resolveSignal, (resolution) => {
      resolutionResult = resolution;
    });

    setHandler(reassignSignal, async (newAgentId) => {
      await assignTicket(ticketId, classification);
    });

    await condition(() => resolutionResult !== undefined);

    if (resolutionResult) {
      await markAsResolved(ticketId, resolutionResult);
      status = 'resolved';

      await closeTicket(ticketId);
      status = 'closed';
    }
  } catch (error) {
    // Обработка ошибок классификации/назначения
    console.error('Error in classifyAndAssignTicket:', error);
  }
}

// Workflow для обработки закрытия тикета агентом
export async function resolveTicketWorkflow(ticketId: string): Promise<void> {
  await closeTicket(ticketId);
  
  // Здесь можно добавить дополнительные действия после закрытия,
  // например, отправку уведомления пользователю
} 