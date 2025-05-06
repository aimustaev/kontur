import { Worker } from '@temporalio/worker';
import { createClient } from './temporal';
import { logger } from './utils';
import { ticketActivities, createWorker } from "./temporal";

async function run(): Promise<void> {
  try {
    const worker = await createWorker({
      workflowsPath: new URL("./temporal/workflows/ticket.ts", import.meta.url)
        .pathname,
      activities: ticketActivities,
      taskQueue: "workflow-ticket",
    });

    logger.info("Worker запущен");
    await worker.run();
  } catch (error) {
    logger.error("Ошибка при запуске worker:", error);
    process.exit(1);
  }
}

run().catch((err: Error) => {
  logger.error("Критическая ошибка:", err);
  process.exit(1);
}); 