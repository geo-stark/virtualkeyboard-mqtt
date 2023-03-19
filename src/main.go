package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	virtualkeyboard "virtualkeyboard-mqtt/modules/virtualkeyboard"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	getopt "github.com/pborman/getopt/v2"
	"github.com/rs/xid"
)

func main() {
	argHelp := getopt.BoolLong("help", 'h', "print help")
	argURL := getopt.StringLong("url", 'u', "tcp://127.0.0.1:1883",
		"broker url, e.g. tcp://127.0.0.1:8183")
	argTopic := getopt.StringLong("topic", 't', "virtualkeyboard/emit",
		"mqtt topic")

	getopt.Parse()
	if *argHelp {
		getopt.PrintUsage(os.Stdout)
		return
	}

	if err := virtualkeyboard.OpenEx(&virtualkeyboard.Options{}); err != nil {
		log.Fatalf("virtual keyboard, open error: %v", err)
	}
	defer virtualkeyboard.Close()

	opts := mqtt.NewClientOptions()
	opts.AddBroker(*argURL)
	opts.SetMaxReconnectInterval(24 * time.Hour)

	uid := "virtualkeyboard-" + xid.New().String()
	opts.SetClientID(uid)

	choke := make(chan [2]string)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		choke <- [2]string{msg.Topic(), string(msg.Payload())}
	})

	opts.SetOnConnectHandler(func(client mqtt.Client) {
		log.Printf("broker: connected")
		token := client.Subscribe(*argTopic, byte(2), nil)
		if token.Wait() && token.Error() != nil {
			log.Printf("broker: subscribe failed: %v", token.Error())
		}
	})

	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Printf("broker: connection lost: %v", err)
	})

	log.Printf("connecting to %v", opts.Servers[0])
	client := mqtt.NewClient(opts)
	for !client.IsConnected() {
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			time.Sleep(opts.PingTimeout)
		}
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	receiveCount := 0

	for {
		select {
		case incoming := <-choke:
			log.Printf("#%d received: %s\n", receiveCount, incoming[1])
			keys := strings.Split(incoming[1], ",")
			if err := virtualkeyboard.Emit(keys); err != nil {
				log.Printf("virtual keyboard, emit failed: %v", err)
			}
			receiveCount++
			continue
		case sig := <-sigs:
			log.Printf("got signal: %v", sig)
			break
		}
		break
	}
	log.Printf("disconnecting")
	client.Disconnect(250)
}
