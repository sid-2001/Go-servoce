package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// KafkaClient uses Kafka REST Proxy so we can produce/consume without third-party Go libs.
type KafkaClient struct {
	BaseURL string
	Client  *http.Client
}

func NewKafkaClient(baseURL string) *KafkaClient {
	return &KafkaClient{
		BaseURL: baseURL,
		Client:  &http.Client{Timeout: 5 * time.Second},
	}
}

// Publish sends an event into a Kafka topic via REST Proxy.
func (k *KafkaClient) Publish(topic string, value any) {
	payload := map[string]any{"records": []map[string]any{{"value": value}}}
	raw, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/topics/%s", k.BaseURL, topic), bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/vnd.kafka.json.v2+json")
	req.Header.Set("Accept", "application/vnd.kafka.v2+json")

	res, err := k.Client.Do(req)
	if err != nil {
		log.Printf("kafka publish error: %v", err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		body, _ := io.ReadAll(res.Body)
		log.Printf("kafka publish status=%d body=%s", res.StatusCode, string(body))
	}
}

// StartConsumer subscribes to a topic and streams records to handler.
func (k *KafkaClient) StartConsumer(groupID, consumerName, topic string, handler func([]byte)) {
	go func() {
		base := fmt.Sprintf("%s/consumers/%s", k.BaseURL, groupID)
		createBody := map[string]any{
			"name":              consumerName,
			"format":            "json",
			"auto.offset.reset": "earliest",
		}
		raw, _ := json.Marshal(createBody)
		req, _ := http.NewRequest(http.MethodPost, base, bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/vnd.kafka.v2+json")
		res, err := k.Client.Do(req)
		if err != nil {
			log.Printf("consumer create error: %v", err)
			return
		}
		defer res.Body.Close()

		subBody := map[string]any{"topics": []string{topic}}
		raw, _ = json.Marshal(subBody)
		req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/instances/%s/subscription", base, consumerName), bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/vnd.kafka.v2+json")
		if _, err := k.Client.Do(req); err != nil {
			log.Printf("consumer subscribe error: %v", err)
			return
		}

		for {
			req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/instances/%s/records?timeout=1000", base, consumerName), nil)
			req.Header.Set("Accept", "application/vnd.kafka.json.v2+json")
			res, err := k.Client.Do(req)
			if err != nil {
				log.Printf("consumer read error: %v", err)
				time.Sleep(2 * time.Second)
				continue
			}
			var records []struct {
				Value json.RawMessage `json:"value"`
			}
			if err := json.NewDecoder(res.Body).Decode(&records); err != nil {
				res.Body.Close()
				time.Sleep(1 * time.Second)
				continue
			}
			res.Body.Close()
			for _, record := range records {
				handler(record.Value)
			}
		}
	}()
}
