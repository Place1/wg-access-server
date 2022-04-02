package services

import (
	"context"

	"github.com/freifunkMUC/wg-access-server/internal/devices"
	"github.com/freifunkMUC/wg-access-server/internal/storage"
	"github.com/freifunkMUC/wg-access-server/pkg/authnz/authsession"
	"github.com/freifunkMUC/wg-access-server/proto/proto"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DeviceService struct {
	proto.UnimplementedDevicesServer
	DeviceManager *devices.DeviceManager
}

func (d *DeviceService) AddDevice(ctx context.Context, req *proto.AddDeviceReq) (*proto.Device, error) {
	user, err := authsession.CurrentUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "not authenticated")
	}

	device, err := d.DeviceManager.AddDevice(user, req.GetName(), req.GetPublicKey())
	if err != nil {
		ctxlogrus.Extract(ctx).Error(err)
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
		ctxlogrus.Extract(ctx).Error(err)
		return nil, status.Errorf(codes.Internal, "failed to retrieve devices")
	}
	return &proto.ListDevicesRes{
		Items: mapDevices(devices),
	}, nil
}

func (d *DeviceService) DeleteDevice(ctx context.Context, req *proto.DeleteDeviceReq) (*emptypb.Empty, error) {
	user, err := authsession.CurrentUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "not authenticated")
	}

	deviceOwner := user.Subject

	if req.Owner != nil {
		if user.Claims.Contains("admin") {
			deviceOwner = req.Owner.Value
		} else {
			return nil, status.Errorf(codes.PermissionDenied, "must be an admin")
		}
	}

	if err := d.DeviceManager.DeleteDevice(deviceOwner, req.GetName()); err != nil {
		ctxlogrus.Extract(ctx).Error(err)
		return nil, status.Errorf(codes.Internal, "failed to delete device")
	}

	return &emptypb.Empty{}, nil
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
		ctxlogrus.Extract(ctx).Error(err)
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
		Connected: d.LastHandshakeTime != nil && devices.IsConnected(*d.LastHandshakeTime),
	}
}

func mapDevices(devices []*storage.Device) []*proto.Device {
	items := []*proto.Device{}
	for _, d := range devices {
		items = append(items, mapDevice(d))
	}
	return items
}
