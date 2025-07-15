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
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Server struct {
	taskpb.UnimplementedTaskServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func ToProto(t *queue.IntTask, queueType taskpb.QueueType) *taskpb.IntTask {
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
		QueueType: queueType,
	}
}

func (s *Server) GetTask(ctx context.Context, req *taskpb.GetTaskRequest) (*taskpb.IntTask, error) {
	var task queue.IntTask
	var ok bool
	var queueType taskpb.QueueType

	switch req.WorkerType {
	case taskpb.WorkerType_WORKER_TYPE_PYTHON:
		task, ok = queue.PythonLocalQueue.PQ.Pop()
		queueType = taskpb.QueueType_QUEUE_TYPE_PYTHON
	case taskpb.WorkerType_WORKER_TYPE_GO:
		task, ok = queue.GoLocalQueue.PQ.Pop()
		queueType = taskpb.QueueType_QUEUE_TYPE_GO
	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid worker type: %v", req.WorkerType)
	}

	if !ok {
		return nil, status.Error(codes.NotFound, "queue empty")
	}

	return ToProto(&task, queueType), nil
}

func (s *Server) GetPythonTask(ctx context.Context, e *emptypb.Empty) (*taskpb.IntTask, error) {
	task, ok := queue.PythonLocalQueue.PQ.Pop()
	if !ok {
		return nil, status.Error(codes.NotFound, "Python queue empty")
	}
	logging.DebugLog(fmt.Sprintf("Dispatching task: %s", task.ID))
	return ToProto(&task, taskpb.QueueType_QUEUE_TYPE_PYTHON), nil
}

func (s *Server) GetGoTask(ctx context.Context, e *emptypb.Empty) (*taskpb.IntTask, error) {
	task, ok := queue.GoLocalQueue.PQ.Pop()
	if !ok {
		return nil, status.Error(codes.NotFound, "Go queue empty")
	}
	return ToProto(&task, taskpb.QueueType_QUEUE_TYPE_GO), nil
}

func Start() error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return fmt.Errorf("gRPC listen failed: %w", err)
	}

	gs := grpc.NewServer()
	taskpb.RegisterTaskServiceServer(gs, NewServer())

	logging.DebugLog("gRPC server started on :50051")

	if err := gs.Serve(lis); err != nil {
		return fmt.Errorf("gRPC serve failed: %w", err)
	}
	return nil
}
