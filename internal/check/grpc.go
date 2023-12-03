package check

import (
	"context"
	"errors"

	"github.com/shashimalcse/cronuseo/internal/util"
	"github.com/shashimalcse/cronuseo/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

func NewGrpcService(service Service, logger *zap.Logger) proto.CheckServer {

	return grpcService{service: service, logger: logger}
}

type grpcService struct {
	service Service
	logger  *zap.Logger
}

func (s grpcService) Check(ctx context.Context, req *proto.GrpcCheckRequest) (*proto.GrpcCheckResponse, error) {

	s.logger.Info("GRPC method : Check", zap.String("method", "Check"))
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata from request")
	}
	apiKey := md.Get("API_KEY")[0]

	input := CheckRequest{
		Identifier: req.Username,
		Action:     req.Action,
		Resource:   req.Resource,
	}

	allow, err := s.service.Check(context.Background(), req.Organization, input, apiKey, false)
	if err != nil {
		return nil, util.HandleError(err)
	}

	return &proto.GrpcCheckResponse{Allow: allow.Allow}, nil
}
