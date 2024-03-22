
sample-50:
	go run ./cmd/ --file $(PWD)/internal/reader/sample_data/measurements_50.txt --workers-count 10 --queue-size 30 --trace

sample-1b:
	go run ./cmd/cli.go --file $(PWD)/internal/reader/sample_data/measurements_1b.txt --workers-count 10 --queue-size 30

sample-1k:
	go run ./cmd/cli.go --file $(PWD)/internal/reader/sample_data/measurements_1k.txt --workers-count 50 --queue-size 30 --trace

sample-1m:
	go run ./cmd/cli.go --file $(PWD)/internal/reader/sample_data/measurements_1m.txt --workers-count 50 --queue-size 75 --trace

build-docker:
	docker \
		buildx \
		build \
		--tag eldius/1brc:latest \
			.

run-docker: build-docker
	docker \
		run \
		-m 512m \
		--cpus=4 \
		-v "$(PWD)/internal/reader/sample_data/measurements_1b.txt:/app/measurements.txt:ro" \
		--rm \
			eldius/1brc:latest --file /app/measurements.txt --workers-count 5 --queue-size 10

run-docker-exportlogs:
	$(eval FLUENTBIT_HOST := $(shell ./fetch_ports.sh fluent-bit 24224 observability))
	$(eval COLLECTOR_TRACE_HOST := $(shell ./fetch_ports.sh otel-collector 55689 observability))
	$(eval COLLECTOR_METRICS_HOST := $(shell ./fetch_ports.sh otel-collector 55690 observability))
	$(eval DB_HOST := $(shell ./fetch_ports.sh postgres 5432 databases))

	@echo "FLUENTBIT_HOST:        $(FLUENTBIT_HOST)"
	@echo "COLLECTOR_TRACE_HOST:  $(COLLECTOR_TRACE_HOST)"
	@echo "COLLECTOR_TRACE_HOST:  $(COLLECTOR_TRACE_HOST)"
	@echo "DB_HOST:               $(DB_HOST)"

	docker \
		run \
		-m 512m \
		--cpus=4 \
		-v "$(PWD)/internal/reader/sample_data/measurements_1b.txt:/app/measurements.txt:ro" \
		--log-driver=fluentd \
		--log-opt fluentd-address=$(FLUENTBIT_HOST) \
		--rm \
			eldius/1brc:latest --file /app/measurements.txt --workers-count 5 --queue-size 10
