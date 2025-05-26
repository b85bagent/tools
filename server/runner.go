package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Runner struct {
	servers []Server
}

func NewRunner(servers ...Server) *Runner {
	return &Runner{servers: servers}
}

func (r *Runner) Run() {
	var wg sync.WaitGroup

	for _, s := range r.servers {
		wg.Add(1)
		go func(s Server) {
			defer wg.Done()
			if err := s.Start(); err != nil {
				log.Printf("server exited with error: %v", err)
			}
		}(s)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Printf("received signal: %s, shutting down...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, s := range r.servers {
		if err := s.Stop(ctx); err != nil {
			log.Printf("error stopping server: %v", err)
		}
	}

	wg.Wait()
	log.Println("all servers stopped gracefully")
}
