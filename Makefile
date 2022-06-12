migrate_up=go run main.go migrate --direction=up --step=0
migrate_down=go run main.go migrate --direction=down --step=0

check-cognitive-complexity:
	-gocognit -over 15 .

lint: check-cognitive-complexity
	golangci-lint run --print-issued-lines=false --exclude-use-default=false --enable=revive --enable=goimports  --enable=unconvert --enable=unparam --concurrency=2

check-modd-exists:
	@modd --version > /dev/null

run: check-modd-exists
	@modd -f ./.modd/server.modd.conf

test-only:
	SVC_ENV=test SVC_DISABLE_CACHING=true $(test_command)

test: lint test-only

migrate:
	@if [ "$(DIRECTION)" = "" ] || [ "$(STEP)" = "" ]; then\
    	$(migrate_up);\
	else\
		go run main.go migrate --direction=$(DIRECTION) --step=$(STEP);\
    fi
