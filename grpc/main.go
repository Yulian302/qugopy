package grpc

import (
	"context"
	"fmt"
	"net"

	taskpb "github.com/Yulian302/qugopy/github.com/Yulian302/qugopy/proto"
	"github.com/Yulian302/qugopy/internal/queue"
	"github.com/Yulian302/qugopy/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Server struct {
	taskpb.UnimplementedTaskServiceServer
}

func ToProto(t *queue.IntTask) *taskpb.IntTask {
	if t == nil {
		return nil
	}
	var deadline *timestamppb.Timestamp

	if t.Task.Deadline != nil {
		deadline = timestamppb.New(*t.Task.Deadline)
	}

	var recurring *wrapperspb.BoolValue
	if t.Task.Recurring != nil {
		recurring = wrapperspb.Bool(*t.Task.Recurring)
	}

	return &taskpb.IntTask{
		Id: t.ID,
		Task: &taskpb.Task{
			Type:      t.Task.Type,
			Payload:   t.Task.Payload,
			Priority:  uint32(t.Task.Priority),
			Deadline:  deadline,
			Recurring: recurring,
		},
	}
}

func (s *Server) GetTask(ctx context.Context, _ *taskpb.Empty) (*taskpb.IntTask, error) {
	t, ok := queue.DefaultLocalQueue.PQ.Pop()
	if !ok {
		return nil, status.Error(codes.NotFound, "queue empty")
	}
	return ToProto(&t), nil
}

func Start() error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return fmt.Errorf("gRPC listen failed: %w", err)
	}
	gs := grpc.NewServer()
	taskpb.RegisterTaskServiceServer(gs, &Server{})

	logging.DebugLog("gRPC listening on :50051")

	if err := gs.Serve(lis); err != nil {
		return fmt.Errorf("gRPC serve failed: %w", err)
	}
	return nil
}
