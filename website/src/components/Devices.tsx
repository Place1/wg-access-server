import React from 'react';
import Grid from '@material-ui/core/Grid';
import { observable, makeObservable } from 'mobx';
import { observer } from 'mobx-react';
import { grpc } from '../Api';
import { autorefresh } from '../Util';
import { DeviceListItem } from './DeviceListItem';
import { AddDevice } from './AddDevice';

export const Devices = observer(class Devices extends React.Component {
  devices = autorefresh(30, async () => {
    return (await grpc.devices.listDevices({})).items;
  });

  constructor(props: {}) {
    super(props);

    makeObservable(this, {
      devices: observable
    });
  }

  componentWillUnmount() {
    this.devices.dispose();
  }

  render() {
    if (!this.devices.current) {
      return <p>loading...</p>;
    }
    return (
      <Grid container spacing={3} justify="center">
        <Grid item xs={12}>
          <Grid container spacing={3}>
            {this.devices.current.map((device, i) => (
              <Grid key={i} item xs={12} sm={6} md={4} lg={3}>
                <DeviceListItem device={device} onRemove={() => this.devices.refresh()} />
              </Grid>
            ))}
          </Grid>
        </Grid>
        <Grid item xs={12} sm={10} md={10} lg={6}>
          <AddDevice onAdd={() => this.devices.refresh()} />
        </Grid>
      </Grid>
    );
  }
});
