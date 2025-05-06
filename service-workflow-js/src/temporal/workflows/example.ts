import { proxyActivities } from "@temporalio/workflow";

interface ProcessedData {
  processed: boolean;
  processedAt: string;
  [key: string]: any;
}

interface WorkflowResult {
  status: string;
  processedData: ProcessedData;
  timestamp: string;
}

// Определение воркфлоу
export const exampleWorkflow = async (input: Record<string, any>): Promise<WorkflowResult> => {
  const { processData } = proxyActivities({ startToCloseTimeout: "1 minute" });

  // Шаг 1: Получение данных
  const data = await processData(input);

  // Шаг 2: Обработка данных
  const result: WorkflowResult = {
    status: "completed",
    processedData: data,
    timestamp: new Date().toISOString(),
  };

  return result;
}; 