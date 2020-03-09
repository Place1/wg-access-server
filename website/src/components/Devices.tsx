import React from 'react';
import Grid from '@material-ui/core/Grid';
import { view } from 'react-easy-state';
import { AppState } from '../Store';
import DeviceListItem from './DeviceListItem';
import { grpc } from '../Api';
import { sleep } from '../Util';

class Devices extends React.Component {

  private mounted = false;

  componentDidMount() {
    this.mounted = true;
    this.poll();
  }

  componentWillUnmount() {
    this.mounted = false;
  }

  async poll() {
    while (this.mounted) {
      // we sleep first because we pre-load the list of
      // devices when the app starts up (index.tsx)
      await sleep(5);
      const res = await grpc.devices.listDevices({});
      AppState.devices = res.items;
    }
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
