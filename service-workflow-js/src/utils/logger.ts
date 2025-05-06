import winston from 'winston';

// Настройка логгера
export const logger: winston.Logger = winston.createLogger({
  level: 'info',
  format: winston.format.json(),
  transports: [
    new winston.transports.Console()
  ]
}); 