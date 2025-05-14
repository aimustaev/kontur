package engine

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"

	"github.com/aimustaev/service-workflow/internal/generated/proto"
)

type WorkflowDefinition struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	States      []StateDefinition `json:"states"`
	Timeouts    Timeouts          `json:"timeouts"`
	InputSchema json.RawMessage   `json:"inputSchema"`
}

type StateDefinition struct {
	Name          string            `json:"name"`
	Type          string            `json:"type"` // activity, signal, timer, etc.
	ActivityName  string            `json:"activityName,omitempty"`
	Input         json.RawMessage   `json:"input,omitempty"`
	Output        string            `json:"output,omitempty"`
	OutputSchema  json.RawMessage   `json:"outputSchema,omitempty"`
	SignalName    string            `json:"signalName,omitempty"`
	Actions       []StateDefinition `json:"actions,omitempty"`
	Concurrent    bool              `json:"concurrent,omitempty"`
	Timeouts      Timeouts          `json:"timeouts,omitempty"`
	TimerDuration string            `json:"timerDuration,omitempty"` // Например: "10s", "1m30s"
}

type Timeouts struct {
	StartToClose    string `json:"startToClose"`
	ScheduleToClose string `json:"scheduleToClose"`
}

type WorkflowEngine struct {
	temporalClient client.Client
	activities     map[string]interface{}
}

func NewEngine(temporalClient client.Client, activities map[string]interface{}) *WorkflowEngine {
	return &WorkflowEngine{
		temporalClient: temporalClient,
		activities:     activities,
	}
}

func (e *WorkflowEngine) ExecuteWorkflow(ctx workflow.Context, def WorkflowDefinition, input interface{}) (interface{}, error) {
	state := make(map[string]interface{})
	state["input"] = input

	// Устанавливаем таймауты по умолчанию
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout:    time.Minute * 5,  // 5 минут на выполнение активности
		ScheduleToCloseTimeout: time.Minute * 10, // 10 минут от планирования до завершения
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	for _, stateDef := range def.States {
		switch stateDef.Type {
		case "activity":
			err := e.executeActivity(ctx, stateDef, state)
			if err != nil {
				return nil, fmt.Errorf("state %s failed: %w", stateDef.Name, err)
			}

		case "signal":
			e.executeSignalHandler(ctx, stateDef, state)

		case "timer":
			err := e.executeTimer(ctx, stateDef, state)
			if err != nil {
				return nil, fmt.Errorf("timer %s failed: %w", stateDef.Name, err)
			}

		default:
			return nil, fmt.Errorf("unknown state type: %s", stateDef.Type)
		}
	}

	return state["output"], nil
}

func (e *WorkflowEngine) executeActivity(ctx workflow.Context, def StateDefinition, state map[string]interface{}) error {
	logger := workflow.GetLogger(ctx)
	// Парсим входные данные
	input, err := parseInput(def.Input, state)
	if err != nil {
		return fmt.Errorf("failed to parse input: %w", err)
	}

	activity := e.activities[def.ActivityName]

	// Определяем тип выходных данных на основе имени activity
	var output interface{}
	switch def.ActivityName {
	case "GetOrCreateTicketActivity",
		"AddMassageToTicketActivity",
		"ClassifierAcitivity",
		"SolveTicketAcitivity",
		"GetTicketByUserActivity",
		"CreateTicketActivity":
		var ticketOutput *proto.TicketResponse
		err = workflow.ExecuteActivity(ctx, activity, input...).Get(ctx, &ticketOutput)
		output = ticketOutput
	case "ProcessMessageActivity",
		"WaitActivity":
		var errOutput error
		err = workflow.ExecuteActivity(ctx, activity, input...).Get(ctx, &errOutput)
		output = errOutput
	default:
		err = workflow.ExecuteActivity(ctx, activity, input...).Get(ctx, &output)
	}

	if err != nil {
		return err
	}

	//// Если есть схема выходных данных, валидируем результат
	//if len(def.OutputSchema) > 0 {
	//	// TODO: Добавить валидацию по схеме
	//	// Пока просто логируем, что схема есть
	//	logger.Info("Activity output schema defined", "schema", string(def.OutputSchema))
	//}

	// Сохраняем результат если нужно
	if def.Output != "" {
		state[def.Output] = output
	}

	logger.Info("Activity completed", "name", def.Name, "output", output)
	return nil
}

func (e *WorkflowEngine) executeSignalHandler(ctx workflow.Context, def StateDefinition, state map[string]interface{}) {
	workflow.Go(ctx, func(ctx workflow.Context) {
		signalChan := workflow.GetSignalChannel(ctx, def.SignalName)
		for {
			var signalData interface{}
			signalChan.Receive(ctx, &signalData)
			state["signalPayload"] = signalData

			for _, action := range def.Actions {
				if action.Type == "activity" {
					err := e.executeActivity(ctx, action, state)
					if err != nil {
						workflow.GetLogger(ctx).Error("Signal handler activity failed", "error", err)
					}
				}
			}
		}
	})
}

