.PHONY: all

all: clean build

clean:
	@rm -rf bin public
	@packr2 clean
	@echo "[✔️] Clean complete!"

build:
	@cd ./web && yarn && yarn build && cd -
	@cp public/index.html public/static
	@packr2
	@go build -o bin/app
	@echo "[✔️] Build complete!"

run:
	@echo "[✔️] App is running!"
	@./bin/app

dev:
	@go run main.go

dev-web:
	@cd ./web && yarn serve
