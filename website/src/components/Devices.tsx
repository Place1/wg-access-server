import React from 'react';
import Grid from '@material-ui/core/Grid';
import { view } from 'react-easy-state';
import { AppState } from '../Store';
import DeviceListItem from './DeviceListItem';
import { grpc } from '../Api';

class Devices extends React.Component {

  componentDidMount() {
    this.load();
  }

  async load() {
    const res = await grpc.devices.listDevices({});
    AppState.devices = res.items;
  }

  render() {
    return (
      <Grid container spacing={3}>
        {AppState.devices.map((device, i) => (
          <Grid key={i} item xs={12} sm={6} md={4} lg={3}>
            <DeviceListItem device={device} />
          </Grid>
        ))}
      </Grid>
    );
  }
}

export default view(Devices);
