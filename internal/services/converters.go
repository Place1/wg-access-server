package services

import (
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TimestampToTime(value *timestamppb.Timestamp) time.Time {
	return time.Unix(value.Seconds, int64(value.Nanos))
}

func TimeToTimestamp(value *time.Time) *timestamppb.Timestamp {
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

func stringValue(value *string) *wrapperspb.StringValue {
	if value != nil {
		return &wrapperspb.StringValue{
			Value: *value,
		}
	}
	return nil
}
