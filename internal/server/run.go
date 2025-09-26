package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"basic-gin/internal/cache"
	rediscache "basic-gin/internal/cache/redis"
	"basic-gin/internal/config"
	"basic-gin/internal/db"
	"basic-gin/internal/handler"
	"basic-gin/internal/repository"
	"basic-gin/internal/service"
)

func Run(ctx context.Context) error {
	config.Load()

	pool, err := db.Connect(ctx, config.App.PostgresDSN)

	if err != nil {
		return fmt.Errorf("db connect: %w", err)
	}

	defer pool.Close()

	var c cache.Cache
	if addr := config.App.RedisAddr; addr != "" {
		rc := rediscache.New(addr, config.App.RedisPass, 0)

		pctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		if err := rc.Ping(pctx); err != nil {
			log.Println("redis disabled: %w", err)
		} else {
			c = rc

			defer rc.Close()
			log.Println("reddis connected")
		}
	}

	client_repo := repository.NewClientRepository(pool)
	client_service := service.NewClientService(*client_repo, c)

	account_repo := repository.NewAccountRepository(pool)
	account_service := service.NewAccountService(account_repo, client_service, c)

	transaction_repo := repository.NewTransactionRepository(pool)
	transaction_service := service.NewTransactionService(transaction_repo, account_repo)

	deps := handler.NewDependencies(client_service, account_service, transaction_service)

	router := newRouter(deps)

	srv := &http.Server{
		Addr:              ":" + config.App.ServerPort,
		Handler:           router,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)

	go func() {
		log.Println("🚀 Server listening on", srv.Addr)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Println("shutdown signal recieved")

		shCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_ = srv.Shutdown(shCtx)
		return nil
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	}
}
