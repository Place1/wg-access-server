import { store } from 'react-easy-state';
import { Device } from './sdk/devices_pb';
import { InfoRes } from './sdk/server_pb';

interface State {
  info?: InfoRes.AsObject,
  devices: Array<Device.AsObject>,
}

export const AppState = store<State>({
  devices: [],
});
