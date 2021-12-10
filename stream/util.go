package stream

func buildEmitterFunc[U, T any](closeChannel chan struct{}, baseSize uint32, prevEmitter EmitterFunction[T], prevCloseChan chan struct{}, processingFunction func(chan U, T)) EmitterFunction[U] {
	return func() chan U {
		newChannel := make(chan U, baseSize)
		go func(c chan T) {
			defer func() { closeChannel <- struct{}{} }()
			for {
				select {
				case e := <-c:
					processingFunction(newChannel, e)
				case <-prevCloseChan:
					queueSize := len(c)
					for i := 0; i < queueSize; i++ {
						processingFunction(newChannel, <-c)
					}
					return
				}
			}
		}(prevEmitter())
		return newChannel
	}
}
