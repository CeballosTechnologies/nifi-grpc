package nifi

import (
	context "context"
	"fmt"
	"log"
	"net"
	"testing"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestServer(t *testing.T) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 4000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	clientConn, err := grpc.Dial(fmt.Sprintf("localhost:%d", 4000), opts...)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	nifiClient := NewFlowFileServiceClient(clientConn)

	ts := &TestStruct{}
	ts.id = 1

	server := NewServer(ts.flowFileHandler)

	go server.Serve(lis)

	t.Run("TestFlowFileHandler", func(t *testing.T) {
		ctx := context.Background()

		flowFileRequest := &FlowFileRequest{Content: []byte("test")}

		reply, err := nifiClient.Send(ctx, flowFileRequest)
		if err != nil {
			t.Errorf(err.Error())
		}

		fmt.Println(reply.ResponseCode.String())
	})
}

func (ts *TestStruct) flowFileHandler(context.Context, *FlowFileRequest) (*FlowFileReply, error) {
	return &FlowFileReply{ResponseCode: FlowFileReply_SUCCESS}, nil
}

type TestStruct struct {
	id int
}
