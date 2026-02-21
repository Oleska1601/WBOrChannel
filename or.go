package or

// Or объединяет несколько каналов завершения (done-каналов) в один.
// Возвращаемый канал закрывается, как только закрывается любой из переданных каналов.
func Or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		done := make(chan interface{})
		close(done)
		return done
	case 1:
		return channels[0]
	case 2:
		done := make(chan interface{})
		go func() {
			defer close(done)

			select {
			case <-channels[0]:
			case <-channels[1]:
			}

		}()

		return done
	default:
		k := len(channels) / 2
		return Or(Or(channels[:k]...), Or(channels[k:]...))
	}
}
