import express, { Request, Response } from "express";
import bodyParser from "body-parser";
import { createClient } from "./temporal";
import { logger } from "./utils";
import { Client } from "@temporalio/client";

const PORT = process.env.PORT || 3002;

const app = express();
app.use(bodyParser.json());

interface Message {
  id: string;
  [key: string]: any;
}

interface WorkflowRequest {
  workflowType: string;
  message: Message;
}

// Инициализация Temporal клиента
let client: Client;
const setup = async (): Promise<void> => {
  client = await createClient();
};
setup();

// Эндпоинт для запуска воркфлоу
app.post("/message", async (req: Request<{}, {}, WorkflowRequest>, res: Response) => {
  try {
    const { workflowType, message } = req.body;

    const handle = await client.workflow.start(workflowType, {
      workflowId: `process-message-${message.id}`,
      taskQueue: "workflow-ticket",
      args: [message],
    });

    logger.info(`Запущен воркфлоу ${workflowType} с ID: ${handle.workflowId}`);

    res.json({
      workflowId: handle.workflowId,
      status: "started",
    });
  } catch (error) {
    logger.error("Ошибка при запуске воркфлоу:", error);
    res.status(500).json({
      error: "Ошибка при запуске воркфлоу",
      details: error instanceof Error ? error.message : String(error),
    });
  }
});

app.post("/done", async (req: Request<{}, {}, WorkflowRequest>, res: Response) => {
  try {
    const { workflowType, message } = req.body;
    console.log(`process-message-${message.id}`);

    const handle2 = client.workflow.getHandle(`process-message-${message.id}`);

    await handle2.signal('resolve', {
      workflowId: `process-message-${message.id}`,
      args: ['Проблема решена обновлением системы'],
    });

    console.log('handle2', handle2);

    logger.info(`Запущен воркфлоу ${workflowType} с ID: ${handle2.workflowId}`);

    res.json({
      workflowId: handle2.workflowId,
      status: "started",
    });
  } catch (error) {
    logger.error("Ошибка при запуске воркфлоу:", error);
    res.status(500).json({
      error: "Ошибка при запуске воркфлоу",
      details: error instanceof Error ? error.message : String(error),
    });
  }
});

app.listen(PORT, () => {
  logger.info(`Сервер запущен на порту ${PORT}`);
}); 