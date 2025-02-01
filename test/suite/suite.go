package suite

import (
	"auth-api/internal/config"
	"context"
	auth "github.com/Iliuxa/protos/gen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient auth.AuthServiceClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoad()

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	grpcAddress := net.JoinHostPort("localhost", strconv.Itoa(cfg.GRPC.Port))

	cc, err := grpc.DialContext(context.Background(), grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal("grpc server connection failed: %v", err)
	}

	authClient := auth.NewAuthServiceClient(cc)

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: authClient,
	}
}
