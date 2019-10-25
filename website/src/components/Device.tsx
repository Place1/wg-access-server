import React from 'react';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardContent from '@material-ui/core/CardContent';
import Typography from '@material-ui/core/Typography';
import Avatar from '@material-ui/core/Avatar';
import DonutSmallIcon from '@material-ui/icons/DonutSmall';
import MenuItem from '@material-ui/core/MenuItem';
import formatDistanceToNow from 'date-fns/formatDistanceToNow'
import { view } from 'react-easy-state';
import { IDevice, AppState } from '../Store';
import { IconMenu } from './IconMenu';
import { PopoverDisplay } from './PopoverDisplay';

interface Props {
  device: IDevice;
}

class Device extends React.Component<Props> {

  dateString(date: Date) {
    if (date.getUTCMilliseconds() === 0) {
      return 'never';
    }
    return formatDistanceToNow(date, { addSuffix: true });
  }

  removeDevice = async () => {
    const res = await fetch(`/api/devices/${this.props.device.name}`, {
      method: 'DELETE',
    });
    if (res.status === 204) {
      AppState.devices = AppState.devices.filter(device => device.name !== this.props.device.name);
    } else {
      window.alert(await res.text());
    }
  }

  render() {
    const device = this.props.device;
    return (
      <Card>
        <CardHeader
          title={device.name}
          avatar={<Avatar><DonutSmallIcon /></Avatar>}
          action={
            <IconMenu>
              <MenuItem style={{ color: 'red' }} onClick={this.removeDevice}>Delete</MenuItem>
            </IconMenu>
          }
        />
        <CardContent>
          <Typography component="p">
            Public Key: <PopoverDisplay label="show">{device.publicKey}</PopoverDisplay>
          </Typography>
          <Typography component="p">
            Endpoint: {device.endpoint}
          </Typography>
        </CardContent>
      </Card>
    );
  }
}

export default view(Device);
