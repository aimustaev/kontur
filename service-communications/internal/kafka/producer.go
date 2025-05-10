package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"

	"github.com/aimustaev/service-communications/internal/adapter"
)

// Producer represents a Kafka message producer
type Producer struct {
	producer sarama.SyncProducer
	topic    string
	logger   *logrus.Logger
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string, topic string, logger *logrus.Logger) (*Producer, error) {
	logger.Infof("Initializing Kafka producer with brokers: %v, topic: %s", brokers, topic)

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	// Enable debug logging
	config.Net.MaxOpenRequests = 1
	config.Net.DialTimeout = 10 * time.Second
	config.Net.ReadTimeout = 10 * time.Second
	config.Net.WriteTimeout = 10 * time.Second

	// Print full broker configuration
	logger.Infof("Kafka config - MaxOpenRequests: %d", config.Net.MaxOpenRequests)
	logger.Infof("Kafka config - DialTimeout: %v", config.Net.DialTimeout)
	logger.Infof("Kafka config - ReadTimeout: %v", config.Net.ReadTimeout)
	logger.Infof("Kafka config - WriteTimeout: %v", config.Net.WriteTimeout)

	// Try to connect to each broker individually
	for _, broker := range brokers {
		logger.Infof("Attempting to connect to broker: %s", broker)
	}

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		logger.Errorf("Failed to create producer with brokers %v: %v", brokers, err)
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	logger.Infof("Successfully connected to Kafka brokers: %v", brokers)
	return &Producer{
		producer: producer,
		topic:    topic,
		logger:   logger,
	}, nil
}

// Close closes the Kafka producer
func (p *Producer) Close() error {
	return p.producer.Close()
}

// SendMessage sends a message to Kafka
func (p *Producer) SendMessage(ctx context.Context, msg adapter.Message) error {
	// Convert message to JSON
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Create Kafka message
	kafkaMsg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(jsonData),
	}

	// Send message
	partition, offset, err := p.producer.SendMessage(kafkaMsg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	p.logger.Infof("Message sent to Kafka: topic=%s partition=%d offset=%d", p.topic, partition, offset)
	return nil
}
