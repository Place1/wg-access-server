import { store } from 'react-easy-state';
import { Device } from './sdk/devices_pb';

export const AppState = store({
  devices: new Array<Device.AsObject>(),
});
