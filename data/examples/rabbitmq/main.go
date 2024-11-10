// main.go
package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strings"
)

var amazonRootCA1 = `-----BEGIN CERTIFICATE-----
MIIDQTCCAimgAwIBAgITBmyfz5m/jAo54vB4ikPmljZbyjANBgkqhkiG9w0BAQsF
ADA5MQswCQYDVQQGEwJVUzEPMA0GA1UEChMGQW1hem9uMRkwFwYDVQQDExBBbWF6
b24gUm9vdCBDQSAxMB4XDTE1MDUyNjAwMDAwMFoXDTM4MDExNzAwMDAwMFowOTEL
MAkGA1UEBhMCVVMxDzANBgNVBAoTBkFtYXpvbjEZMBcGA1UEAxMQQW1hem9uIFJv
b3QgQ0EgMTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALJ4gHHKeNXj
ca9HgFB0fW7Y14h29Jlo91ghYPl0hAEvrAIthtOgQ3pOsqTQNroBvo3bSMgHFzZM
9O6II8c+6zf1tRn4SWiw3te5djgdYZ6k/oI2peVKVuRF4fn9tBb6dNqcmzU5L/qw
IFAGbHrQgLKm+a/sRxmPUDgH3KKHOVj4utWp+UhnMJbulHheb4mjUcAwhmahRWa6
VOujw5H5SNz/0egwLX0tdHA114gk957EWW67c4cX8jJGKLhD+rcdqsq08p8kDi1L
93FcXmn/6pUCyziKrlA4b9v7LWIbxcceVOF34GfID5yHI9Y/QCB/IIDEgEw+OyQm
jgSubJrIqg0CAwEAAaNCMEAwDwYDVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8EBAMC
AYYwHQYDVR0OBBYEFIQYzIU07LwMlJQuCFmcx7IQTgoIMA0GCSqGSIb3DQEBCwUA
A4IBAQCY8jdaQZChGsV2USggNiMOruYou6r4lK5IpDB/G/wkjUu0yKGX9rbxenDI
U5PMCCjjmCXPI6T53iHTfIUJrU6adTrCC2qJeHZERxhlbI1Bjjt/msv0tadQ1wUs
N+gDS63pYaACbvXy8MWy7Vu33PqUXHeeE6V/Uq2V8viTO96LXFvKWlJbYK8U90vv
o/ufQJVtMVT8QtPHRh8jrdkPSHCa2XV4cdFyQzR1bldZwgJcJmApzyMZFo6IQ6XU
5MsI+yMRQ+hDKXJioaldXgjUkK642M4UwtBV8ob2xJNDd2ZhwLnoQdeXeGADbkpy
rqXRfboQnoZsG4q5WTP468SQvvG5
-----END CERTIFICATE-----`

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func loadAmazonCACert() *x509.CertPool {
	// Create a new certificate pool and add the Amazon CA certificate
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM([]byte(amazonRootCA1))
	return pool
}

func main() {
	// Connect to RabbitMQ server
	user := os.Getenv("RABBITMQ_USER")
	pass := os.Getenv("RABBITMQ_PASS")
	host := os.Getenv("RABBITMQ_HOST")
	log.Println("RABBITMQ_HOST: ", host)
	uri := ""
	var conn *amqp.Connection
	var err error
	if strings.Contains(host, "amazonaws.com") {
		// AWS managed
		awsPref := "amqps://"
		uri = strings.TrimPrefix(host, awsPref)
		uri = fmt.Sprintf("%s%s:%s@%s", awsPref, user, pass, uri)
		log.Println("Connecting to RabbitMQ: ", uri)
		conn, err = amqp.DialTLS(uri, &tls.Config{
			RootCAs: loadAmazonCACert(),
		})
	} else {
		uri = fmt.Sprintf("amqp://%s:%s@%s", user, pass, host)
		log.Println("Connecting to RabbitMQ: ", uri)
		conn, err = amqp.Dial(uri)
	}

	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Create a channel
	log.Println("Creating a channel...")
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare a queue
	log.Println("Declaring a queue...")
	q, err := ch.QueueDeclare(
		"hello", // queue name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Publish a message to the queue
	log.Println("Publishing a message to the queue...")
	body := "Hello, RabbitMQ!"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	fmt.Printf(" [x] Sent %s\n", body)

	// Consume messages from the queue
	log.Println("Consuming messages from the queue...")
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	// Receive and print messages
	go func() {
		for d := range msgs {
			fmt.Printf(" [x] Received %s\n", d.Body)
		}
	}()

	fmt.Println(" [*] Waiting for messages. To exit, press CTRL+C")
	select {}
}
