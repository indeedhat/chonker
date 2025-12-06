package hardware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/indeedhat/chonker/internal/types"
	"github.com/indeedhat/dotenv"
	"periph.io/x/host/v3"
)

func init() {
	if _, err := host.Init(); err != nil {
		log.Fatalf("failed to init periph.io %s", err)
	}
}

const (
	leeway = 0.5
)

var (
	dispenseTimeout dotenv.Int = "DISPENSE_TIMEOUT"
)

func Disponse(weight float64) Report {
	ctx, cancel := context.WithTimeout(time.Second * time.Duration(dispenseTimeout.Get()))

	svo := newServo()
	if err := svo.Run(ctx); err != nil {
		return Report{err: err}
	}
	defer svo.Close()

	scl := newScale()
	defer scl.Close()

	select {
	case <-ctx.Done():
		return Report{err: errors.New("dispense timeout reached")}
	case <-scl.Await(weight - leeway):
		cancel()

		time.Sleep(time.Second)

		w, err := scl.Weigh()
		if err != nil {
			return Report{
				title: "Food dispensed",
				message: fmt.Sprint("dispensed UNKNOWN")
			}
		}

		return Report{
			title: "Food dispensed",
			message: fmt.Sprintf("dispensed %0.2fg", w)
		}
	}
}

type Report struct {
	err     error
	title   string
	message string
}

// Error implements types.Report.
func (r Report) Error() error {
	return r.err
}

// Message implements types.Report.
func (r Report) Message() string {
	return r.message
}

// Title implements types.Report.
func (r Report) Title() string {
	return r.title
}

var _ types.Report = (*Report)(nil)
