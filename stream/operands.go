package stream

func (s *Stream[T]) Filter(predicate Predicate[T]) *Stream[T] {
	newCloseChannel := make(chan struct{})
	filteredEmitter := buildEmitterFunc(newCloseChannel, s.baseSize, s.emitter.emitter, s.emitter.closeChannel, func(us chan T, t T) {
		if predicate(t) {
			us <- t
		}
	})
	return FromEmitter(
		&Emitter[T]{
			emitter:      filteredEmitter,
			closeChannel: newCloseChannel,
		},
		s.baseSize)
}

func Map[T, U any](s *Stream[T], mapping Mapper[T, U]) *Stream[U] {
	newCloseChannel := make(chan struct{})
	filteredEmitter := buildEmitterFunc(newCloseChannel, s.baseSize, s.emitter.emitter, s.emitter.closeChannel,
		func(c chan U, t T) {
			c <- mapping(t)
		})
	return FromEmitter(
		&Emitter[U]{
			emitter:      filteredEmitter,
			closeChannel: newCloseChannel,
		},
		s.baseSize)
}

func (s *Stream[T]) ForEach(consumer Consumer[T]) {
	emitterChannel := s.emitter.emitter()
	for {
		select {
		case e := <-emitterChannel:
			consumer(e)
		case <-s.emitter.closeChannel:
			queueSize := len(emitterChannel)
			for i := 0; i < queueSize; i++ {
				consumer(<-emitterChannel)
			}
			return
		}
	}
}

func (s *Stream[T]) ToSlice() []T {
	emitterChannel := s.emitter.emitter()
	var result []T
	for {
		select {
		case e := <-emitterChannel:
			result = append(result, e)
		case <-s.emitter.closeChannel:
			queueSize := len(emitterChannel)
			for i := 0; i < queueSize; i++ {
				result = append(result, <-emitterChannel)
			}
			return result
		}
	}
}
