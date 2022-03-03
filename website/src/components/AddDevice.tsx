import Button from '@material-ui/core/Button';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import CardHeader from '@material-ui/core/CardHeader';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormControl from '@material-ui/core/FormControl';
import FormHelperText from '@material-ui/core/FormHelperText';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import Typography from '@material-ui/core/Typography';
import AddIcon from '@material-ui/icons/Add';
import { codeBlock } from 'common-tags';
import { observable } from 'mobx';
import { observer } from 'mobx-react';
import React from 'react';
import { box_keyPair } from 'tweetnacl-ts';
import { grpc } from '../Api';
import { AppState } from '../AppState';
import { GetConnected } from './GetConnected';
import { Info } from './Info';

interface Props {
  onAdd: () => void;
}

@observer
export class AddDevice extends React.Component<Props> {
  @observable
  dialogOpen = false;

  @observable
  error?: string;

  @observable
  deviceName = '';

  @observable
  configFile?: string;

  submit = async (event: React.FormEvent) => {
    event.preventDefault();

    const keypair = box_keyPair();
    const publicKey = window.btoa(String.fromCharCode(...(new Uint8Array(keypair.publicKey) as any)));
    const privateKey = window.btoa(String.fromCharCode(...(new Uint8Array(keypair.secretKey) as any)));

    try {
      const device = await grpc.devices.addDevice({
        name: this.deviceName,
        publicKey,
      });
      this.props.onAdd();

      const info = AppState.info!;
      const configFile = codeBlock`
        [Interface]
        PrivateKey = ${privateKey}
        Address = ${device.address}
        ${info.dnsEnabled && `DNS = ${info.dnsAddress}`}
        MTU = ${info.mtu}

        [Peer]
        PublicKey = ${info.publicKey}
        AllowedIPs = ${info.allowedIps}
        Endpoint = ${`${info.host?.value || window.location.hostname}:${info.port || '51820'}`}
      `;

      this.configFile = configFile;
      this.dialogOpen = true;
      this.reset();
    } catch (error) {
      console.log(error);
      // TODO: unwrap grpc error message
      this.error = 'failed';
    }
  };

  reset = () => {
    this.deviceName = '';
  };

  render() {
    return (
      <>
        <Card>
          <CardHeader title="Add A Device" />
          <CardContent>
            <form onSubmit={this.submit}>
              <FormControl error={!!this.error} fullWidth>
                <InputLabel htmlFor="device-name">Device Name</InputLabel>
                <Input
                  id="device-name"
                  value={this.deviceName}
                  onChange={(event) => (this.deviceName = event.currentTarget.value)}
                  aria-describedby="device-name-text"
                />
                <FormHelperText id="device-name-text">{this.error}</FormHelperText>
              </FormControl>
              <Typography component="div" align="right">
                <Button color="secondary" type="button" onClick={this.reset}>
                  Cancel
                </Button>
                <Button color="primary" variant="contained" endIcon={<AddIcon />} type="submit">
                  Add
                </Button>
              </Typography>
            </form>
          </CardContent>
        </Card>
        <Dialog disableBackdropClick disableEscapeKeyDown maxWidth="xl" open={this.dialogOpen}>
          <DialogTitle>
            Get Connected
            <Info>
              <Typography component="p" style={{ paddingBottom: 8 }}>
                Your VPN connection file is not stored by this portal.
              </Typography>
              <Typography component="p" style={{ paddingBottom: 8 }}>
                If you lose this file you can simply create a new device on this portal to generate a new connection
                file.
              </Typography>
              <Typography component="p">
                The connection file contains your WireGuard Private Key (i.e. password) and should{' '}
                <strong>never</strong> be shared.
              </Typography>
            </Info>
          </DialogTitle>
          <DialogContent>
            <GetConnected configFile={this.configFile!} />
          </DialogContent>
          <DialogActions>
            <Button color="secondary" variant="outlined" onClick={() => (this.dialogOpen = false)}>
              Done
            </Button>
          </DialogActions>
        </Dialog>
      </>
    );
  }
}
