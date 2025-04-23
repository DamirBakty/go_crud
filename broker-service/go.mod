module broker-service

go 1.24.2

require (
	auth-service v0.0.0
	blog-service v0.0.0
	google.golang.org/grpc v1.71.1
)

require (
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250414145226-207652e42e2e // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace auth-service => ../auth-service

replace blog-service => ../blog-service
