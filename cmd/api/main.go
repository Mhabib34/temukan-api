package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"temukan-api/cmd/wire"
	"time"
)

func main() {
	// ── 1. Wire — inisialisasi semua dependency ───────────────────────────────
	engine, cleanup, err := wire.InitializeApp()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	// ── 2. Port dari env, default 8080 ────────────────────────────────────────
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// ── 3. HTTP server ────────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("[Server] listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[Server] listen error: %v", err)
		}
	}()

	// ── 4. Graceful shutdown ──────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[Server] shutting down...")

	// Stop worker (cancel context worker di dalam cleanup)
	cleanup()

	// Beri waktu 10 detik untuk request yang sedang berjalan selesai
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("[Server] forced shutdown: %v", err)
	}

	log.Println("[Server] exited cleanly")
}
