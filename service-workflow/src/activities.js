const winston = require('winston');

// Настройка логгера
const logger = winston.createLogger({
  level: 'info',
  format: winston.format.json(),
  transports: [
    new winston.transports.Console()
  ]
});

// Активность для обработки данных
async function processData(input) {
  logger.info('Обработка данных:', input);
  
  // Здесь ваша логика обработки данных
  // Например:
  // 1. Валидация данных
  // 2. Преобразование данных
  // 3. Сохранение в базу данных
  
  return {
    ...input,
    processed: true,
    processedAt: new Date().toISOString()
  };
}

module.exports = {
  processData
}; 