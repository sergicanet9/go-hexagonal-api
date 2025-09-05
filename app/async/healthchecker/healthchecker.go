package healthchecker

import (
	"context"
	"net/http"
	"time"

	"github.com/sergicanet9/go-hexagonal-api/proto/gen/go/pb"
	"github.com/sergicanet9/scv-go-tools/v3/observability"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const contentType = "application/json"

func RunHTTP(ctx context.Context, cancel context.CancelFunc, url string, interval time.Duration) {
	defer cancel()
	defer func() {
		if rec := recover(); rec != nil {
			observability.Logger().Printf("FATAL - recovered panic in HTTP healthchecker process: %v", rec)
		}
	}()

	for ctx.Err() == nil {
		<-time.After(interval)

		func() {
			req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
			if err != nil {
				observability.Logger().Printf("HTTP healthchecker process - error: %s", err)
				return
			}
			req.Header.Set("Content-Type", contentType)

			start := time.Now()
			resp, err := http.DefaultClient.Do(req)
			elapsed := time.Since(start)

			if err != nil {
				observability.Logger().Printf("HTTP healthchecker process - error: %s", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				observability.Logger().Printf("HTTP healthchecker process - error: %s", err)
				return
			}

			observability.Logger().Printf("HTTP healthchecker process - health Check complete, time elapsed: %s", elapsed)
		}()
	}
}

func RunGRPC(ctx context.Context, cancel context.CancelFunc, target string, interval time.Duration) {
	defer cancel()
	defer func() {
		if rec := recover(); rec != nil {
			observability.Logger().Printf("FATAL - recovered panic in gRPC healthchecker process: %v", rec)
		}
	}()

	for ctx.Err() == nil {
		<-time.After(interval)

		func() {
			conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				observability.Logger().Printf("gRPC healthchecker process - error: %v", err)
				return
			}
			defer conn.Close()
			client := pb.NewHealthServiceClient(conn)

			start := time.Now()
			_, err = client.HealthCheck(ctx, nil)
			elapsed := time.Since(start)

			if err != nil {
				observability.Logger().Printf("gRPC healthchecker process - error: %v", err)
				return
			}

			if _, ok := status.FromError(err); !ok {
				observability.Logger().Printf("gRPC healthchecker process - error: %v", err)
			}

			observability.Logger().Printf("gRPC healthchecker process - health Check complete, time elapsed: %s", elapsed)
		}()
	}
}
