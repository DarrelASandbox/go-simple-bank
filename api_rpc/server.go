package api_rpc

import (
	"fmt"

	db "github.com/DarrelASandbox/go-simple-bank/db/sqlc"
	"github.com/DarrelASandbox/go-simple-bank/util"
	"github.com/DarrelASandbox/go-simple-bank/pb"
	"github.com/DarrelASandbox/go-simple-bank/token"
)

// Server serves gRPC requests for our banking service
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer creates a new gRPC server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
