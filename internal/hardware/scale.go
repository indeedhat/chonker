package hardware

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/indeedhat/dotenv"
	"periph.io/x/devices/v3/hx711"
	"periph.io/x/host/v3/rpi"
)

const readTimeout = time.Millisecond * 500

var (
	adjustToGrams dotenv.Int = "ADJUST_TO_G"
)

type ringBuffer struct {
	data []int
	i    int
	full bool
	mux  sync.Mutex
}

func newRingBuffer() *ringBuffer {
	return &ringBuffer{data: make([]int, 11)}
}

func (b *ringBuffer) Set(v int32) {
	b.mux.Lock()
	defer b.mux.Unlock()

	b.i++

	if b.i == 10 {
		b.full = true
	}

	i := b.i % 10
	b.data[i] = int(v)
}

func (b *ringBuffer) Data() []int {
	b.mux.Lock()
	defer b.mux.Unlock()

	out := make([]int, 11)
	copy(out, b.data)

	return out
}

type scale struct {
	hx     *hx711.Dev
	buf    *ringBuffer
	cancel context.CancelFunc
	Zero   float64
	Scale  float64
}

func NewScale() (*scale, error) {
	hx, err := hx711.New(
		rpi.P1_29,
		rpi.P1_31,
	)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	scale := &scale{
		hx:     hx,
		buf:    newRingBuffer(),
		cancel: cancel,
		Scale:  float64(adjustToGrams.Get(1)),
	}

	go func() {
		readings := scale.hx.ReadContinuous()

		for {
			select {
			case <-ctx.Done():
				return
			case r := <-readings:
				if r.Raw == -1 {
					continue
				}
				scale.buf.Set(r.Raw)
			}
		}
	}()

	return scale, nil
}

func (s *scale) AutoZero() error {
	_, _ = s.Weigh()

	adjustScale := s.Scale
	s.Scale = 1
	s.Zero = 0

	time.Sleep(time.Second * 2)

	w, err := s.Weigh()
	if err != nil {
		return err
	}

	s.Scale = adjustScale
	s.Zero = w

	return nil
}

func (s *scale) Await(ctx context.Context, weight float64) <-chan struct{} {
	var done chan struct{}

	go func() {
		ticker := time.NewTicker(readTimeout)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if !s.buf.full {
					continue
				}

				w, _ := s.Weigh()
				if w >= weight {
					done <- struct{}{}
					ticker.Stop()
					return
				}
			}
		}
	}()

	return done
}

func (s *scale) Weigh() (float64, error) {
	for !s.buf.full {
		time.Sleep(time.Millisecond * 100)
	}

	data := s.buf.Data()

	sort.Ints(data)

	return (float64(data[len(data)/2]) - s.Zero) / s.Scale, nil
}

func (s *scale) Close() {
	s.hx.Halt()
	s.cancel()
}
