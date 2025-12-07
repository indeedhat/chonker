package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/indeedhat/chonker/internal/hardware"
	"github.com/indeedhat/chonker/internal/schedule"
	"github.com/indeedhat/chonker/internal/types"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sched, err := schedule.Load()
	if err != nil {
		log.Fatalf("failed to load schedule: %s", err)
	}

	runner := schedule.NewRunner(
		sched,
		func(f float64) types.Report {
			return nil
		},
		func(r types.Report) {},
	)

	runner.Start(ctx)

	svr := &http.Server{
		Addr:    ":8087",
		Handler: http.DefaultServeMux,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("ListenAndServer: %v", svr.ListenAndServe())
	}()

	<-quit
	log.Print("Shutting down server...")
	cancel()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := svr.Shutdown(shutdownCtx); err != nil {
		log.Print("Server forced to shutdown after timeout")
	}
}

func calibrateScale(weight1, weight2 float64) {
	var (
		adjustZero float64
		reading1   float64
		reading2   float64
	)

	scale, err := hardware.NewScale()
	if err != nil {
		log.Fatalf("failed to init scale: %s", err)
	}

	// NB: the adjustment values is pulled from the env so we need to make sure to reset it back to 1
	//     to get accurate calibration values
	scale.Scale = 1

	fmt.Println("Make sure scale is working and empty, getting weight in 5 seconds...")
	time.Sleep(5 * time.Second)

	fmt.Println("Getting weight...")
	adjustZero, err = scale.Weigh()
	if err != nil {
		fmt.Println("ReadDataMedianRaw error:", err)
		return
	}

	fmt.Println("Raw weight is:", adjustZero)
	fmt.Println("")

	fmt.Printf("Put first weight of %.2f on scale, getting weight in 15 seconds...\n", weight1)
	time.Sleep(15 * time.Second)

	fmt.Println("Getting weight...")
	reading1, err = scale.Weigh()
	if err != nil {
		fmt.Println("ReadDataMedianRaw error:", err)
		return
	}

	fmt.Println("Raw weight is:", reading1)
	fmt.Println("")

	fmt.Printf("Put second weight of %.2f on scale, getting weight in 15 seconds...\n", weight2)
	time.Sleep(15 * time.Second)

	fmt.Println("Getting weight...")
	reading2, err = scale.Weigh()
	if err != nil {
		fmt.Println("ReadDataMedianRaw error:", err)
		return
	}
	fmt.Println("Raw weight is ", reading2)
	fmt.Println("")

	adjust1 := float64(reading1-adjustZero) / weight1
	adjust2 := float64(reading2-adjustZero) / weight2

	fmt.Println("AdjustZero should be set to:", adjustZero)
	fmt.Printf("AdjustScale should be set to a value between %f and %f\n", adjust1, adjust2)
	fmt.Println("")
}
