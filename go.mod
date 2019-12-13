module goim_pro

go 1.13

require (
	github.com/golang/protobuf v1.3.2
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553
	golang.org/x/sys v0.0.0-20191210023423-ac6580df4449 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20191206224255-0243a4be9c8f // indirect
	google.golang.org/grpc v1.25.1
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.25.1
