const { proxyActivities } = require('@temporalio/workflow');

// Определение воркфлоу
async function exampleWorkflow(input) {
  const { processData } = proxyActivities({ startToCloseTimeout: '1 minute' });

  // Шаг 1: Получение данных
  const data = await processData(input);

  // Шаг 2: Обработка данных
  const result = {
    status: 'completed',
    processedData: data,
    timestamp: new Date().toISOString()
  };

  return result;
}

module.exports = {
  exampleWorkflow
}; 