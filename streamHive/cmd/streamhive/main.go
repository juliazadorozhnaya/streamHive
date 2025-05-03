package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"streamHive/streamHive/internal/delivery/grpc"
	httpdelivery "streamHive/streamHive/internal/delivery/http"
	"streamHive/streamHive/internal/infrastructure/config"
	"streamHive/streamHive/internal/infrastructure/pubsub"
	"streamHive/streamHive/internal/infrastructure/rtsp"
	"streamHive/streamHive/internal/infrastructure/state"
	"streamHive/streamHive/internal/infrastructure/storage"
	"streamHive/streamHive/internal/usecase"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	recorder := rtsp.New()
	store := storage.New(cfg.StoragePath, cfg.Retention.String())
	store.StartCleanupWorker()
	pub := pubsub.NewNATS(cfg.NATSURL)
	stateStore := state.NewMongo(cfg.MongoURI, cfg.MongoDB, cfg.MongoCollection)

	uc := usecase.New(recorder, store, pub, stateStore)
	uc.StartBackgroundConsumer(ctx)

	handler := httpdelivery.NewHandler(uc)
	httpServer := &http.Server{
		Addr:    cfg.HTTPPort,
		Handler: httpdelivery.NewRouter(handler),
	}

	go func() {
		log.Printf("HTTP server listening on %s", cfg.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("error while serving HTTP: %v", err)
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	grpcSrv := grpc.NewHandler(uc)
	grpcServer, listener := grpc.ServeWithInstance(cfg.GRPCPort, grpcSrv)
	go func() {
		log.Printf("gRPC server listening on %s", cfg.GRPCPort)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP shutdown error: %v", err)
	} else {
		log.Println("HTTP server stopped gracefully")
	}

	grpcServer.GracefulStop()
	log.Println("gRPC server stopped gracefully")
}
