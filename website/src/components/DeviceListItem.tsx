import React from 'react';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardContent from '@material-ui/core/CardContent';
import Avatar from '@material-ui/core/Avatar';
import WifiIcon from '@material-ui/icons/Wifi';
import WifiOffIcon from '@material-ui/icons/WifiOff';
import MenuItem from '@material-ui/core/MenuItem';
import formatDistanceToNow from 'date-fns/formatDistanceToNow';
import numeral from 'numeral';
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

  lastSeen() {
    if (this.props.device.lastHandshakeTime === undefined) {
      return 'Never';
    }
    return formatDistanceToNow(toDate(this.props.device.lastHandshakeTime!), {
      includeSeconds: true,
      addSuffix: true,
    });
  }

  render() {
    const device = this.props.device;
    return (
      <Card>
        <CardHeader
          title={device.name}
          avatar={
            <Avatar style={{ backgroundColor: device.connected ? '#76de8a' : '#bdbdbd' }}>
              {/* <DonutSmallIcon /> */}
              {device.connected
                ? <WifiIcon />
                : <WifiOffIcon />
              }
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
          <table cellPadding="5">
            <tbody>
              {AppState.info?.metadataEnabled && device.connected &&
                <>
                  <tr>
                    <td>Endpoint</td>
                    <td>{device.endpoint}</td>
                  </tr>
                  <tr>
                    <td>Sent</td>
                    <td>{numeral(device.transmitBytes).format('0b')}</td>
                  </tr>
                  <tr>
                    <td>Received</td>
                    <td>{numeral(device.receiveBytes).format('0b')}</td>
                  </tr>
                </>
              }
              {AppState.info?.metadataEnabled && !device.connected &&
                <>
                  <tr>
                    <td>Disconnected</td>
                  </tr>
                  <tr>
                    <td>Last Seen</td>
                    <td>{this.lastSeen()}</td>
                  </tr>
                </>
              }
              <tr>
                <td>Public key</td>
                <td><PopoverDisplay label="show">{device.publicKey}</PopoverDisplay></td>
              </tr>
            </tbody>
          </table>
        </CardContent>
      </Card>
    );
  }
}

export default view(DeviceListItem);
