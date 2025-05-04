const { Worker, NativeConnection } = require("@temporalio/worker");
const winston = require("winston");
const { processData } = require("./activities");
console.log("TEMPORAL_ADDRESS:", process.env.TEMPORAL_ADDRESS);
// Настройка логгера
const logger = winston.createLogger({
  level: "info",
  format: winston.format.json(),
  transports: [new winston.transports.Console()],
});

// Определение воркфлоу
async function exampleWorkflow(input) {
  logger.info("Начало выполнения воркфлоу с входными данными:", input);

  // Здесь ваша логика воркфлоу
  // Например:
  // 1. Получение данных
  // 2. Обработка
  // 3. Сохранение результатов

  logger.info("Воркфлоу успешно завершен");
  return { status: "completed", result: "Результат обработки" };
}

async function run() {
  try {
    const connection = await NativeConnection.connect({
      address: process.env.TEMPORAL_ADDRESS ?? "localhost:3005",
    });
    const worker = await Worker.create({
      workflowsPath: require.resolve("./workflows"),
      activities: {
        processData,
      },
      taskQueue: "workflow-queue",
      connection,
    });

    logger.info("Worker запущен");
    await worker.run();
  } catch (error) {
    logger.error("Ошибка при запуске worker:", error);
    process.exit(1);
  }
}

run().catch((err) => {
  logger.error("Критическая ошибка:", err);
  process.exit(1);
});
