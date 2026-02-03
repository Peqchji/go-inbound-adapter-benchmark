package main

import (
	"log"
	"os"

	adapterinmemory "github.com/Peqchji/go-inbound-adapter-benchmark/internal/adapter/inmemory"
	adapterrest "github.com/Peqchji/go-inbound-adapter-benchmark/internal/adapter/rest"
	"github.com/Peqchji/go-inbound-adapter-benchmark/internal/client/database/inmemory"
	"github.com/Peqchji/go-inbound-adapter-benchmark/internal/domain/wallet"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	dbClient := inmemory.NewInMemoryClient()
	_ = dbClient.CreateTable("wallet")

	walletTable, err := dbClient.GetTable("wallet")
	if err != nil {
		log.Fatalf("failed to get wallet table: %v", err)
	}

	repo := adapterinmemory.NewInMemoryWalletAdapter(walletTable)
	svc := wallet.NewWalletService(repo)

	handler := adapterrest.NewWalletHandler(svc)

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	handler.RegisterRoutes(e)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
