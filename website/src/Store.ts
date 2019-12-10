import { store } from 'react-easy-state';

export interface IDevice {
  name: string;
  publicKey: string;
  endpoint: string;
  address: string;
  dns: string;
  createdAt: string;
  serverPublicKey: string;
  // TODO: these fields on backend
  // receiveBytes: number;
  // transmitBytes: number;
  // lastHandshakeTime: string;
}

export const AppState = store({
  devices: new Array<IDevice>(),
});
