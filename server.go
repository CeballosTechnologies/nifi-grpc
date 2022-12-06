package nifi

import (
	context "context"

	"google.golang.org/grpc"
)

type Server struct {
	UnimplementedFlowFileServiceServer
	flowFileHandler FlowFileHandler
}

type FlowFileHandler func(context.Context, *FlowFileRequest) (*FlowFileReply, error)

func NewServer(f FlowFileHandler) *grpc.Server {
	server := new(Server)
	server.flowFileHandler = f

	grpcServer := grpc.NewServer()
	RegisterFlowFileServiceServer(grpcServer, server)

	return grpcServer
}

// Implements Nifi gRPC service definition. Accepts flowfiles sent from Nifi.
// https://raw.githubusercontent.com/apache/nifi/main/nifi-nar-bundles/nifi-grpc-bundle/nifi-grpc-processors/src/main/resources/proto/flowfile_service.proto
func (s *Server) Send(ctx context.Context, req *FlowFileRequest) (reply *FlowFileReply, err error) {
	return s.flowFileHandler(ctx, req)
}
