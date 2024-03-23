
sample-50:
	go run ./cmd/ --file $(PWD)/internal/reader/sample_data/measurements_50.txt --workers-count 10 --queue-size 30

sample-1b:
	go run ./cmd/cli.go --file $(PWD)/internal/reader/sample_data/measurements_1b.txt --workers-count 10 --queue-size 30

sample-1k:
	go run ./cmd/cli.go --file $(PWD)/internal/reader/sample_data/measurements_1k.txt --workers-count 50 --queue-size 30

sample-1m:
	go run ./cmd/cli.go --file $(PWD)/internal/reader/sample_data/measurements_1m.txt --workers-count 6 --queue-size 24


trace-50:
	go run ./cmd/ --file $(PWD)/internal/reader/sample_data/measurements_50.txt --workers-count 10 --queue-size 30 --trace

trace-1b:
	go run ./cmd/cli.go --file $(PWD)/internal/reader/sample_data/measurements_1b.txt --workers-count 10 --queue-size 30 --trace

trace-1k:
	go run ./cmd/cli.go --file $(PWD)/internal/reader/sample_data/measurements_1k.txt --workers-count 50 --queue-size 30 --trace

trace-1m:
	go run ./cmd/cli.go --file $(PWD)/internal/reader/sample_data/measurements_1m.txt --workers-count 6 --queue-size 24 --trace

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

benchmark:
	go test -bench='.' -cpuprofile='cpu_wg.prof' -memprofile='mem_wg.prof' -bench='BenchmarkConsumeWG' ./internal/service/
	go test -bench='.' -cpuprofile='cpu_eg.prof' -memprofile='mem_eg.prof' -bench='BenchmarkConsumeEG' ./internal/service/
	go test -bench='.' -cpuprofile='cpu.prof' -memprofile='mem.prof' ./internal/service/

pprof:
	go tool pprof -http=0.0.0.0:12345 -show_from='service.Consume' cpu.prof

pprof-mem:
	go tool pprof -http=localhost:12345 -show_from='service.Consume' -diff_base=mem_wg.prof mem_eg.prof

pprof-cpu:
	go tool pprof -http=localhost:12345 -show_from='service.Consume' -diff_base=cpu_wg.prof cpu_eg.prof


compare:
	go tool pprof -text -hide=runtime -show=Consume.* cpu_wg.prof
	go tool pprof -text -hide=runtime -show=Consume.* cpu_eg.prof
	go tool pprof -text -hide=runtime -show=Consume.* cpu.prof
	go tool pprof -text -hide=runtime -show=Consume.* mem_wg.prof
	go tool pprof -text -hide=runtime -show=Consume.* mem_eg.prof
	go tool pprof -text -hide=runtime -show=Consume.* mem.prof
	go tool pprof -svg -output=cpu_wg.svg -hide=runtime -show=Consume.* cpu_wg.prof
	go tool pprof -svg -output=cpu_eg.svg -hide=runtime -show=Consume.* cpu_eg.prof
	go tool pprof -svg -output=cpu.svg -hide=runtime -show=Consume.* cpu.prof
	go tool pprof -svg -output=mem_wg.svg -hide=runtime -show=Consume.* mem_wg.prof
	go tool pprof -svg -output=mem_eg.svg -hide=runtime -show=Consume.* mem_eg.prof
	go tool pprof -svg -output=mem.svg -hide=runtime -show=Consume.* mem.prof
	go tool pprof -gif -output=cpu_wg.gif -hide=runtime -show=Consume.* cpu_wg.prof
	go tool pprof -gif -output=cpu_eg.gif -hide=runtime -show=Consume.* cpu_eg.prof
	go tool pprof -gif -output=cpu.gif -hide=runtime -show=Consume.* cpu.prof
	go tool pprof -gif -output=mem_wg.gif -hide=runtime -show=Consume.* mem_wg.prof
	go tool pprof -gif -output=mem_eg.gif -hide=runtime -show=Consume.* mem_eg.prof
	go tool pprof -gif -output=mem.gif -hide=runtime -show=Consume.* mem.prof

clean:
	-rm *.prof
	-rm *.svg
	-rm *.gif
