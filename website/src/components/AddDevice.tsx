import React from 'react';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardContent from '@material-ui/core/CardContent';
import AddIcon from '@material-ui/icons/Add';
import CardActions from '@material-ui/core/CardActions';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';
import DialogActions from '@material-ui/core/DialogActions';
import qrcode from 'qrcode';
import { view } from 'react-easy-state';
import { box_keyPair } from 'tweetnacl-ts';
import { codeBlock } from 'common-tags';
import { FormHelperText } from '@material-ui/core';
import { GetConnected } from './GetConnected';
import { IDevice, AppState } from '../Store';

class AddDevice extends React.Component {

  state = {
    name: '',
    open: false,
    qrCodeUri: '',
    configFileUri: '',
    error: '',
  };

  onAdd = async (event: React.FormEvent) => {
    event.preventDefault();
    const keypair = box_keyPair();
    const b64PublicKey = window.btoa(String.fromCharCode(...new Uint8Array(keypair.publicKey) as any));
    const b64PrivateKey = window.btoa(String.fromCharCode(...new Uint8Array(keypair.secretKey) as any));

    const res = await fetch('/api/devices', {
      method: 'POST',
      body: JSON.stringify({
        name: this.state.name,
        publicKey: b64PublicKey,
      }),
    });
    if (res.status >= 400) {
      this.setState({ error: await res.text() });
      return;
    }
    const { device } = await res.json() as { device: IDevice };

    AppState.devices.push(device);

    const configFile = codeBlock`
      [Interface]
      PrivateKey = ${b64PrivateKey}
      Address = ${device.address}
      DNS = ${'1.1.1.1, 8.8.8.8'}

      [Peer]
      PublicKey = ${device.serverPublicKey}
      AllowedIPs = 0.0.0.0/1, 128.0.0.0/1
      Endpoint = ${device.endpoint}
    `;

    this.setState({
      open: true,
      qrCodeUri: await qrcode.toDataURL(configFile),
      configFileUri: URL.createObjectURL(new Blob([configFile])),
    });
  }

  render() {
    return (
      <form onSubmit={this.onAdd}>
        <Card>
          <CardHeader
            title="Add a device"
          />
          <CardContent>
            <TextField
              label="Device Name"
              error={this.state.error !== ''}
              value={this.state.name}
              onChange={(event) => this.setState({ name: event.currentTarget.value })}
              style={{ marginTop: -20, marginBottom: 8 }}
              fullWidth
            />
            {this.state.error !== '' && <FormHelperText>{this.state.error}</FormHelperText>}
          </CardContent>
          <CardActions>
            <Button
              color="primary"
              variant="contained"
              endIcon={<AddIcon />}
              type="submit"
            >
              Add
            </Button>
            <Dialog
              disableBackdropClick
              disableEscapeKeyDown
              maxWidth="xl"
              open={this.state.open}
            >
              <DialogTitle>Get Connected</DialogTitle>
              <DialogContent>
                <GetConnected
                  qrCodeUri={this.state.qrCodeUri}
                  configFileUri={this.state.configFileUri}
                />
              </DialogContent>
              <DialogActions>
                <Button color="secondary" variant="outlined" onClick={() => this.setState({ open: false })}>
                  Done
                </Button>
              </DialogActions>
            </Dialog>
          </CardActions>
        </Card>
      </form>
    );
  }
}

export default view(AddDevice);
