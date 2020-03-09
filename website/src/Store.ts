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

console.info('see global app state by typing "window.AppState"');
Object.assign(window as any, {
  get AppState() {
    return JSON.parse(JSON.stringify(AppState));
  }
});
