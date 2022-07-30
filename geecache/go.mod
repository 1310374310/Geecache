module geecache

go 1.18

require lru v0.0.0
require singleflight v0.0.0

replace singleflight => ./singleflight
replace lru => ./lru