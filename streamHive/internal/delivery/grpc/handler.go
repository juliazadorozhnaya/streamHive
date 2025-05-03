package grpc

import (
	"context"
	"streamHive/streamHive/internal/delivery/grpc/pb"
	"streamHive/streamHive/internal/domain"
)

type Server struct {
	pb.UnimplementedStreamServiceServer
	uc domain.UseCase
}

func NewHandler(uc domain.UseCase) *Server {
	return &Server{uc: uc}
}

func (s *Server) Start(ctx context.Context, req *pb.StartRequest) (*pb.StartResponse, error) {
	jobID, err := s.uc.StartRecording(ctx, req.StreamUrl)
	if err != nil {
		return nil, err
	}
	return &pb.StartResponse{JobId: jobID}, nil
}

func (s *Server) Stop(_ context.Context, req *pb.StopRequest) (*pb.StopResponse, error) {
	err := s.uc.StopRecording(req.JobId)
	if err != nil {
		return nil, err
	}
	return &pb.StopResponse{Status: pb.JobStatus_STOPPED}, nil
}

func (s *Server) GetStatus(_ context.Context, req *pb.StatusRequest) (*pb.StatusResponse, error) {
	status, err := s.uc.GetJobStatus(req.JobId)
	if err != nil {
		return nil, err
	}

	var enum pb.JobStatus
	switch status {
	case "started":
		enum = pb.JobStatus_STARTED
	case "finished":
		enum = pb.JobStatus_FINISHED
	case "stopped":
		enum = pb.JobStatus_STOPPED
	case "failed":
		enum = pb.JobStatus_FAILED
	default:
		enum = pb.JobStatus_UNKNOWN
	}

	return &pb.StatusResponse{Status: enum}, nil
}
