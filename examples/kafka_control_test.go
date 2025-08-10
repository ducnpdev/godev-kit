package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

const baseURL = "http://localhost:10000/api/v1/kafka"

// KafkaMessage represents a Kafka message
type KafkaMessage struct {
	Topic string      `json:"topic"`
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// KafkaStatus represents Kafka status response
type KafkaStatus struct {
	Status string                 `json:"status"`
	Data   map[string]interface{} `json:"data"`
}

// TestKafkaControl demonstrates the Kafka control functionality
func TestKafkaControl(t *testing.T) {
	fmt.Println("🚀 Kafka Control Demo")
	fmt.Println("=====================")

	// 1. Check initial status
	fmt.Println("\n1. Checking initial Kafka status...")
	status := getKafkaStatus()
	printStatus(status)

	// 2. Send a test message (should work)
	fmt.Println("\n2. Sending test message (should work)...")
	message := KafkaMessage{
		Topic: "test-topic",
		Key:   "demo-key",
		Value: "Hello Kafka!",
	}
	err := sendMessage(message)
	if err != nil {
		fmt.Printf("❌ Error sending message: %v\n", err)
	} else {
		fmt.Println("✅ Message sent successfully")
	}

	// 3. Disable producer
	fmt.Println("\n3. Disabling Kafka producer...")
	disableProducer()

	// 4. Try to send message (should fail)
	fmt.Println("\n4. Trying to send message (should fail)...")
	err = sendMessage(message)
	if err != nil {
		fmt.Printf("❌ Expected error: %v\n", err)
	} else {
		fmt.Println("⚠️  Unexpected: Message sent successfully")
	}

	// 5. Check status after disabling producer
	fmt.Println("\n5. Checking status after disabling producer...")
	status = getKafkaStatus()
	printStatus(status)

	// 6. Enable producer
	fmt.Println("\n6. Enabling Kafka producer...")
	enableProducer()

	// 7. Send message again (should work)
	fmt.Println("\n7. Sending message again (should work)...")
	err = sendMessage(message)
	if err != nil {
		fmt.Printf("❌ Error sending message: %v\n", err)
	} else {
		fmt.Println("✅ Message sent successfully")
	}

	// 8. Disable consumer
	fmt.Println("\n8. Disabling Kafka consumer...")
	disableConsumer()

	// 9. Check final status
	fmt.Println("\n9. Final status...")
	status = getKafkaStatus()
	printStatus(status)

	// 10. Enable consumer
	fmt.Println("\n10. Enabling Kafka consumer...")
	enableConsumer()

	fmt.Println("\n🎉 Demo completed!")
}

func getKafkaStatus() *KafkaStatus {
	resp, err := http.Get(baseURL + "/status")
	if err != nil {
		fmt.Printf("❌ Error getting status: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Error reading response: %v\n", err)
		return nil
	}

	var status KafkaStatus
	if err := json.Unmarshal(body, &status); err != nil {
		fmt.Printf("❌ Error parsing status: %v\n", err)
		return nil
	}

	return &status
}

func printStatus(status *KafkaStatus) {
	if status == nil {
		fmt.Println("❌ No status available")
		return
	}

	fmt.Printf("📊 Kafka Status:\n")
	fmt.Printf("   Producer Enabled: %v\n", status.Data["producer_enabled"])
	fmt.Printf("   Consumer Enabled: %v\n", status.Data["consumer_enabled"])
	fmt.Printf("   Consumer Count: %v\n", status.Data["consumer_count"])
	fmt.Printf("   Brokers: %v\n", status.Data["brokers"])
}

func sendMessage(message KafkaMessage) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling message: %w", err)
	}

	resp, err := http.Post(baseURL+"/producer/request", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]string
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return fmt.Errorf("error parsing error response: %w", err)
		}
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, errorResp["error"])
	}

	return nil
}

func enableProducer() {
	resp, err := http.Post(baseURL+"/producer/enable", "application/json", nil)
	if err != nil {
		fmt.Printf("❌ Error enabling producer: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ Producer enabled")
	} else {
		fmt.Printf("❌ Failed to enable producer: HTTP %d\n", resp.StatusCode)
	}
}

func disableProducer() {
	resp, err := http.Post(baseURL+"/producer/disable", "application/json", nil)
	if err != nil {
		fmt.Printf("❌ Error disabling producer: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ Producer disabled")
	} else {
		fmt.Printf("❌ Failed to disable producer: HTTP %d\n", resp.StatusCode)
	}
}

func enableConsumer() {
	resp, err := http.Post(baseURL+"/consumer/enable", "application/json", nil)
	if err != nil {
		fmt.Printf("❌ Error enabling consumer: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ Consumer enabled")
	} else {
		fmt.Printf("❌ Failed to enable consumer: HTTP %d\n", resp.StatusCode)
	}
}

func disableConsumer() {
	resp, err := http.Post(baseURL+"/consumer/disable", "application/json", nil)
	if err != nil {
		fmt.Printf("❌ Error disabling consumer: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ Consumer disabled")
	} else {
		fmt.Printf("❌ Failed to disable consumer: HTTP %d\n", resp.StatusCode)
	}
}
