package infra

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ronannnn/infra/cfg"
	"go.uber.org/zap"
)

type HttpServerBaseRunner interface {
	RegisterRoutes() http.Handler
	Addr(httpAddr string, httpPort int) string
}

type HttpServerRunner struct {
	HttpServerBaseRunner
}

func (hs HttpServerRunner) Addr(httpAddr string, httpPort int) string {
	return fmt.Sprintf("%s:%d", httpAddr, httpPort)
}

type BaseHttpServer struct {
	HttpServerRunner
	Sys *cfg.Sys
	Log *zap.SugaredLogger
}

func (hs *BaseHttpServer) Run() {
	addr := hs.Addr(hs.Sys.HttpAddr, hs.Sys.HttpPort)
	server := &http.Server{Addr: addr, Handler: hs.RegisterRoutes()}
	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancelFn := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				hs.Log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			hs.Log.Fatal(err)
		}
		cancelFn()
		serverStopCtx()
	}()

	// Run the server
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		hs.Log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
