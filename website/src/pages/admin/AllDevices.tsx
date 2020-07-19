import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Button from '@material-ui/core/Button';
import { observer } from 'mobx-react';
import { grpc } from '../../Api';
import { lastSeen, lazy } from '../../Util';
import { Device } from '../../sdk/devices_pb';
import { confirm } from '../../components/Present';

@observer
export class AllDevices extends React.Component {
  devices = lazy(async () => {
    const res = await grpc.devices.listAllDevices({});
    return res.items;
  });

  deleteDevice = async (device: Device.AsObject) => {
    if (await confirm('Are you sure?')) {
      await grpc.devices.deleteDevice({
        name: device.name,
        owner: { value: device.owner },
      });
      await this.devices.refresh();
    }
  };

  render() {
    if (!this.devices.current) {
      return <p>loading...</p>;
    }

    const rows = this.devices.current;

    // show the provider column
    // when there is more than 1 provider in use
    // i.e. not all devices are from the same auth provider.
    const showProviderCol = rows.length >= 2 && rows.some((r) => r.ownerProvider !== rows[0].ownerProvider);

    return (
      <TableContainer>
        <Table stickyHeader>
          <TableHead>
            <TableRow>
              <TableCell>Owner</TableCell>
              {showProviderCol && <TableCell>Auth Provider</TableCell>}
              <TableCell>Device</TableCell>
              <TableCell>Connected</TableCell>
              <TableCell>Last Seen</TableCell>
              <TableCell>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {rows.map((row, i) => (
              <TableRow key={i}>
                <TableCell component="th" scope="row">
                  {row.ownerName || row.ownerEmail || row.owner}
                </TableCell>
                {showProviderCol && <TableCell>{row.ownerProvider}</TableCell>}
                <TableCell>{row.name}</TableCell>
                <TableCell>{row.connected ? 'yes' : 'no'}</TableCell>
                <TableCell>{lastSeen(row.lastHandshakeTime)}</TableCell>
                <TableCell>
                  <Button variant="outlined" color="secondary" onClick={() => this.deleteDevice(row)}>
                    Delete
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    );
  }
}
