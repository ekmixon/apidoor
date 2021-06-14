module gateway

go 1.16

require (
	github.com/go-chi/chi/v5 v5.0.3
	github.com/go-redis/redis/v8 v8.9.0
	github.com/lib/pq v1.10.2
	local.packages/gateway v0.0.0-00010101000000-000000000000
)

replace local.packages/gateway => ./