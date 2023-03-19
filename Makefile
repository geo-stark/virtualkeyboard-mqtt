NAME=virtualkeyboard-mqtt

all:
	cd src; go build -o ../bin/$(NAME)

arm:
	# CC varibale should define cross cross-compiler
	export GOARCH=arm GOOS=linux CGO_ENABLED=1
	cd src;  go build -a -o ../bin/$(NAME) -ldflags="-extldflags=-static -w"

install:
	cp bin/$(NAME) /usr/local/sbin
	cp $(NAME).service /lib/systemd/system
	systemctl daemon-reload
