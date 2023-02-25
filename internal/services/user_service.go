package services

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/freifunkMUC/wg-access-server/internal/devices"
	"github.com/freifunkMUC/wg-access-server/pkg/authnz/authsession"
	"github.com/freifunkMUC/wg-access-server/proto/proto"
)

type UserService struct {
	proto.UnimplementedUsersServer
	DeviceManager *devices.DeviceManager
}

func (d *UserService) ListUsers(ctx context.Context, req *proto.ListUsersReq) (*proto.ListUsersRes, error) {
	user, err := authsession.CurrentUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "not authenticated")
	}

	if !user.Claims.Has("admin", "true") {
		return nil, status.Errorf(codes.PermissionDenied, "must be an admin")
	}

	users, err := d.DeviceManager.ListUsers()
	if err != nil {
		ctxlogrus.Extract(ctx).Error(err)
		return nil, status.Errorf(codes.Internal, "failed to retrieve users")
	}

	return &proto.ListUsersRes{
		Items: mapUsers(users),
	}, nil
}

func (d *UserService) DeleteUser(ctx context.Context, req *proto.DeleteUserReq) (*emptypb.Empty, error) {
	user, err := authsession.CurrentUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "not authenticated")
	}

	if !user.Claims.Has("admin", "true") {
		return nil, status.Errorf(codes.PermissionDenied, "must be an admin")
	}

	if err := d.DeviceManager.DeleteDevicesForUser(req.Name); err != nil {
		ctxlogrus.Extract(ctx).Error(err)
		return nil, status.Errorf(codes.Internal, "failed to delete user")
	}

	return &emptypb.Empty{}, nil
}

func mapUser(u *devices.User) *proto.User {
	return &proto.User{
		Name: u.Name,
		DisplayName: u.DisplayName,
	}
}

func mapUsers(users []*devices.User) []*proto.User {
	items := []*proto.User{}
	for _, u := range users {
		items = append(items, mapUser(u))
	}
	return items
}
