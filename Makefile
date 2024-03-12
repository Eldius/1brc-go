
sample-50:
	go run ./cmd/

sample-1b:
	go run ./cmd/ --file ./internal/parser/sample_data/measurements_1b.txt --workers-count 10

sample-1k:
	go run ./cmd/ --file ./internal/parser/sample_data/measurements_1k.txt --workers-count 10

build-docker:
	docker \
		buildx \
		build \
		--tag eldius/1brc:latest \
			.

run-docker: build-docker
#	$(eval FLUENTBIT_HOST := $(shell ./fetch_ports.sh fluent-bit 24224 observability))
#	$(eval COLLECTOR_TRACE_HOST := $(shell ./fetch_ports.sh otel-collector 55689 observability))
#	$(eval COLLECTOR_METRICS_HOST := $(shell ./fetch_ports.sh otel-collector 55690 observability))
#	$(eval DB_HOST := $(shell ./fetch_ports.sh postgres 5432 databases))
#
#	@echo "FLUENTBIT_HOST:        $(FLUENTBIT_HOST)"
#	@echo "COLLECTOR_TRACE_HOST:  $(COLLECTOR_TRACE_HOST)"
#	@echo "COLLECTOR_TRACE_HOST:  $(COLLECTOR_TRACE_HOST)"
#	@echo "DB_HOST:               $(DB_HOST)"

#	docker \
#		run \
#		-m 512m \
#		--cpus=4 \
#		-v "$(PWD)/internal/parser/sample_data/measurements_1b.txt:/app/measurements.txt:ro" \
#		--log-driver=fluentd \
#		--log-opt fluentd-address=$(FLUENTBIT_HOST) \
#		--rm \
#			eldius/1brc:latest --file /app/measurements.txt --workers-count 5 --queue-size 10

	docker \
		run \
		-m 512m \
		--cpus=4 \
		-v "$(PWD)/internal/parser/sample_data/measurements_1b.txt:/app/measurements.txt:ro" \
		--rm \
			eldius/1brc:latest --file /app/measurements.txt --workers-count 5 --queue-size 10
