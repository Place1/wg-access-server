import formatDistance from "date-fns/formatDistance";
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { toDate } from "./Api";
import { fromResource } from "mobx-utils";

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

export function updatingDatasource<T>(seconds: number, cb: () => Promise<T>) {
  let running = false;
  let sink: ((next: T) => void) | undefined;
  return {
    update: async () => {
      if (sink) {
        sink(await cb());
      }
    },
    ...fromResource<T>(
      async s => {
        sink = s;
        running = true;
        while (running) {
          sink(await cb());
          await sleep(seconds);
        }
      },
      () => {
        running = false;
      }
    )
  }
}
