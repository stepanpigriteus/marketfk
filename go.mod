module marketfuck

go 1.22

require (
	github.com/lib/pq v1.10.9
	github.com/prometheus/client_golang v1.19.0 // Совместима с Go 1.22
	github.com/redis/go-redis/v9 v9.12.1
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect; Понизили для совместимости
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/prometheus/client_model v0.5.0 // indirect; Понизили для совместимости
	github.com/prometheus/common v0.48.0 // indirect; Понизили для совместимости
	github.com/prometheus/procfs v0.12.0 // indirect; Понизили для совместимости
	golang.org/x/sys v0.16.0 // indirect; Понизили для совместимости
	google.golang.org/protobuf v1.33.0 // indirect; Понизили для совместимости
)
