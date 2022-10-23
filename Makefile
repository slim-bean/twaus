

twaus-pizero:
	docker run --rm -v "${PWD}":/usr/src/twaus -w /usr/src/twaus slimbean/raspi-cross /bin/sh -c 'CC=/opt/cross-pi-gcc/bin/arm-linux-gnueabihf-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=6 go build -o cmd/twaus/twaus ./cmd/twaus/main.go'