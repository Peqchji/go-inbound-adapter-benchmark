package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/Peqchji/go-inbound-adapter-benchmark/cmd/grpcserver/proto"
	adapterinmemory "github.com/Peqchji/go-inbound-adapter-benchmark/internal/adapter/inmemory"
	"github.com/Peqchji/go-inbound-adapter-benchmark/internal/client/database/inmemory"
	"github.com/Peqchji/go-inbound-adapter-benchmark/internal/domain/wallet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	proto.UnimplementedWalletServiceServer
	service *wallet.WalletService
}

func (s *server) GetWallet(ctx context.Context, req *proto.GetWalletRequest) (*proto.GetWalletResponse, error) {
	result := s.service.GetWallet(req.Id)
	if result.Err != nil {
		return nil, result.Err
	}

	w := result.Res
	return &proto.GetWalletResponse{
		Wallet: &proto.Wallet{
			Id:      w.ID(),
			Balance: w.Balance(),
			Owner: &proto.Owner{
				Id:        w.Owner().ID(),
				Firstname: w.Owner().Firstname(),
				Lastname:  w.Owner().Lastname(),
			},
		},
	}, nil
}

func (s *server) CreateWallet(ctx context.Context, req *proto.CreateWalletRequest) (*proto.CreateWalletResponse, error) {
	result := s.service.CreateWallet(req.OwnerId, req.Firstname, req.Lastname)
	if result.Err != nil {
		return nil, result.Err
	}

	w := result.Res
	return &proto.CreateWalletResponse{
		Wallet: &proto.Wallet{
			Id:      w.ID(),
			Balance: w.Balance(),
			Owner: &proto.Owner{
				Id:        w.Owner().ID(),
				Firstname: w.Owner().Firstname(),
				Lastname:  w.Owner().Lastname(),
			},
		},
	}, nil
}

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Initialize dependencies
	dbClient := inmemory.NewInMemoryClient()
	_ = dbClient.CreateTable("wallet")
	walletTable, err := dbClient.GetTable("wallet")
	if err != nil {
		log.Fatalf("failed to get wallet table: %v", err)
	}
	repo := adapterinmemory.NewInMemoryWalletAdapter(walletTable)
	svc := wallet.NewWalletService(repo)

	s := grpc.NewServer()
	proto.RegisterWalletServiceServer(s, &server{service: svc})
	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
