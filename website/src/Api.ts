import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { Devices } from './sdk/devices_pb';
import { Server } from './sdk/server_pb';

const backend = window.location.origin + '/api';

export const grpc = {
  server: new Server(backend),
  devices: new Devices(backend),
}

// https://github.com/SafetyCulture/grpc-web-devtools
const devtools = (window as any).__GRPCWEB_DEVTOOLS__;
if (devtools) {
  devtools(Object.values(grpc));
}

// utils
export function toDate(timestamp: Timestamp.AsObject): Date {
  const t = new Timestamp();
  t.setSeconds(timestamp.seconds);
  t.setNanos(timestamp.nanos);
  return t.toDate();
}

export function dateToTimestamp(date: Date): Timestamp.AsObject {
  return {
    seconds: Math.round(date.getTime() / 1000),
    nanos: 0,
  };
}
