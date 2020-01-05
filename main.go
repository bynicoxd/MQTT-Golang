package main

import (
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//Broker for mqtt
type Broker struct {
	host     string
	port     string
	username string
	password string
	clientID string
	QoS      byte
	protocol string
	topic    string
	message  string
}

func main() {
	//Define parameters of the connection
	var br Broker
	br.host = "test.mosquitto.org"
	br.port = "1883"
	// this is only if you have your own broker
	/*
		br.username = "user"
		br.password = "pass"
	*/
	br.QoS = byte(0)
	// in the protocol you can set ssl or tcp

	br.protocol = "tcp"
	br.topic = "/yourTopic"
	br.message = "message"

	mode := ""
	// for run this code, you must set the values of the connection, and run it into the command line passing the mode 'sub' to suscribe
	// or 'pub' to publish
	// in two cases you must set only the clientID
	// in this format ./main mode clientID

	if len(os.Args) > 1 {
		mode = os.Args[1]
		br.clientID = os.Args[2]
	} else {
		fmt.Println("Insert a option to start...")
		os.Exit(0)
	}
	client, err := inicialize(br)
	if err != nil {
		panic(err)
	}
	switch mode {
	case "sub":
		susbcriber(client, br)

	case "pub":
		publisher(client, br)
	}

}
func inicialize(br Broker) (mqtt.Client, error) {
	fmt.Println("MQTT Example")
	opts := mqtt.NewClientOptions()
	host := (br.protocol + "://" + br.host + ":" + br.port)
	fmt.Println(host)
	opts.AddBroker(host)
	opts.SetUsername(br.username)
	opts.SetPassword(br.password)
	opts.SetClientID(br.clientID)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	err := token.Error()
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
		return nil, err
	} else {
		fmt.Println(br.clientID + " connected\n ")
		return client, nil
	}

}

func susbcriber(client mqtt.Client, br Broker) {
	for {
		client.Subscribe("#", br.QoS, func(client mqtt.Client, msg mqtt.Message) {
			fmt.Printf(" Listening from topic %s message :%s\n", msg.Topic(), string(msg.Payload()))
		})
		time.Sleep(1 * time.Second)
	}

}

func publisher(client mqtt.Client, br Broker) {

	for {
		client.Publish(br.topic, br.QoS, false, br.message)
		fmt.Printf(" Sending: %s  to topic %s from the publisher \n", br.message, br.topic)
		time.Sleep(3 * time.Second)
	}
}
