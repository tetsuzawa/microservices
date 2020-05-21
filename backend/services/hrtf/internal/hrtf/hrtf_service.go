package hrtf

import (
	"context"

	"github.com/tetsuzawa/microservices/backend/pkg/api"
)

type HRTFServiceServer struct {
	r Repository
}

func NewHRTFServiceServer(r Repository) api.HRTFServiceServer {
	return &HRTFServiceServer{r}
}

func (s *HRTFServiceServer) CreateHRTF(ctx context.Context, request *api.CreateHRTFRequest) (*api.CreateHRTFResponse, error) {
	//TODO
	return &api.CreateHRTFResponse{}, nil
}

func (s *HRTFServiceServer) GetHRTF(ctx context.Context, request *api.GetHRTFRequest) (*api.GetHRTFResponse, error) {
	//TODO
	return &api.GetHRTFResponse{}, nil
}

func (s *HRTFServiceServer) UpdateHRTF(ctx context.Context, request *api.UpdateHRTFRequest) (*api.UpdateHRTFResponse, error) {
	//TODO
	return &api.UpdateHRTFResponse{}, nil
}

func (s *HRTFServiceServer) DeleteHRTF(ctx context.Context, request *api.DeleteHRTFRequest) (*api.DeleteHRTFResponse, error) {
	//TODO
	return &api.DeleteHRTFResponse{}, nil
}

func (s *HRTFServiceServer) ListHRTFs(ctx context.Context, request *api.ListHRTFsRequest) (*api.ListHRTFsResponse, error) {
	//TODO
	return &api.ListHRTFsResponse{}, nil
}
