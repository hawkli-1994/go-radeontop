
build-example-for-riscv:
	GOOS=linux GOARCH=riscv64 go build -o examples/basic_monitoring/basic_monitoring_riscv examples/basic_monitoring/main.go

build-example-for-amd64:
	GOOS=linux GOARCH=amd64 go build -o examples/basic_monitoring/basic_monitoring_amd64 examples/basic_monitoring/main.go

