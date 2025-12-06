package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
