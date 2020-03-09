package services

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/sirupsen/logrus"
)

func TimestampToTime(value *timestamp.Timestamp) time.Time {
	return time.Unix(value.Seconds, int64(value.Nanos))
}

func TimeToTimestamp(value *time.Time) *timestamp.Timestamp {
	if value == nil {
		return nil
	}
	t, err := ptypes.TimestampProto(*value)
	if err != nil {
		logrus.Error("bad time value")
		t = ptypes.TimestampNow()
	}
	return t
}

func stringValue(value *string) *wrappers.StringValue {
	if value != nil {
		return &wrappers.StringValue{
			Value: *value,
		}
	}
	return nil
}