func (e *WorkflowEngine) executeTimer(ctx workflow.Context, def StateDefinition, state map[string]interface{}) error {
	logger := workflow.GetLogger(ctx)

	if def.TimerDuration == "" {
		return fmt.Errorf("timer duration not specified for state %s", def.Name)
	}

	duration, err := time.ParseDuration(def.TimerDuration)
	if err != nil {
		return fmt.Errorf("invalid timer duration for state %s: %w", def.Name, err)
	}

	logger.Info("Starting timer", "name", def.Name, "duration", duration)

	// Создаем таймер
	timerCtx, cancel := workflow.WithCancel(ctx)
	timer := workflow.NewTimer(timerCtx, duration)

	// Ожидаем срабатывания таймера или отмены
	selector := workflow.NewSelector(ctx)
	selector.AddFuture(timer, func(f workflow.Future) {
		// Таймер сработал
		logger.Info("Timer completed", "name", def.Name)

		// Выполняем действия после таймера, если они есть
		for _, action := range def.Actions {
			switch action.Type {
			case "activity":
				if err := e.executeActivity(ctx, action, state); err != nil {
					logger.Error("Post-timer activity failed", "error", err)
				}
			case "signal":
				e.executeSignalHandler(ctx, action, state)
			}
		}
	})

	// Добавляем возможность отмены через сигнал
	if def.SignalName != "" {
		cancelChan := workflow.GetSignalChannel(ctx, def.SignalName+"_cancel")
		selector.AddReceive(cancelChan, func(c workflow.ReceiveChannel, more bool) {
			cancel() // Отменяем таймер
			logger.Info("Timer cancelled by signal", "name", def.Name)
		})
	}

	selector.Select(ctx) // Блокируемся до срабатывания таймера или отмены

	return nil
}

func parseInput(inputJSON json.RawMessage, state map[string]interface{}) ([]interface{}, error) {
	var inputs []interface{}

	// Если input - массив
	if len(inputJSON) > 0 && inputJSON[0] == '[' {
		var inputArray []json.RawMessage
		if err := json.Unmarshal(inputJSON, &inputArray); err != nil {
			return nil, err
		}

		for _, item := range inputArray {
			val, err := parseValue(item, state)
			if err != nil {
				return nil, err
			}
			inputs = append(inputs, val)
		}
	} else {
		// Если input - одиночное значение
		val, err := parseValue(inputJSON, state)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, val)
	}

	return inputs, nil
}

func parseValue(value json.RawMessage, state map[string]interface{}) (interface{}, error) {
	// Проверяем, является ли значение ссылкой на state (начинается с "$.")
	if len(value) > 2 && value[0] == '"' && value[1] == '$' {
		var ref string
		if err := json.Unmarshal(value, &ref); err != nil {
			return nil, err
		}

		if ref[:2] == "$." {
			key := ref[2:]
			if val, ok := getNestedValue(state, key); ok {
				return val, nil
			}
			return nil, fmt.Errorf("state key not found: %s", key)
		}
	}

	// Если это не ссылка, разбираем как обычное JSON значение
	var result interface{}
	if err := json.Unmarshal(value, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func parseDuration(durStr string) time.Duration {
	if durStr == "" {
		return 0
	}
	dur, err := time.ParseDuration(durStr)
	if err != nil {
		return 0
	}
	return dur
}

// getNestedValue позволяет получить значение по ключу с точками (например, "input.Message") из вложенных map[string]interface{}.
func getNestedValue(state map[string]interface{}, key string) (interface{}, bool) {
	keys := strings.Split(key, ".")
	current := interface{}(state)

	for _, k := range keys {
		// Пытаемся обработать как map[string]interface{}
		if m, ok := current.(map[string]interface{}); ok {
			val, exists := m[k]
			if !exists {
				return nil, false
			}
			current = val
			continue
		}

		// Пытаемся обработать как protobuf-сообщение (через рефлексию)
		val := reflect.ValueOf(current)
		if val.Kind() == reflect.Ptr {
			val = val.Elem() // Разыменовываем указатель (*proto.Ticket → proto.Ticket)
		}

		if val.Kind() != reflect.Struct {
			return nil, false // Не мапа и не структура — ошибка
		}

		// Ищем поле в protobuf-структуре
		field := val.FieldByName(k)
		if !field.IsValid() {
			return nil, false
		}

		current = field.Interface()
	}

	return current, true
}
