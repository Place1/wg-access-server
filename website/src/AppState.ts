import { observable, makeObservable } from 'mobx';
import { InfoRes } from './sdk/server_pb';

class GlobalAppState {
  info?: InfoRes.AsObject;

  constructor() {
    makeObservable(this, {
      info: observable
    });
  }
}

export const AppState = new GlobalAppState();

console.info('see global app state by typing "window.AppState"');

Object.assign(window as any, {
  get AppState() {
    return JSON.parse(JSON.stringify(AppState));
  },
});
