/*
 * Copyright (c) 2021 IBM Corp and others.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v2.0
 * and Eclipse Distribution License v1.0 which accompany this distribution.
 *
 * The Eclipse Public License is available at
 *    https://www.eclipse.org/legal/epl-2.0/
 * and the Eclipse Distribution License is available at
 *   http://www.eclipse.org/org/documents/edl-v10.php.
 *
 * Contributors:
 *    Seth Hoenig
 *    Allan Stockdill-Mander
 *    Mike Robertson
 */

package main

import (
	"flag"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"os"
)

/*
Options:

	[-help]                      Display help
	[-a pub|sub]                 Action pub (publish) or sub (subscribe)
	[-m <message>]               Payload to send
	[-n <number>]                Number of messages to send or receive
	[-q 0|1|2]                   Quality of Service
	[-clean]                     CleanSession (true if -clean is present)
	[-id <clientid>]             CliendID
	[-user <user>]               User
	[-password <password>]       Password
	[-broker <uri>]              Broker URI
	[-topic <topic>]             Topic
	[-store <path>]              Store Directory
*/
var (
	client MQTT.Client
)

func startMQTT() {
	topic := flag.String("topic", "robert/robertstestsensor/message/testmessage", "The topic name to/from which to publish/subscribe")
	broker := flag.String("broker", "tcp://test.mosquitto.org:1884", "The broker URI. ex: tcp://10.10.1.1:1883")
	password := flag.String("password", "readwrite", "The password (optional)")
	user := flag.String("user", "rw", "The User (optional)")
	id := flag.String("id", "robertssupercoolclientidthatisnotused", "The ClientID (optional)")
	cleansess := flag.Bool("clean", true, "Set Clean Session (default false)")
	qos := flag.Int("qos", 2, "The Quality of Service 0,1,2 (default 0)")
	num := flag.Int("num", 1, "The number of messages to publish or subscribe (default 1)")
	payload := flag.String("message", "{\"test\":\"robert\"}", "The message text to publish (default empty)")
	action := flag.String("action", "pub", "Action publish or subscribe (required)")
	store := flag.String("store", ":memory:", "The Store Directory (default use memory store)")
	flag.Parse()

	if *action != "pub" && *action != "sub" {
		fmt.Println("Invalid setting for -action, must be pub or sub")
		return
	}

	if *topic == "" {
		fmt.Println("Invalid setting for -topic, must not be empty")
		return
	}

	fmt.Printf("Sample Info:\n")
	fmt.Printf("\taction:    %s\n", *action)
	fmt.Printf("\tbroker:    %s\n", *broker)
	fmt.Printf("\tclientid:  %s\n", *id)
	fmt.Printf("\tuser:      %s\n", *user)
	fmt.Printf("\tpassword:  %s\n", *password)
	fmt.Printf("\ttopic:     %s\n", *topic)
	fmt.Printf("\tmessage:   %s\n", *payload)
	fmt.Printf("\tqos:       %d\n", *qos)
	fmt.Printf("\tcleansess: %v\n", *cleansess)
	fmt.Printf("\tnum:       %d\n", *num)
	fmt.Printf("\tstore:     %s\n", *store)

	opts := MQTT.NewClientOptions()
	opts.AddBroker(*broker)
	// opts.SetClientID(*id)
	opts.SetUsername(*user)
	opts.SetPassword(*password)
	opts.SetCleanSession(*cleansess)
	if *store != ":memory:" {
		opts.SetStore(MQTT.NewFileStore(*store))
	}
	connectMQTT(opts, action, num, payload, topic, qos)
}

func connectMQTT(opts *MQTT.ClientOptions, action *string, num *int, payload *string, topic *string, qos *int) {
	if *action == "pub" {
		client = MQTT.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println("Sample Publisher Started")
		for i := 0; i < *num; i++ {
			fmt.Println("---- doing publish ----")
			token := client.Publish(*topic, byte(*qos), false, *payload)
			token.Wait()
		}

		// client.Disconnect(250)
		// fmt.Println("Sample Publisher Disconnected")
	} else {
		receiveCount := 0
		choke := make(chan [2]string)

		opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
			choke <- [2]string{msg.Topic(), string(msg.Payload())}
		})

		client = MQTT.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

		if token := client.Subscribe(*topic, byte(*qos), nil); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}

		for receiveCount < *num {
			incoming := <-choke
			fmt.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])
			receiveCount++
		}

		client.Disconnect(250)
		fmt.Println("Sample Subscriber Disconnected")
	}
}
