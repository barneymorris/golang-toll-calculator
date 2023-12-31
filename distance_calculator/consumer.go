package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/betelgeusexru/golang-toll-calculator/aggregator/client"
	"github.com/betelgeusexru/golang-toll-calculator/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	Consumer *kafka.Consumer
	IsRunning bool
	CalcService CalculatorServicer
	AggClient client.Client
}

func NewKafkaConsumer(topic string, svc CalculatorServicer, aggClient client.Client) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{topic}, nil)

	return &KafkaConsumer{
		Consumer: c,
		CalcService: svc,
		AggClient: aggClient,
	}, nil
}

func (c *KafkaConsumer) Start() {
	logrus.Info("kafka transport started")
	c.IsRunning = true
	c.ReadMessageLoop()
}

func (c *KafkaConsumer) ReadMessageLoop() {
	for c.IsRunning {
		msg, err := c.Consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("kafka consume error %s", err)
			continue
		}

		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("json serialization error: %s", err)
			continue
		}

		distance, err := c.CalcService.CalculateDistance(data);
		if err != nil {
			logrus.Errorf("calculation error: %s", err)
			continue
		}

		req := &types.AggregateRequest{
			Value: distance,
			Unix: time.Now().UnixNano(),
			ObuID: int32(data.OBUID),
		}
		if err := c.AggClient.Aggregate(context.Background(), req); err != nil {
			logrus.Errorf("aggregate error: ", err)
			continue
		}
	}
}