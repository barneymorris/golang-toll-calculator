package main

import (
	"encoding/json"
	"fmt"

	"github.com/betelgeusexru/golang-toll-calculator/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	Consumer *kafka.Consumer
	IsRunning bool
	CalcService CalculatorServicer
}

func NewKafkaConsumer(topic string, svc CalculatorServicer) (*KafkaConsumer, error) {
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

		fmt.Println(distance)
	}
}