const ACTIVITY_NODES = [
  {
    type: "activity",
    label: "GetOrCreateTicket",
    description: "Создать или получить тикет",
    activityName: "GetOrCreateTicketActivity",
    argIn: ["message"],
    argOut: ["ticket"],
  },
  {
    type: "activity",
    label: "AddInitialMessage",
    description: "Добавить первое сообщение",
    activityName: "AddMassageToTicketActivity",
    argIn: ["message", "ticketId"],
    argOut: [],
  },
  {
    type: "activity",
    label: "MessageListener",
    description: "Слушатель сообщений",
    activityName: "MessageListener",
  },
  {
    type: "activity",
    label: "ClassifyTicket",
    description: "Классификация тикета",
    activityName: "ClassifierAcitivity",
    argIn: ["ticket"],
    argOut: ["ticket"],
  },
  {
    type: "activity",
    label: "SolveTicket",
    description: "Решить тикет",
    activityName: "SolveTicketAcitivity",
    argIn: ["ticket"],
    argOut: ["ticket"],
  },
];
const TIMER_NODES = [
  {
    type: "timer",
    label: "Timer",
    description: "Таймер ожидания",
    argIn: ["timerDuration"],
    argOut: [],
  },
];
const SIGNAL_NODES = [
  {
    type: "signal",
    label: "Signal",
    description: "Сигнал внешнего события",
    signalName: "NewMessage",
    argIn: [],
    argOut: [],
  },
];

export { ACTIVITY_NODES, TIMER_NODES, SIGNAL_NODES };
