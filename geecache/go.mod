module geecache

go 1.18

require lru v0.0.0

require singleflight v0.0.0

require (
	github.com/golang/protobuf v1.5.2 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

replace singleflight => ./singleflight

replace lru => ./lru
