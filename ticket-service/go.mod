module github.com/eventhub/ticket-service

go 1.23.0

require (
	github.com/eventhub/pkg v0.0.0
	github.com/eventhub/proto v0.0.0
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.68.0
	gorm.io/driver/postgres v1.5.9
	gorm.io/gorm v1.25.12
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241104194629-dd2ea8efbc28 // indirect
	google.golang.org/protobuf v1.36.8 // indirect
)

replace (
	github.com/eventhub/pkg => ../pkg
	github.com/eventhub/proto => ../proto
)
