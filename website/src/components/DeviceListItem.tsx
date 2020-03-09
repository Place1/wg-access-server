import React from 'react';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardContent from '@material-ui/core/CardContent';
import Typography from '@material-ui/core/Typography';
import Avatar from '@material-ui/core/Avatar';
import DonutSmallIcon from '@material-ui/icons/DonutSmall';
import MenuItem from '@material-ui/core/MenuItem';
import formatDistanceToNow from 'date-fns/formatDistanceToNow';
import { view } from 'react-easy-state';
import { AppState } from '../Store';
import { IconMenu } from './IconMenu';
import { PopoverDisplay } from './PopoverDisplay';
import { Device } from '../sdk/devices_pb'
import { grpc, toDate } from '../Api';

interface Props {
  device: Device.AsObject;
}

class DeviceListItem extends React.Component<Props> {
  dateString(date: Date) {
    if (date.getUTCMilliseconds() === 0) {
      return 'never';
    }
    return formatDistanceToNow(date, { addSuffix: true });
  }

  removeDevice = async () => {
    try {
      await grpc.devices.deleteDevice({
        name: this.props.device.name,
      });
      AppState.devices = AppState.devices.filter(device => device.name !== this.props.device.name);
    } catch {
      window.alert('api request failed');
    }
  };

  render() {
    const device = this.props.device;
    return (
      <Card>
        <CardHeader
          title={device.name}
          avatar={
            <Avatar>
              <DonutSmallIcon />
            </Avatar>
          }
          action={
            <IconMenu>
              <MenuItem style={{ color: 'red' }} onClick={this.removeDevice}>
                Delete
              </MenuItem>
            </IconMenu>
          }
        />
        <CardContent>
          <Typography component="p">
            Connected: {device.connected ? 'yes' : 'no'}
          </Typography>
          <Typography component="p">
            Endpoint: {device.endpoint}
          </Typography>
          <Typography component="p">
            Sent: {device.transmitBytes} bytes
          </Typography>
          <Typography component="p">
            Received: {device.receiveBytes} bytes
          </Typography>
          <Typography component="p">
            Last Handshake Time: {toDate(device.lastHandshakeTime!).toLocaleString()}
          </Typography>
          <Typography component="p">
            Public Key: <PopoverDisplay label="show">{device.publicKey}</PopoverDisplay>
          </Typography>
        </CardContent>
      </Card>
    );
  }
}

export default view(DeviceListItem);
