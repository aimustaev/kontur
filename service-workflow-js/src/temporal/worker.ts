import { Worker, NativeConnection } from "@temporalio/worker";

interface WorkerConfig {
  workflowsPath: string;
  activities: Record<string, Function>;
  taskQueue: string;
}

export const createWorker = async ({
  workflowsPath,
  activities,
  taskQueue,
}: WorkerConfig): Promise<Worker> => {
  const connection = await NativeConnection.connect({
    address: process.env.TEMPORAL_ADDRESS ?? "localhost:7233",
  });

  const worker = await Worker.create({
    workflowsPath,
    activities,
    taskQueue,
    connection,
  });

  return worker;
}; 