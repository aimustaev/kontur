import { Kafka, Consumer, EachMessagePayload } from "kafkajs";
import { createClient } from "./temporal";
import { logger } from "./utils";
import { WorkflowClient } from "@temporalio/client";

interface KafkaEvent {
  data: {
    workflowId: string;
    workflowType: string;
    input?: any;
  };
}

// Инициализация Temporal клиента
let client: WorkflowClient;
const setup = async (): Promise<void> => {
  client = await createClient();
};
setup();

// Инициализация Kafka consumer
const kafka = new Kafka({
  clientId: "workflow-consumer",
  brokers: [process.env.KAFKA_BROKER ?? "localhost:9092"],
});

const consumer: Consumer = kafka.consumer({ groupId: "workflow-group" });

const run = async (): Promise<void> => {
  try {
    // Подключаемся к Kafka
    await consumer.connect();
    logger.info("Connected to Kafka");

    // Подписываемся на топик
    await consumer.subscribe({
      topic: process.env.KAFKA_TOPIC ?? "workflow-events",
      fromBeginning: false,
    });
    logger.info("Subscribed to Kafka topic");

    // Запускаем обработку сообщений
    await consumer.run({
      eachMessage: async ({ topic, partition, message }: EachMessagePayload) => {
        try {
          const event: KafkaEvent = JSON.parse(message.value?.toString() ?? "{}");
          logger.info(`Received event from Kafka: ${JSON.stringify(event)}`);

          const {
            data: { workflowId, workflowType, input },
          } = event;
          console.log("1231", workflowId, workflowType, input, event);
          if (!workflowId || !workflowType) {
            logger.error("workflowId and workflowType are required");
            return;
          }

          const handle = await client.workflow.start(workflowType, {
            workflowId,
            taskQueue: "workflow-queue",
            args: input ? [input] : [],
          });

          logger.info(
            `Started workflow ${workflowType} with ID: ${workflowId}`
          );
        } catch (error) {
          logger.error("Error processing Kafka message:", error);
        }
      },
    });
  } catch (error) {
    logger.error("Fatal error:", error);
    process.exit(1);
  }
};

// Обработка завершения процесса
process.on("SIGTERM", async () => {
  logger.info("SIGTERM received. Shutting down...");
  await consumer.disconnect();
  process.exit(0);
});

process.on("SIGINT", async () => {
  logger.info("SIGINT received. Shutting down...");
  await consumer.disconnect();
  process.exit(0);
});

run().catch((e: Error) => {
  logger.error(e);
  process.exit(1);
}); 