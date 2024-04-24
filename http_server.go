package infra

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ronannnn/infra/cfg"
	"github.com/ronannnn/infra/utils"
	"go.uber.org/zap"
)

// golang abstract class reference: https://adrianwit.medium.com/abstract-class-reinvented-with-go-4a7326525034

type HttpServerBaseRunner interface {
	RegisterRoutes() *chi.Mux
	LogRegisteredRoutes(routes *chi.Mux, log *zap.SugaredLogger)
	Addr(httpAddr string, httpPort int) string
}

type HttpServerRunner struct {
	HttpServerBaseRunner
}

func (hs HttpServerRunner) Addr(httpAddr string, httpPort int) string {
	return fmt.Sprintf("%s:%d", httpAddr, httpPort)
}

func (hs HttpServerRunner) LogRegisteredRoutes(routes *chi.Mux, log *zap.SugaredLogger) {
	var routesInfo [][]string
	if err := chi.Walk(routes, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		routesInfo = append(routesInfo, []string{method, route})
		return nil
	}); err != nil {
		log.Errorf("Error while walking routes: %s", err.Error())
	}
	utils.LeftJustifyingPrint(routesInfo, log)
}

type BaseHttpServer struct {
	HttpServerRunner
	Sys *cfg.Sys
	Log *zap.SugaredLogger
}

func (hs *BaseHttpServer) Run() {
	addr := hs.Addr(hs.Sys.HttpAddr, hs.Sys.HttpPort)
	routes := hs.RegisterRoutes()
	if hs.Sys.PrintRoutes {
		hs.LogRegisteredRoutes(routes, hs.Log)
	}
	server := &http.Server{Addr: addr, Handler: routes}
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
	hs.Log.Info(fmt.Sprintf("Server running on %s", addr))
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		hs.Log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
