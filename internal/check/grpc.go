package check

import (
	"context"
	"errors"

	"github.com/shashimalcse/cronuseo/internal/util"
	cronuseo "github.com/shashimalcse/cronuseo/proto"
	"google.golang.org/grpc/metadata"
)

func NewGrpcService(service Service) cronuseo.CheckServer {

	return grpcService{service: service}
}

type grpcService struct {
	service Service
}

func (s grpcService) Check(ctx context.Context, req *cronuseo.GrpcCheckRequest) (*cronuseo.GrpcCheckResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata from request")
	}
	apiKey := md.Get("API_KEY")[0]

	input := CheckRequest{
		Username: req.Username,
		Action:   req.Action,
		Resource: req.Resource,
	}

	allow, err := s.service.Check(context.Background(), req.Organization, input, apiKey)
	if err != nil {
		return nil, util.HandleError(err)
	}

	return &cronuseo.GrpcCheckResponse{Allow: allow}, nil
}
