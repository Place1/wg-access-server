import React from 'react';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardContent from '@material-ui/core/CardContent';
import Avatar from '@material-ui/core/Avatar';
import WifiIcon from '@material-ui/icons/Wifi';
import WifiOffIcon from '@material-ui/icons/WifiOff';
import MenuItem from '@material-ui/core/MenuItem';
import numeral from 'numeral';
import { lastSeen } from '../Util';
import { AppState } from '../AppState';
import { IconMenu } from './IconMenu';
import { PopoverDisplay } from './PopoverDisplay';
import { Device } from '../sdk/devices_pb'
import { grpc } from '../Api';
import { observer } from 'mobx-react';

interface Props {
  device: Device.AsObject;
  onRemove: () => void;
}

@observer
export class DeviceListItem extends React.Component<Props> {
  removeDevice = async () => {
    try {
      await grpc.devices.deleteDevice({
        name: this.props.device.name,
      });
      this.props.onRemove();
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
                    <td>Download</td>
                    <td>{numeral(device.transmitBytes).format('0b')}</td>
                  </tr>
                  <tr>
                    <td>Upload</td>
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
                    <td>{lastSeen(device.lastHandshakeTime)}</td>
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
