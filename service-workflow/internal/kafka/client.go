package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"

	"github.com/aimustaev/service-workflow/internal/config"
)

// Client represents a Kafka client
type Client struct {
	consumer sarama.ConsumerGroup
	config   *sarama.Config
	brokers  []string
	topic    string
	groupID  string
}

// NewClient creates a new Kafka client
func NewClient(cfg *config.Config) (*Client, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Version = sarama.V2_8_0_0

	// Verify brokers availability
	admin, err := sarama.NewClusterAdmin(cfg.Kafka.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster admin: %w", err)
	}
	defer admin.Close()

	// Get topics list
	topics, err := admin.ListTopics()
	if err != nil {
		return nil, fmt.Errorf("failed to list topics: %w", err)
	}
	log.Printf("Available topics: %v", topics)

	// Create consumer group
	consumer, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, cfg.Kafka.GroupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &Client{
		consumer: consumer,
		config:   config,
		brokers:  cfg.Kafka.Brokers,
		topic:    cfg.Kafka.Topic,
		groupID:  cfg.Kafka.GroupID,
	}, nil
}

// Consume starts consuming messages from all partitions
func (c *Client) Consume(ctx context.Context, messageHandler func(msg *sarama.ConsumerMessage)) error {
	topics := []string{c.topic}
	consumer := &consumerGroupHandler{
		messageHandler: messageHandler,
	}

	for {
		err := c.consumer.Consume(ctx, topics, consumer)
		if err != nil {
			return fmt.Errorf("error from consumer: %w", err)
		}
	}
}

// Close closes the Kafka consumer
func (c *Client) Close() error {
	return c.consumer.Close()
}

// consumerGroupHandler represents a Sarama consumer group consumer
type consumerGroupHandler struct {
	messageHandler func(msg *sarama.ConsumerMessage)
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages()
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		h.messageHandler(message)
		// Mark message as processed
		session.MarkMessage(message, "")
	}
	return nil
}
