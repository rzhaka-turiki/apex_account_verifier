package grpc

import (
	"context"

	"github.com/rzhaka-turiki/apex_account_verifier/internal/service"
	"github.com/rzhaka-turiki/apex_account_verifier/proto/apexpb"
)

type Server struct {
	apexpb.UnimplementedApexVerifierServer
	verifier *service.Verifier
}

func NewServer(verifier *service.Verifier) *Server {
	return &Server{verifier: verifier}
}

func (s *Server) VerifyAccount(ctx context.Context, req apexpb.VerifyAccountRequest) (*apexpb.VerifyAccountResponse, error) {
	account, err := s.verifier.VerifyAccount(ctx, req.GetPlayer(), req.GetPlatform(), int(req.GetLevel()))
	if err != nil {
		return nil, err
	}
	return &apexpb.VerifyAccountResponse{
		Uid:      account.UID,
		Player:   account.Player,
		Platform: account.Platform,
		Level:    int32(account.Level),
	}, nil
}
