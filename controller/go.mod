module github.com/RuslanSemenikhin/Effective-Mobile-TestCase.git

go 1.24.3

require (
	github.com/RuslanSemenikhin/Effective-Mobile-TestCase.git/rpc v0.0.0
	github.com/joho/godotenv v1.5.1
	github.com/kelseyhightower/envconfig v1.4.0
	google.golang.org/grpc v1.74.2
)

require (
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250528174236-200df99c418a // indirect
	google.golang.org/protobuf v1.36.7 // indirect
)

replace github.com/RuslanSemenikhin/Effective-Mobile-TestCase.git/rpc => ../rpc
