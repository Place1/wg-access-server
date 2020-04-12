package services

import (
	"context"
	"time"

	"github.com/place1/wg-access-server/pkg/authnz/authsession"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/place1/wg-access-server/internal/devices"
	"github.com/place1/wg-access-server/internal/storage"
	"github.com/place1/wg-access-server/proto/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeviceService struct {
	DeviceManager *devices.DeviceManager
}

func (d *DeviceService) AddDevice(ctx context.Context, req *proto.AddDeviceReq) (*proto.Device, error) {
	user, err := authsession.CurrentUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "not authenticated")
	}

	device, err := d.DeviceManager.AddDevice(user, req.GetName(), req.GetPublicKey())
	if err != nil {
		logrus.Error(err)
		return nil, status.Errorf(codes.Internal, "failed to add device")
	}

	return mapDevice(device), nil
}

func (d *DeviceService) ListDevices(ctx context.Context, req *proto.ListDevicesReq) (*proto.ListDevicesRes, error) {
	user, err := authsession.CurrentUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "not authenticated")
	}

	devices, err := d.DeviceManager.ListDevices(user.Subject)
	if err != nil {
		logrus.Error(err)
		return nil, status.Errorf(codes.Internal, "failed to retrieve devices")
	}
	return &proto.ListDevicesRes{
		Items: mapDevices(devices),
	}, nil
}

func (d *DeviceService) DeleteDevice(ctx context.Context, req *proto.DeleteDeviceReq) (*empty.Empty, error) {
	user, err := authsession.CurrentUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "not authenticated")
	}

	if err := d.DeviceManager.DeleteDevice(user.Subject, req.GetName()); err != nil {
		logrus.Error(err)
		return nil, status.Errorf(codes.Internal, "failed to delete device")
	}

	return &empty.Empty{}, nil
}

func (d *DeviceService) ListAllDevices(ctx context.Context, req *proto.ListAllDevicesReq) (*proto.ListAllDevicesRes, error) {
	user, err := authsession.CurrentUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "not authenticated")
	}

	if !user.Claims.Contains("admin") {
		return nil, status.Errorf(codes.PermissionDenied, "must be an admin")
	}

	devices, err := d.DeviceManager.ListAllDevices()
	if err != nil {
		logrus.Error(err)
		return nil, status.Errorf(codes.Internal, "failed to retrieve devices")
	}

	return &proto.ListAllDevicesRes{
		Items: mapDevices(devices),
	}, nil
}

func mapDevice(d *storage.Device) *proto.Device {
	return &proto.Device{
		Name:              d.Name,
		Owner:             d.Owner,
		OwnerName:         d.OwnerName,
		OwnerEmail:        d.OwnerEmail,
		OwnerProvider:     d.OwnerProvider,
		PublicKey:         d.PublicKey,
		Address:           d.Address,
		CreatedAt:         TimeToTimestamp(&d.CreatedAt),
		LastHandshakeTime: TimeToTimestamp(d.LastHandshakeTime),
		ReceiveBytes:      d.ReceiveBytes,
		TransmitBytes:     d.TransmitBytes,
		Endpoint:          d.Endpoint,
		/**
		 * Wireguard is a connectionless UDP protocol - data is only
		 * sent over the wire when the client is sending real traffic.
		 * Wireguard has no keep alive packets by default to remain as
		 * silent as possible.
		 *
		 */
		Connected: isConnected(d.LastHandshakeTime),
	}
}

func mapDevices(devices []*storage.Device) []*proto.Device {
	items := []*proto.Device{}
	for _, d := range devices {
		items = append(items, mapDevice(d))
	}
	return items
}

func isConnected(lastHandshake *time.Time) bool {
	if lastHandshake == nil {
		return false
	}
	return lastHandshake.After(time.Now().Add(-1 * time.Minute))
}
