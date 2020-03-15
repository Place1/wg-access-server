import { observable } from 'mobx';
import { InfoRes } from './sdk/server_pb';

class GlobalAppState {

  @observable
  info?: InfoRes.AsObject;

}

export const AppState = new GlobalAppState();

console.info('see global app state by typing "window.AppState"');

Object.assign(window as any, {
  get AppState() {
    return JSON.parse(JSON.stringify(AppState));
  }
});
