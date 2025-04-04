package mqtt

import (
	"fmt"
	"log"
	"os"
	// "regexp"
	// "strconv"
	// "strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"time"
)

type MqttClient struct {
	mqttClient mqtt.Client
}

func (sh *MqttClient) init() {
	sh.initMqttClient()
}

func (sh *MqttClient) initMqttClient() {
	//Message handler for incoming messages
	var messagePubHandler mqtt.MessageHandler = sh.messageHandler

	// Connection handler when connected
	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		fmt.Println("Connected to mqtt broker")
	}

	// Connection lost handler
	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		fmt.Printf("Connection lost: %v\n", err)
	}

	opts := mqtt.NewClientOptions().AddBroker("mqtt://" + os.Getenv("MOSQUITTO_IP") + ":31883")
	opts.SetClientID("stats_mqtt_client")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	// Create and start a new mqtt client
	sh.mqttClient = mqtt.NewClient(opts)
	if token := sh.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

func (sh *MqttClient) close() {
	time.Sleep(5 * time.Second)
	sh.mqttClient.Disconnect(250)
}

// func (sh *MqttClient) subscribeStats() {
// 	for key := range statsMap {
// 		sh.subscribeMessage(key)
// 	}
// }

func (sh *MqttClient) subscribeMessage(topic string) {
	fmt.Printf("Subscription: %s\n", topic)
	if token := sh.mqttClient.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

// func (sh *MqttClient) valueExtractor(input string) (string, string) {
// 	re := regexp.MustCompile(`([-+]?\d*\.\d+|[-+]?\d+|on|off)(°C|%|W|kWh|V|A|;ok)?`)
// 	var value, unit string

// 	matches := re.FindStringSubmatch(input)
// 	if len(matches) > 0 {
// 		value = matches[1]
// 		unit = matches[2]
// 		fmt.Printf("Value: %s, Unit: %s\n", value, unit)
// 	} else {
// 		fmt.Println("Value or unit not found")
// 	}

// 	return value, unit
// }

func (sh *MqttClient) messageHandler(client mqtt.Client, msg mqtt.Message) {

	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	// valueStr, _ := sh.valueExtractor(string(msg.Payload()))

	// valueStr = strings.ReplaceAll(valueStr, "on", "1")
	// valueStr = strings.ReplaceAll(valueStr, "off", "0")

	// dataPoint := statsMap[msg.Topic()]
	// value, _ := strconv.ParseFloat(valueStr, 64)
	// dataPoint.Fields = map[string]interface{}{"value": value}
	// sh.statsArchiver.saveStat(dataPoint)
}

// Funkcja do wysyłania komunikatu MQTT
func SendMessage(message string, topic string) {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883") // Zmien na odpowiedni adres brokera MQTT
	opts.SetClientID("home_automation")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	client.Publish(topic, 0, false, message)
	fmt.Printf("Sent message to MQTT topic %s: %s\n", topic, message)
}
