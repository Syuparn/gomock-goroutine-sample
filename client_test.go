package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/syuparn/gomock-goroutine-sample/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

func testClient(
	ctx context.Context,
	mockSrv *proto.MockPersonServer,
) (proto.PersonClient, error) {
	lis := bufconn.Listen(bufSize)

	// create client
	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer))
	if err != nil {
		return nil, fmt.Errorf("%+v", err)
	}

	s := grpc.NewServer()
	proto.RegisterPersonServer(s, mockSrv)

	// run server concurrently
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	return proto.NewPersonClient(conn), nil
}

func TestGetName(t *testing.T) {
	ctx := context.TODO()

	// create mock server
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSrv := proto.NewMockPersonServer(ctrl)
	mockSrv.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(&proto.GetResponse{
			Id:   1234,
			Name: "Taro",
		}, nil)

	// create client
	client, err := testClient(ctx, mockSrv)

	// test
	h := NewPersonHandler(client)
	actual, err := h.GetName(ctx, PersonID(1234))
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, PersonName("Taro"), actual)
}

// NOTE: this is a sample of failed test
// test never ends when unimplemented method is called because t.Fatalf does not work
func TestGetNameValidationErr(t *testing.T) {
	ctx := context.TODO()

	// create mock server
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSrv := proto.NewMockPersonServer(ctrl) // no methods are expected to be called

	// create client
	client, err := testClient(ctx, mockSrv)

	// test
	h := NewPersonHandler(client)
	_, err = h.GetName(ctx, PersonID(-1))
	if err == nil {
		t.Error("validation error must be occurred")
	}
}

// NOTE: this is a sample of failed test
// test immediately fails when unimplemented method is called
func TestGetNameValidationErrRevised(t *testing.T) {
	ctx := context.TODO()

	// create mock server
	ctrl := gomock.NewController(NewConcurrentTestReporter(t))
	defer ctrl.Finish()

	mockSrv := proto.NewMockPersonServer(ctrl) // no methods are expected to be called

	// create client
	client, err := testClient(ctx, mockSrv)

	// test
	h := NewPersonHandler(client)
	_, err = h.GetName(ctx, PersonID(-1))
	if err == nil {
		t.Error("validation error must be occurred")
	}
}

type ConcurrentTestReporter struct {
	*testing.T
}

func NewConcurrentTestReporter(t *testing.T) *ConcurrentTestReporter {
	return &ConcurrentTestReporter{t}
}

func (r *ConcurrentTestReporter) Fatalf(format string, args ...interface{}) {
	// use os.Exit(1) to kill all goroutines
	log.Fatalf(format, args...)
}
