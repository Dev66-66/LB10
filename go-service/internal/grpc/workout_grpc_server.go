package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Dev66-66/LB10/go-service/internal/grpc/pb"
	"github.com/Dev66-66/LB10/go-service/internal/models"
	"github.com/Dev66-66/LB10/go-service/internal/store"
)

// WorkoutGRPCServer implements pb.WorkoutServiceServer using the shared in-memory store.
type WorkoutGRPCServer struct {
	pb.UnimplementedWorkoutServiceServer
	store *store.WorkoutStore
}

func NewWorkoutGRPCServer(s *store.WorkoutStore) *WorkoutGRPCServer {
	return &WorkoutGRPCServer{store: s}
}

func (s *WorkoutGRPCServer) GetWorkout(_ context.Context, req *pb.GetWorkoutRequest) (*pb.WorkoutResponse, error) {
	w, err := s.store.GetByID(int(req.Id))
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "workout with id %d not found", req.Id)
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &pb.WorkoutResponse{Workout: workoutToProto(w)}, nil
}

func (s *WorkoutGRPCServer) ListWorkouts(_ context.Context, _ *pb.ListWorkoutsRequest) (*pb.ListWorkoutsResponse, error) {
	all := s.store.GetAll()
	protos := make([]*pb.Workout, 0, len(all))
	for _, w := range all {
		protos = append(protos, workoutToProto(w))
	}
	return &pb.ListWorkoutsResponse{Workouts: protos}, nil
}

func (s *WorkoutGRPCServer) CreateWorkout(_ context.Context, req *pb.CreateWorkoutRequest) (*pb.WorkoutResponse, error) {
	w, err := s.store.Create(models.Workout{
		Name:           req.Name,
		Type:           models.WorkoutType(req.Type),
		Duration:       int(req.Duration),
		Difficulty:     models.WorkoutDifficulty(req.Difficulty),
		CaloriesBurned: int(req.CaloriesBurned),
	})
	if err != nil {
		if errors.Is(err, store.ErrInvalidName) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &pb.WorkoutResponse{Workout: workoutToProto(w)}, nil
}

func workoutToProto(w models.Workout) *pb.Workout {
	return &pb.Workout{
		Id:             int32(w.ID),
		Name:           w.Name,
		Type:           string(w.Type),
		Duration:       int32(w.Duration),
		Difficulty:     string(w.Difficulty),
		CaloriesBurned: int32(w.CaloriesBurned),
		CreatedAt:      w.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
