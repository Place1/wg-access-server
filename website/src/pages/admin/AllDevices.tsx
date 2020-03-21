import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import { observer } from 'mobx-react';
import { lazyObservable } from 'mobx-utils';
import { grpc } from '../../Api';
import { Device } from '../../sdk/devices_pb';
import { lastSeen } from '../../Util';

@observer
export class AllDevices extends React.Component {

  devices = lazyObservable<Device.AsObject[]>(async sink => {
    const res = await grpc.devices.listAllDevices({});
    sink(res.items);
  });

  render() {
    if (!this.devices.current()) {
      return <p>loading...</p>
    }

    const rows = this.devices.current();

    return (
      <TableContainer>
        <Table stickyHeader>
          <TableHead>
            <TableRow>
              <TableCell>Owner</TableCell>
              <TableCell>Device</TableCell>
              <TableCell>Connected</TableCell>
              <TableCell>Last Seen</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {rows.map((row, i) => (
              <TableRow key={i}>
                <TableCell component="th" scope="row">
                  {row.owner}
                </TableCell>
                <TableCell>{row.name}</TableCell>
                <TableCell>{row.connected ? 'yes' : 'no'}</TableCell>
                <TableCell>{lastSeen(row.lastHandshakeTime)}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    );
  }
}
