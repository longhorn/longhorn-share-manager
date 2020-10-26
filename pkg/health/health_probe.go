package health

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/longhorn/longhorn-share-manager/pkg/server"
)

type CheckServer struct {
	manager *server.ShareManager
}

func NewHealthCheckServer(manager *server.ShareManager) *CheckServer {
	// TODO: implement a proper health checker that takes the nfs servers status in consideration
	return &CheckServer{manager: manager}
}

func (hc *CheckServer) Check(context.Context, *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	if hc.manager == nil {
		return &healthpb.HealthCheckResponse{
			Status: healthpb.HealthCheckResponse_NOT_SERVING,
		}, fmt.Errorf("share manager is not running")
	}

	return &healthpb.HealthCheckResponse{
		Status: healthpb.HealthCheckResponse_SERVING,
	}, nil
}

func (hc *CheckServer) Watch(req *healthpb.HealthCheckRequest, ws healthpb.Health_WatchServer) error {
	for ; ; time.Sleep(time.Second) {
		if hc.manager == nil {
			if err := ws.Send(&healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_NOT_SERVING}); err != nil {
				logrus.Errorf("failed to send share manager health check result %v: %v",
					healthpb.HealthCheckResponse_NOT_SERVING, err)
			}
			continue
		}

		if err := ws.Send(&healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}); err != nil {
			logrus.Errorf("failed to send share manager health check result %v: %v",
				healthpb.HealthCheckResponse_SERVING, err)
		}
	}
}
