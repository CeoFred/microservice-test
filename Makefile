build:
	cd src/cmd && go build -o service

tidy:
	cd src && go mod tidy

start:
	export SERVICE_PORT=3002 && export SERVICE_BIND=127.0.0.1 && ./service