

.PHONY: get-protobuf protoc-gen

build: get-protobuf protoc-gen build-exporter

clean:
	rm -rf build
get-protobuf:
	mkdir build
	grpcurl -plaintext -protoset-out build/dish.protoset 192.168.100.1:9200 describe SpaceX.API.Device.Device
protoc-gen:
	protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=build/dish.protoset spacex/api/device/device.proto
	protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=build/dish.protoset spacex/api/common/status/status.proto
	protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=build/dish.protoset spacex/api/device/command.proto
	protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=build/dish.protoset spacex/api/device/common.proto
	protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=build/dish.protoset spacex/api/device/dish.proto
	protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=build/dish.protoset spacex/api/device/wifi.proto
	protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=build/dish.protoset spacex/api/device/wifi_config.proto
	protoc --go_out=pkg/    --go-grpc_out=pkg/  --descriptor_set_in=build/dish.protoset spacex/api/device/transceiver.proto

	find pkg/spacex.com -name "*.go" | xargs sed -i.bak 's|spacex.com/api|github.com/dev-mull/starlink_exporter/pkg/spacex.com/api|g'
	find pkg/spacex.com -name "*.bak" | xargs rm
build-exporter:
	GOOS=linux GOARCH=amd64 go build -o build/starlink_exporter-linux-amd64


