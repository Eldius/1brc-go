
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
	docker \
		run \
		-m 512m \
		--cpus=4 \
		-v "$(PWD)/internal/parser/sample_data/measurements_1b.txt:/app/measurements.txt:ro" \
		--rm \
			eldius/1brc:latest --file /app/measurements.txt --workers-count 5 --queue-size 10
