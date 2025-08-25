package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/jstaud/edgelink-go/internal/api"
	"github.com/jstaud/edgelink-go/internal/config"
	"github.com/jstaud/edgelink-go/internal/device/fake"
	"github.com/jstaud/edgelink-go/internal/log"
	"github.com/jstaud/edgelink-go/internal/observability"
	"github.com/jstaud/edgelink-go/internal/pipeline"
	"github.com/jstaud/edgelink-go/internal/pipeline/publisher"
)

func main() {
	// Command-line flag for config file path
	var cfgPath string
	
	// Create root command using Cobra CLI framework
	root := &cobra.Command{
		Use: "edgelink",
		Short: "EdgeLink SCADA gateway",
		RunE: func(_ *cobra.Command, _ []string) error {
			// Create logger
			l := log.New()
			
			// Load configuration
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return err
			}
			
			// Initialize Prometheus metrics
			observability.Init()
			
			// Create in-memory cache for latest readings
			cache := pipeline.NewMemoryCache()
			
			// Create MQTT publisher (or console for testing)
			// pub := publisher.NewMQTT(cfg.Broker.MQTTURL, cfg.Broker.ClientID, cfg.Broker.BaseTopic)
			pub := publisher.NewConsole(cfg.Broker.BaseTopic)  // Use console for testing
			defer pub.Close()  // Clean up connection when main exits
			
			// Create context that cancels on SIGINT/SIGTERM (Ctrl+C or kill signal)
			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()
			
			// Create and start HTTP server in a goroutine
			srv := &http.Server{
				Addr:    cfg.HTTPAddr,
				Handler: api.NewREST(cache),
			}
			go func() {
				l.Info("HTTP server starting", zap.String("addr", cfg.HTTPAddr))
				if err := srv.ListenAndServe(); err != http.ErrServerClosed {
					l.Error("HTTP server error", zap.Error(err))
				}
			}()
			
			// Start one poller goroutine for each configured device
			for _, spec := range cfg.Devices {
				spec := spec  // Capture loop variable for goroutine
				drv := fake.New(spec)  // Create fake device driver
				
				go func() {
					err := pipeline.RunPoller(ctx, l, drv, spec, pub, cache)
					if err != nil && err != context.Canceled {
						l.Error("Poller error", zap.String("device", spec.ID), zap.Error(err))
					}
				}()
			}
			
			// Wait for shutdown signal
			<-ctx.Done()
			l.Info("Shutdown signal received")
			
			// Give things time to clean up gracefully
			time.Sleep(cfg.ShutdownWait)
			
			// Shutdown HTTP server gracefully
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = srv.Shutdown(shutdownCtx)
			
			l.Info("Shutdown complete")
			return nil
		},
	}
	
	// Add command-line flag for config file
	root.Flags().StringVarP(&cfgPath, "config", "c", "configs/example.yaml", "config file path")
	
	// Execute the command
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
