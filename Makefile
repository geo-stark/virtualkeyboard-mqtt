NAME=virtualkeyboard-mqtt

all:
	go build -o bin/$(NAME) ./src

install:
	cp bin/$(NAME) /usr/local/sbin
	cp *.service /lib/systemd/system
	systemctl daemon-reload
    