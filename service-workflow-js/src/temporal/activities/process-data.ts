import { logger } from "../../utils";

interface ProcessedData {
  processed: boolean;
  processedAt: string;
  [key: string]: any;
}

export const processData = async (data: Record<string, any>): Promise<ProcessedData> => {
  logger.info("[Activity: ProcessData] Обработка данных:", data);

  // Здесь ваша логика обработки данных
  // Например:
  // 1. Валидация данных
  // 2. Преобразование данных
  // 3. Сохранение в базу данных

  return {
    ...data,
    processed: true,
    processedAt: new Date().toISOString(),
  };
}; 