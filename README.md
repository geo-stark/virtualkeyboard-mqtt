###Description:
This is MQTT client that emulate keyboard events using linux uinput module (/dev/uinput).

Default topic is virtualkeyboard/emit. Payload is a list of keys separated by comma. All keys in the list will be pressed then released along.
Supported keys listed at the top of keyboard.go
Other keys can be added from <linux/input-event-codes.h>
To get the code of necessary key run `sudo showkey`

###Dependencies:	
	go get github.com/eclipse/paho.mqtt.golang
	go get github.com/pborman/getopt/v2
	go get github.com/rs/xid

###Build:
	make
	sudo make install

###Usage:
    virtualkeyboard-mqtt --url tcp://host:port
	run as root to access /dev/uinput

	or run as systemd service:
	sudo systemctl enable virtualkeyboard-mqtt
	sudo systemctl start virtualkeyboard-mqtt

###Test:
	mosquitto_pub -h 127.0.0.1 -t virtualkeyboard/emit -m 'ctrl,alt,del'
	mosquitto_pub -h 127.0.0.1 -t virtualkeyboard/emit -m 'f1'	
 
###License:
	MIT
