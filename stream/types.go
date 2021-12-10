package stream

type Stream[T any] struct {
	emitter  *Emitter[T]
	baseSize uint32
}

type Emitter[T any] struct {
	emitter      EmitterFunction[T]
	closeChannel chan struct{}
}

type Predicate[T any] func(T) bool
type Consumer[T any] func(T)
type Mapper[T any, U any] func(T) U
type EmitterFunction[T any] func() chan T

func FromIterable[I any](iterable []I) *Stream[I] {
	finished := make(chan struct{})
	emitter := func() chan I {
		ch := make(chan I, len(iterable))
		go func() {
			for _, i := range iterable {
				ch <- i
			}
			finished <- struct{}{}
		}()
		return ch
	}
	return &Stream[I]{
		emitter: &Emitter[I]{
			emitter:      emitter,
			closeChannel: finished,
		},
		baseSize: uint32(len(iterable)),
	}
}

func FromEmitter[I any](emitter *Emitter[I], baseSize ...uint32) *Stream[I] {
	size := uint32(10)
	if len(baseSize) > 0 {
		size = baseSize[0]
	}
	return &Stream[I]{
		emitter:  emitter,
		baseSize: size,
	}
}