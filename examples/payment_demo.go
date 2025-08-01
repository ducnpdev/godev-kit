package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ducnpdev/godev-kit/internal/entity"
)

// PaymentDemo demonstrates the payment system
func PaymentDemo() {
	fmt.Println("=== Payment System Demo ===")

	// Demo 1: Register a payment
	fmt.Println("\n1. Registering a payment...")
	paymentReq := map[string]interface{}{
		"user_id":        1,
		"amount":         500000,
		"currency":       "VND",
		"payment_type":   "electric",
		"meter_number":   "EVN001234567",
		"customer_code":  "CUST001",
		"description":    "Thanh toán tiền điện tháng 12/2024",
		"payment_method": "bank_transfer",
	}

	reqBody, _ := json.Marshal(paymentReq)
	resp, err := http.Post("http://localhost:8080/api/v1/payments", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("Error registering payment: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		var paymentResp entity.PaymentResponse
		json.NewDecoder(resp.Body).Decode(&paymentResp)
		fmt.Printf("Payment registered successfully!\n")
		fmt.Printf("Payment ID: %d\n", paymentResp.ID)
		fmt.Printf("Transaction ID: %s\n", paymentResp.TransactionID)
		fmt.Printf("Status: %s\n", paymentResp.Status)
	} else {
		fmt.Printf("Failed to register payment. Status: %d\n", resp.StatusCode)
	}

	// Demo 2: Get payment by ID
	fmt.Println("\n2. Getting payment by ID...")
	resp, err = http.Get("http://localhost:8080/api/v1/payments/1")
	if err != nil {
		log.Printf("Error getting payment: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var payment entity.PaymentResponse
		json.NewDecoder(resp.Body).Decode(&payment)
		fmt.Printf("Payment details:\n")
		fmt.Printf("  ID: %d\n", payment.ID)
		fmt.Printf("  Amount: %.2f %s\n", payment.Amount, payment.Currency)
		fmt.Printf("  Status: %s\n", payment.Status)
		fmt.Printf("  Meter Number: %s\n", payment.MeterNumber)
		fmt.Printf("  Customer Code: %s\n", payment.CustomerCode)
		fmt.Printf("  Created At: %s\n", payment.CreatedAt.Format(time.RFC3339))
	} else {
		fmt.Printf("Failed to get payment. Status: %d\n", resp.StatusCode)
	}

	// Demo 3: Get payments by user ID
	fmt.Println("\n3. Getting payments by user ID...")
	resp, err = http.Get("http://localhost:8080/api/v1/users/1/payments")
	if err != nil {
		log.Printf("Error getting user payments: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var payments []entity.PaymentResponse
		json.NewDecoder(resp.Body).Decode(&payments)
		fmt.Printf("Found %d payments for user:\n", len(payments))
		for i, payment := range payments {
			fmt.Printf("  %d. Payment ID: %d, Amount: %.2f %s, Status: %s\n",
				i+1, payment.ID, payment.Amount, payment.Currency, payment.Status)
		}
	} else {
		fmt.Printf("Failed to get user payments. Status: %d\n", resp.StatusCode)
	}

	// Demo 4: Wait for payment processing
	fmt.Println("\n4. Waiting for payment processing...")
	time.Sleep(5 * time.Second)

	// Check payment status again
	resp, err = http.Get("http://localhost:8080/api/v1/payments/1")
	if err != nil {
		log.Printf("Error getting payment: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var payment entity.PaymentResponse
		json.NewDecoder(resp.Body).Decode(&payment)
		fmt.Printf("Payment status after processing: %s\n", payment.Status)
		if payment.Status == "completed" {
			fmt.Println("✅ Payment processed successfully!")
		} else if payment.Status == "failed" {
			fmt.Println("❌ Payment processing failed!")
		} else {
			fmt.Printf("⏳ Payment is still %s\n", payment.Status)
		}
	}

	fmt.Println("\n=== Demo completed ===")
}

// KafkaConsumerDemo demonstrates Kafka consumer
func KafkaConsumerDemo() {
	fmt.Println("\n=== Kafka Consumer Demo ===")
	fmt.Println("Starting payment consumer...")
	fmt.Println("Consumer will process payment events from Kafka topic 'payment-events'")
	fmt.Println("Press Ctrl+C to stop the consumer")

	// In a real application, you would start the consumer here
	// consumer := payment.NewPaymentConsumer(brokers, groupID, useCase, logger)
	// consumer.Start(context.Background())
}

// RunPaymentDemo runs the payment demo
func RunPaymentDemo() {
	// Run payment demo
	PaymentDemo()

	// Run Kafka consumer demo
	KafkaConsumerDemo()
}
