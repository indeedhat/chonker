package hardware

import (
	"context"

	"periph.io/x/devices/v3/hx711"
	"periph.io/x/host/v3/rpi"
)

type scale struct {
	hx *hx711.Dev
}

func newScale() (*scale, error) {
	hx, err := hx711.New(
		rpi.P1_5,
		rpi.P1_3,
	)
	if err != nil {
		return nil, err
	}

	return &scale{hx}, nil
}

func (s *scale) Await(ctx context.Context, weight float64) <-chan struct{} {
	reading := s.hx.ReadContinuous()
	var done chan struct{}

	go func() {
		defer s.Close()
		for {
			select {
			case <-ctx.Done():
				return
			case r := <-reading:
				// TODO: get the proper weight
				w := float64(r.Raw)
				if w >= weight {
					done <- struct{}{}
					return
				}
			}
		}
	}()

	return done
}

func (s *scale) Weigh() (float64, error) {
	r, err := s.hx.Read()
	if err != nil {
		return 0, err
	}

	// TODO: find out what i need to do to convert this to g
	return float64(r.Raw), nil
}

func (s *scale) Close() {
	s.hx.Halt()
}
