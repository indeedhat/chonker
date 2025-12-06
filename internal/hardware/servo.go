package hardware

import (
	"context"
	"fmt"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/host/v3/rpi"
)

const (
	frequency = 50 * physic.Hertz
	periodMs  = 20

	minPeriodMs = 0.5
	maxPeriodMs = 3.4

	angleStep = 5
	interval  = 100 * time.Millisecond
)

type servo struct {
	pin     gpio.PinIO
	angle   float64
	reverse bool
}

func newServo() servo {
	return servo{pin: rpi.P1_12}
}

func (s servo) Run(ctx context.Context) error {
	if err := s.pin.Out(gpio.Low); err != nil {
		fmt.Errorf("Failed to set pin to low: %s", err.Error())
	}

	go func() {
		ticker := time.NewTicker(time.Millisecond * 100)
		defer ticker.Stop()
		defer s.Close()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.step()

				duty := angleToDuty(s.angle)
				if err := s.pin.PWM(duty, frequency); err != nil {
					return
				}
			}
		}
	}()

	return nil
}

func (s servo) step() {
	if s.reverse {
		s.angle -= angleStep
	} else {
		s.angle += angleStep
	}

	if s.angle < 0 {
		s.angle = 0
		s.reverse = false
	}

	if s.angle > 180 {
		s.angle = 180
		s.reverse = true
	}
}

func (s servo) Close() error {
	return s.pin.Halt()
}

func angleToDuty(angle float64) gpio.Duty {
	if angle < 0 {
		angle = 0
	}

	if angle > 180 {
		angle = 180
	}

	pulseMs := minPeriodMs + (angle/180.0)*(maxPeriodMs-minPeriodMs)
	dutyFraction := pulseMs / periodMs
	return gpio.Duty(float64(gpio.DutyMax) * dutyFraction)
}
