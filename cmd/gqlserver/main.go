package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Peqchji/go-inbound-adapter-benchmark/cmd/gqlserver/graph"
	adapterinmemory "github.com/Peqchji/go-inbound-adapter-benchmark/internal/adapter/inmemory"
	"github.com/Peqchji/go-inbound-adapter-benchmark/internal/client/database/inmemory"
	"github.com/Peqchji/go-inbound-adapter-benchmark/internal/domain/wallet"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbClient := inmemory.NewInMemoryClient()
	_ = dbClient.CreateTable("wallet")

	walletTable, err := dbClient.GetTable("wallet")
	if err != nil {
		log.Fatalf("failed to get wallet table: %v", err)
	}

	repo := adapterinmemory.NewInMemoryWalletAdapter(walletTable)
	svc := wallet.NewWalletService(repo)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		Service: svc,
	}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
