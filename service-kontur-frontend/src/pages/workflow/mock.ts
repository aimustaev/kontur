export const mockWorkflow = {
    "nodes": [
        {
            "id": "step-1747296656224",
            "type": "activity",
            "data": {
                "label": "GetOrCreateTicket",
                "activityName": "GetOrCreateTicketActivity",
                "timerDuration": "",
                "signalName": "",
                "argIn": [
                    "message"
                ],
                "argOut": [
                    "ticket"
                ],
                "description": "Создать или получить тикет",
                "input": {
                    "message": "$.input.Message"
                },
                "output": {
                    "ticket": "ticket"
                }
            },
            "position": {
                "x": 23.720213298347588,
                "y": -8.46597968371043
            }
        },
        {
            "id": "step-1747296657141",
            "type": "activity",
            "data": {
                "label": "AddInitialMessage",
                "activityName": "AddMassageToTicketActivity",
                "timerDuration": "",
                "signalName": "",
                "argIn": [
                    "message",
                    "ticketId"
                ],
                "argOut": [
                    "ticket"
                ],
                "description": "Добавить первое сообщение",
                "input": {
                    "message": "$.input.Message",
                    "ticketId": "$.ticket.Id"
                },
                "output": {
                    "ticket": "ticket"
                }
            },
            "position": {
                "x": 418.4195143096235,
                "y": -45.76866758669604
            }
        },
        {
            "id": "step-1747297690439",
            "type": "signal",
            "data": {
                "label": "Signal",
                "activityName": "",
                "timerDuration": "",
                "signalName": "NewMessage",
                "argIn": [],
                "argOut": [],
                "description": "Сигнал внешнего события",
                "input": {},
                "output": {}
            },
            "position": {
                "x": 764.3720544784059,
                "y": -16.304013845826624
            }
        },
        {
            "id": "step-1747297694553",
            "type": "activity",
            "data": {
                "label": "ClassifyTicket",
                "activityName": "ClassifierAcitivity",
                "timerDuration": "",
                "signalName": "",
                "argIn": [
                    "ticket"
                ],
                "argOut": [
                    "ticket"
                ],
                "description": "Классификация тикета",
                "input": {
                    "ticket": "$.ticket"
                },
                "output": {
                    "ticket": "ticket"
                }
            },
            "position": {
                "x": 1015.8841680504928,
                "y": -55.10207392299451
            }
        },
        {
            "id": "step-1747297695629",
            "type": "timer",
            "data": {
                "label": "Timer",
                "activityName": "",
                "timerDuration": "",
                "signalName": "",
                "argIn": [
                    "timerDuration"
                ],
                "argOut": [],
                "description": "Таймер ожидания",
                "input": {
                    "timerDuration": "3s"
                },
                "output": {}
            },
            "position": {
                "x": 1359.221801753839,
                "y": -30.114927913195093
            }
        },
        {
            "id": "step-1747297696431",
            "type": "activity",
            "data": {
                "label": "SolveTicket",
                "activityName": "SolveTicketAcitivity",
                "timerDuration": "",
                "signalName": "",
                "argIn": [
                    "ticket"
                ],
                "argOut": [
                    "ticket"
                ],
                "description": "Решить тикет",
                "input": {
                    "ticket": "$.ticket"
                },
                "output": {
                    "ticket": "ticket"
                }
            },
            "position": {
                "x": 1784.6362086173442,
                "y": -84.65512884098544
            }
        }
    ],
    "edges": [
        {
            "id": "reactflow__edge-step-1747296656224-step-1747296657141",
            "source": "step-1747296656224",
            "target": "step-1747296657141"
        },
        {
            "id": "reactflow__edge-step-1747296657141-step-1747297690439",
            "source": "step-1747296657141",
            "target": "step-1747297690439"
        },
        {
            "id": "reactflow__edge-step-1747297690439-step-1747297694553",
            "source": "step-1747297690439",
            "target": "step-1747297694553"
        },
        {
            "id": "reactflow__edge-step-1747297694553-step-1747297695629",
            "source": "step-1747297694553",
            "target": "step-1747297695629"
        },
        {
            "id": "reactflow__edge-step-1747297695629-step-1747297696431",
            "source": "step-1747297695629",
            "target": "step-1747297696431"
        }
    ]
}