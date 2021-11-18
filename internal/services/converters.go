package services

import (
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TimestampToTime(value *timestamp.Timestamp) time.Time {
	return time.Unix(value.Seconds, int64(value.Nanos))
}

func TimeToTimestamp(value *time.Time) *timestamp.Timestamp {
	if value == nil {
		return nil
	}
	t := timestamppb.New(*value)
	if t == nil {
		logrus.Error("bad time value")
		t = timestamppb.Now()
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
