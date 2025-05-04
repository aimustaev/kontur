const express = require("express");
const bodyParser = require("body-parser");
const { Client, Connection } = require("@temporalio/client");
const winston = require("winston");
const { NativeConnection } = require("@temporalio/worker");

// Настройка логгера
const logger = winston.createLogger({
  level: "info",
  format: winston.format.json(),
  transports: [new winston.transports.Console()],
});

const app = express();
app.use(bodyParser.json());
logger.info(process.env.TEMPORAL_ADDRESS);
// Инициализация Temporal клиента
let client;
const run = async () => {
  const connection = await Connection.connect({
    address: process.env.TEMPORAL_ADDRESS ?? "localhost:3005",
  });
  console.log('1231', connection)
  client = new Client({ connection });
  console.log('1231', client)
};

run();

// Эндпоинт для запуска воркфлоу
app.post("/start", async (req, res) => {
  try {
    logger.info(client, '1232131');
    console.log('1231', client)
    const { workflowId, workflowType, input } = req.body;

    if (!workflowId || !workflowType) {
      return res.status(400).json({
        error: "Необходимо указать workflowId и workflowType",
      });
    }

    const handle = await client.workflow.start(workflowType, {
      workflowId,
      taskQueue: "workflow-queue",
      args: ["Temporal"],
    });

    logger.info(`Запущен воркфлоу ${workflowType} с ID: ${workflowId}`);

    res.json({
      workflowId: handle.workflowId,
      status: "started",
    });
  } catch (error) {
    logger.error("Ошибка при запуске воркфлоу:", error);
    res.status(500).json({
      error: "Ошибка при запуске воркфлоу",
      details: error.message,
    });
  }
});

const PORT = process.env.PORT || 3002;
app.listen(PORT, () => {
  logger.info(`Сервер запущен на порту ${PORT}`);
});
