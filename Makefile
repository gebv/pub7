build:
	CGO_ENABLED=0 go build \
			-o bin/tgb \
			-v \
			-ldflags "-X main.VERSION=`date -u +%Y%m%d.%H%M%S`" \
			-a ./main.go
.PHONY: build

builddocker:
	docker-compose up -d buildapp

admin-run:
	go run main.go -file scripts/ff_script.toml admin-run -path_views views/admin/*.tpl

tgbot:
	go run main.go -file scripts/ff_script.toml tgbot-run