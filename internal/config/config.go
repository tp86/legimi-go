package config

type ConfigFn[T any] func(T) T

func New[T any, C ConfigFn[T]](cfgs ...ConfigFn[T]) T {
	var inst T
	for _, cfg := range cfgs {
		inst = cfg(inst)
	}
	return inst
}
