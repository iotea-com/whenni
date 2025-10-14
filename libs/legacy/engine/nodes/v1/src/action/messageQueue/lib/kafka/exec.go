package kafkaActionNode

import (
	"time"

	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

func (n *KafkaActionSubnode) Exec(params node.ExecParams) node.Error {
Loop:
	for {
		select {
		case <-params.Ctx.Done():
			break Loop
		case ioData := <-n.inputChannels[0].Channel:
			err, invalid := n.handlePayload(ioData)
			if err != nil {
				if invalid {
					// Safely notify the runtime that data was processed
					n.notifyChannel <- node.Notification{
						Type:    node.NotifyDataInvalid,
						DataCtx: ioData.Ctx,
						Reason:  err.Error(),
					}
				} else {
					return node.Error{
						Type:   node.FatalError,
						Reason: err.Error(),
					}
				}
			} else {
				// Safely notify the runtime that data was processed
				n.notifyChannel <- node.Notification{
					Type:    node.NotifyDataProcessed,
					DataCtx: ioData.Ctx,
				}
			}
		}
	}

	return node.Error{
		Type: node.NoError,
	}
}

func (n *KafkaActionSubnode) handlePayload(ioData node.IoData) (error, bool) {
	// Verify that there is a request body
	if ioData.Data == nil {
		return fmt.Errorf("no data was passed in input channel"), true // Validation error
	}

	// Validate that the data is a []byte
	payloadData, ok := ioData.Data.([]byte)
	if !ok {
		return fmt.Errorf("expected data of type []byte, got %T", ioData.Data), true // Validation error
	}

	// Prepare the message
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic: &n.config.Topic,
		},
		Value:         payloadData,
		Timestamp:     time.Now(),
		TimestampType: kafka.TimestampCreateTime,
	}

	n.logger.Info().Ctx(ioData.Ctx).
		Str("data", string(payloadData)).
		Msgf("Producing message to Kafka topic %s", *message.TopicPartition.Topic)

	// Create a delivery channel to receive the delivery report
	deliveryChan := make(chan kafka.Event, 1)

	// Produce the message asynchronously
	err := n.producer.Produce(message, deliveryChan)
	if err != nil {
		return fmt.Errorf("error enqueuing message for Kafka cluster: %s", err), false // Fatal error
	}

	// Block until a delivery report is received or timeout occurs
	select {
	case e := <-deliveryChan:
		m := e.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			return fmt.Errorf("error delivering message to Kafka topic %s: %v", *m.TopicPartition.Topic, m.TopicPartition.Error), false // Fatal error
		}
		n.logger.Info().Ctx(ioData.Ctx).
			Msgf("Message delivered to Kafka topic %s, partition %d, offset %v", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	case <-time.After(3 * time.Second):
		return fmt.Errorf("timeout waiting for message delivery"), false // Fatal error
	}

	return nil, false
}
