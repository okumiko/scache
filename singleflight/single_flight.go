package singleflight

type SingleFlight interface {
	Do(key string, fn func() (interface{}, error)) (interface{}, error, bool)
}
