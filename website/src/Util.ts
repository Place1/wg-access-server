import formatDistance from "date-fns/formatDistance";
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { toDate } from "./Api";

export function sleep(seconds: number) {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve();
    }, seconds * 1000);
  });
}

export function lastSeen(timestamp: Timestamp.AsObject | undefined): string {
  if (timestamp === undefined) {
    return 'Never';
  }
  return formatDistance(toDate(timestamp), new Date(), {
    addSuffix: true,
  });
}
