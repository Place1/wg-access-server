import React from 'react';
import { view } from 'react-easy-state';
import TableContainer from '@material-ui/core/TableContainer';
import Table from '@material-ui/core/Table';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import TableCell from '@material-ui/core/TableCell';
import TableBody from '@material-ui/core/TableBody';
import { Device } from '../../sdk/devices_pb';
import { dateToTimestamp } from '../../Api';
import { lastSeen } from '../../Util';

class AllDevices extends React.Component {
  render() {
    const rows: Array<Device.AsObject> = [
      {
        name: 'Example',
        owner: 'John',
        connected: false,
        lastHandshakeTime: dateToTimestamp(new Date()),
        endpoint: "192.160.0.100",
        address: "10.44.0.2",
        publicKey: "example",
        receiveBytes: 0,
        transmitBytes: 0,
        createdAt: dateToTimestamp(new Date()),
      },
    ];

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

export default view(AllDevices);
