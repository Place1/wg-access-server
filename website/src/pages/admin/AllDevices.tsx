import Button from '@mui/material/Button';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Typography from '@mui/material/Typography';
import WifiIcon from '@mui/icons-material/Wifi';
import WifiOffIcon from '@mui/icons-material/WifiOff';
import Avatar from "@mui/material/Avatar";
import { observer } from 'mobx-react';
import React from 'react';
import { grpc } from '../../Api';
import { AppState } from '../../AppState';
import { confirm } from '../../components/Present';
import { Device } from '../../sdk/devices_pb';
import { User } from '../../sdk/users_pb';
import { lastSeen, lazy } from '../../Util';
import numeral from "numeral";

export const AllDevices = observer(class AllDevices extends React.Component {
  users = lazy(async () => {
    const res = await grpc.users.listUsers({});
    return res.items;
  });

  devices = lazy(async () => {
    const res = await grpc.devices.listAllDevices({});
    let deviceList = res.items;
    deviceList.sort((d1, d2) => (d2.lastHandshakeTime ? d2.lastHandshakeTime.seconds : 0) - (d1.lastHandshakeTime ? d1.lastHandshakeTime.seconds : 0));
    return deviceList;
  });

  deleteUser = async (user: User.AsObject) => {
    if (await confirm('Are you sure you want to delete ' + user.name + '?')) {
      await grpc.users.deleteUser({
        name: user.name,
      });
      await this.users.refresh();
      await this.devices.refresh();
    }
  }

  deleteDevice = async (device: Device.AsObject) => {
    if (await confirm('Are you sure you want to delete ' + device.name + ' from ' + device.ownerName + '?')) {
      await grpc.devices.deleteDevice({
        name: device.name,
        owner: { value: device.owner },
      });
      await this.devices.refresh();
    }
  };

  render() {
    if (!this.devices.current || !this.users.current) {
      return <p>loading...</p>;
    }

    const users = this.users.current;
    const devices = this.devices.current;

    // show the provider column
    // when there is more than 1 provider in use
    // i.e. not all devices are from the same auth provider.
    const showProviderCol = devices.length >= 2 && devices.some((d) => d.ownerProvider !== devices[0].ownerProvider);

    return (
      <div style={{ display: 'grid', gridGap: 25, gridAutoFlow: 'row' }}>

        <Typography variant="h5" component="h5">
          Devices
            <Typography component="span"> ({devices.filter(p => p.connected).length}  of {devices.length} online)</Typography>
        </Typography>
        <TableContainer>
          <Table stickyHeader>
            <TableHead>
              <TableRow>
                <TableCell></TableCell>
                <TableCell>Owner</TableCell>
                {showProviderCol && <TableCell>Auth Provider</TableCell>}
                <TableCell>Device</TableCell>
                <TableCell>Connected</TableCell>
                <TableCell>Local Address</TableCell>
                <TableCell>Last Endpoint</TableCell>
                <TableCell>Download / Upload</TableCell>
                <TableCell>Last Seen</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {devices.map((device, i) => (
                <TableRow key={i}>
                  <TableCell>
                    <Avatar style={{ backgroundColor: device.connected ? '#76de8a' : '#bdbdbd' }}>
                      {/* <DonutSmallIcon /> */}
                      {device.connected ? <WifiIcon /> : <WifiOffIcon />}
                    </Avatar>
                  </TableCell>
                  <TableCell component="th" scope="row">
                    {device.ownerName || device.ownerEmail || device.owner}
                  </TableCell>
                  {showProviderCol && <TableCell>{device.ownerProvider}</TableCell>}
                  <TableCell>{device.name}</TableCell>
                  <TableCell>{device.connected ? 'yes' : 'no'}</TableCell>
                  <TableCell>{device.address}</TableCell>
                  <TableCell>{device.endpoint}</TableCell>
                  <TableCell>{ numeral(device.transmitBytes).format('0b') } / { numeral(device.receiveBytes).format('0b') }</TableCell>
                  <TableCell>{lastSeen(device.lastHandshakeTime)}</TableCell>
                  <TableCell>
                    <Button variant="outlined" color="secondary" onClick={() => this.deleteDevice(device)}>
                      Delete
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>

        <Typography variant="h5" component="h5">
            Users
            <Typography component="span"> ({users.length})</Typography>
        </Typography>
        <TableContainer>
            <Table stickyHeader>
                <TableHead>
                    <TableRow>
                        <TableCell>Name</TableCell>
                        <TableCell>Actions</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {users.map((user, i) => (
                        <TableRow key={i}>
                            <TableCell component="th" scope="row">
                                {user.displayName || user.name}
                            </TableCell>
                            <TableCell>
                                <Button variant="outlined" color="secondary" onClick={() => this.deleteUser(user)}>
                                    Delete
                                </Button>
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </TableContainer>

        <Typography variant="h5" component="h5">
          Server Info
        </Typography>
        <code>
          <pre>
            {JSON.stringify(AppState.info, null, 2)}

          </pre>
        </code>

      </div>
    );
  }
});
