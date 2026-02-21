package or_test

import (
	"testing"
	"time"

	or "github.com/Oleska1601/WBOrChannel"
)

func setupChannels(delays []time.Duration) []<-chan interface{} {
	channels := make([]<-chan interface{}, len(delays))

	for i, delay := range delays {
		ch := make(chan interface{})
		go func(d time.Duration, c chan interface{}) {
			time.Sleep(d)
			close(c)
		}(delay, ch)
		channels[i] = ch
	}

	return channels
}

func TestOr(t *testing.T) {
	tests := []struct {
		name                string
		delays              []time.Duration
		expectedMaxDuration time.Duration
	}{
		{
			name:                "0 channels",
			delays:              []time.Duration{},
			expectedMaxDuration: 1 * time.Millisecond,
		},
		{
			name:                "1 channel",
			delays:              []time.Duration{100 * time.Millisecond},
			expectedMaxDuration: 150 * time.Millisecond,
		},
		{
			name:                "2 channels with first fast",
			delays:              []time.Duration{100 * time.Millisecond, 200 * time.Millisecond},
			expectedMaxDuration: 150 * time.Millisecond,
		},
		{
			name:                "2 channels with last fast",
			delays:              []time.Duration{300 * time.Millisecond, 200 * time.Millisecond},
			expectedMaxDuration: 250 * time.Millisecond,
		},
		{
			name:                "3 channels",
			delays:              []time.Duration{300 * time.Millisecond, 100 * time.Millisecond, 200 * time.Millisecond},
			expectedMaxDuration: 150 * time.Millisecond,
		},
		{
			name: "5 channels with first fast",
			delays: []time.Duration{100 * time.Millisecond, 200 * time.Millisecond,
				300 * time.Millisecond, 400 * time.Millisecond, 500 * time.Millisecond},
			expectedMaxDuration: 150 * time.Millisecond,
		},
		{
			name: "5 channels with last fast",
			delays: []time.Duration{500 * time.Millisecond, 400 * time.Millisecond,
				300 * time.Millisecond, 200 * time.Millisecond, 10 * time.Millisecond},
			expectedMaxDuration: 150 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			channels := setupChannels(tt.delays)
			start := time.Now()

			got := or.Or(channels...)
			<-got

			duration := time.Since(start)
			if duration > tt.expectedMaxDuration {
				t.Errorf("Or() duration = %v, want less than %v", duration, tt.expectedMaxDuration)
			}
		})
	}
}

// тест на панику с nil каналом
func TestOrWithNilChannel(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Or panicked with nil channel: %v", r)
		}
	}()

	ch := make(chan interface{})
	close(ch)

	got := or.Or(ch, nil)

	// должен закрыться сразу (ch уже закрыт)
	select {
	case <-got:
		// OK
	case <-time.After(100 * time.Millisecond):
		t.Error("Or with nil channel timed out")
	}
}
