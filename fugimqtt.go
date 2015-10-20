package main

import (
	"fmt"
	"log"
	"os"
	"time"
	// Paho MQTT library
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

var f MQTT.MessageHandler = func(client *MQTT.Client, msg MQTT.Message) {
	fmt.Println("--- NEW MQTT MESSAGE ---")
	fmt.Printf("   TOPIC: %s\n", msg.Topic())
	fmt.Printf("     MSG: %s\n", msg.Payload())
}

var c *MQTT.Client

func updateColor(target, color string, config Config) {
	topic := config.MQTT.Topic + "/color"
	log.Println(target)
	switch target {
	case "all":
		token := c.Publish(topic, 0, false, color)
		token.Wait()
	case "taps", "tower":
		topic += "/tower"
		c.Publish(topic, 0, false, color)
	case "shelves":
		c.Publish(topic+"/rshelf", 0, false, color)
		c.Publish(topic+"/lshelf", 0, false, color)
		c.Publish(topic+"/rear", 0, false, color)
	default:
		log.Println("not caught")
	}
}

// SetupMQTT : connect to server and register callbacks
func setupMQTT(config Config) {
	fmt.Printf("username : %s", config.MQTT.Username)
	if pass := config.MQTT.Password; pass == "" {
		fmt.Println("No password given")
	} else {
		fmt.Printf("password: %s\n", config.MQTT.Password)
	}

	var server = fmt.Sprintf("tcp://%s:%s", config.MQTT.Host, config.MQTT.Port)
	opts := MQTT.NewClientOptions().AddBroker(server)
	opts.SetUsername(config.MQTT.Username)
	opts.SetPassword(config.MQTT.Password)
	opts.SetClientID(config.MQTT.ClientID)

	c = MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("connected to %s\n", server)
	}
	//defer c.Disconnect(250)

	if token := c.Subscribe(config.MQTT.Topic, 0, f); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	text := fmt.Sprintf("Hello from %s who just connected!", config.MQTT.ClientID)
	token := c.Publish(config.MQTT.Topic, 0, false, text)
	token.Wait()
	time.Sleep(2 * time.Second)
}
