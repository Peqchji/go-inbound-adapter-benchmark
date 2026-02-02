package main

import (
	"log"
	"net/http"

	adapterinmemory "github.com/Peqchji/go-inbound-adapter-benchmark/internal/adapter/inmemory"
	adapterrest "github.com/Peqchji/go-inbound-adapter-benchmark/internal/adapter/rest"
	"github.com/Peqchji/go-inbound-adapter-benchmark/internal/client/database/inmemory"
	"github.com/Peqchji/go-inbound-adapter-benchmark/internal/domain/wallet"
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

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	log.Println("Starting HTTP server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
