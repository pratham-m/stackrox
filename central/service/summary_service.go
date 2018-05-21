package service

import (
	"bitbucket.org/stack-rox/apollo/central/datastore"
	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"bitbucket.org/stack-rox/apollo/pkg/grpc/authz/user"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewSummaryService returns the SummaryService object.
func NewSummaryService(alerts datastore.AlertDataStore, clusters datastore.ClusterDataStore, deployments datastore.DeploymentDataStore, images datastore.ImageDataStore) *SummaryService {
	return &SummaryService{
		alerts:      alerts,
		clusters:    clusters,
		deployments: deployments,
		images:      images,
	}
}

// SummaryService serves Summary APIs.
type SummaryService struct {
	alerts      datastore.AlertDataStore
	clusters    datastore.ClusterDataStore
	deployments datastore.DeploymentDataStore
	images      datastore.ImageDataStore
}

// RegisterServiceServer registers this service with the given gRPC Server.
func (s *SummaryService) RegisterServiceServer(grpcServer *grpc.Server) {
	v1.RegisterSummaryServiceServer(grpcServer, s)
}

// RegisterServiceHandlerFromEndpoint registers this service with the given gRPC Gateway endpoint.
func (s *SummaryService) RegisterServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return v1.RegisterSummaryServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}

// AuthFuncOverride specifies the auth criteria for this API.
func (s *SummaryService) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, returnErrorCode(user.Any().Authorized(ctx))
}

// GetSummaryCounts returns the global counts of alerts, clusters, deployments, and images.
func (s *SummaryService) GetSummaryCounts(context.Context, *empty.Empty) (*v1.SummaryCountsResponse, error) {
	alerts, err := s.alerts.CountAlerts()
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	clusters, err := s.clusters.CountClusters()
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	deployments, err := s.deployments.CountDeployments()
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	images, err := s.images.CountImages()
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &v1.SummaryCountsResponse{
		NumAlerts:      int64(alerts),
		NumClusters:    int64(clusters),
		NumDeployments: int64(deployments),
		NumImages:      int64(images),
	}, nil
}
