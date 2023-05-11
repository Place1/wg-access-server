import Button from '@material-ui/core/Button';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import CardHeader from '@material-ui/core/CardHeader';
import Checkbox from '@material-ui/core/Checkbox';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormHelperText from '@material-ui/core/FormHelperText';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import Typography from '@material-ui/core/Typography';
import AddIcon from '@material-ui/icons/Add';
import { codeBlock } from 'common-tags';
import { makeObservable, observable } from 'mobx';
import { observer } from 'mobx-react';
import React from 'react';
import { box_keyPair, randomBytes } from 'tweetnacl-ts';
import { grpc } from '../Api';
import { AppState } from '../AppState';
import { GetConnected } from './GetConnected';
import { Info } from './Info';

import Accordion from '@mui/material/Accordion';
import AccordionSummary from '@mui/material/AccordionSummary';
import AccordionDetails from '@mui/material/AccordionDetails';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import Box from '@material-ui/core/Box';

interface Props {
  onAdd: () => void;
}

export const AddDevice = observer(class AddDevice extends React.Component<Props> {
  dialogOpen = false;

  error?: string;

  deviceName = '';

  devicePublickey = '';

  useDevicePresharekey = false;

  showAdvancedOptions = false;

  configFile?: string;

  showMobile = true;

  submit = async (event: React.FormEvent) => {
    event.preventDefault();

    const keypair = box_keyPair();
    var publicKey: string;
    var privateKey: string;
    if (this.devicePublickey) {
      publicKey = this.devicePublickey
      privateKey = 'pleaseReplaceThisPrivatekey'
      this.showMobile = false;
    } else {
      publicKey = window.btoa(String.fromCharCode(...(new Uint8Array(keypair.publicKey) as any)));
      privateKey = window.btoa(String.fromCharCode(...(new Uint8Array(keypair.secretKey) as any)));
      this.showMobile = true;
    }

    const presharedKey = this.useDevicePresharekey ? window.btoa(String.fromCharCode(...(randomBytes(32) as any))) : '';

    try {
      const device = await grpc.devices.addDevice({
        name: this.deviceName,
        publicKey,
        presharedKey,
      });
      this.props.onAdd();

      const info = AppState.info!;

      const dnsInfo = [];
      if (info.clientConfigDnsServers) {
        // If custom DNS entries are specified via client config, prefer them over the calculated ones.
        dnsInfo.push(info.clientConfigDnsServers);
      }
      else if (info.dnsEnabled) {
        // Otherwise, and if DNS is enabled, use the ones from the server.
        dnsInfo.push(info.dnsAddress);
      }

      if (info.clientConfigDnsSearchDomain) {
        // In any case, if there is a custom search domain configured in the client config, append it to the list of DNS servers.
        dnsInfo.push(info.clientConfigDnsSearchDomain)
      }

      const configFile = codeBlock`
        [Interface]
        PrivateKey = ${privateKey}
        Address = ${device.address}
        ${ 0 < dnsInfo.length && `DNS = ${ dnsInfo.join(", ") }` }

        [Peer]
        PublicKey = ${info.publicKey}
        AllowedIPs = ${info.allowedIps}
        Endpoint = ${`${info.host?.value || window.location.hostname}:${info.port || '51820'}`}
        ${ this.useDevicePresharekey ? `PresharedKey = ${presharedKey}` : `` }
      `;

      this.configFile = configFile;
      this.dialogOpen = true;
      this.reset();
    } catch (error) {
      console.log(error);
      // TODO: unwrap grpc error message
      this.error = 'failed to add device';
    }
  };

  reset = () => {
    this.deviceName = '';
    this.devicePublickey = '';
    this.useDevicePresharekey = false;
    this.showAdvancedOptions = false;
    this.error = '';
  };

  constructor(props: Props) {
    super(props);

    makeObservable(this, {
      dialogOpen: observable,
      error: observable,
      deviceName: observable,
      devicePublickey: observable,
      useDevicePresharekey: observable,
      configFile: observable,
      showMobile: observable
    });
  }

  render() {
    return (
      <>
        <Card>
          <CardHeader title="Add A Device" />
          <CardContent>
            <form onSubmit={this.submit}>
              <FormControl fullWidth>
                <InputLabel htmlFor="device-name">Device Name</InputLabel>
                <Input
                  id="device-name"
                  value={this.deviceName}
                  onChange={(event) => (this.deviceName = event.currentTarget.value)}
                  aria-describedby="device-name-text"
                />
              </FormControl>
              <Box mt={2} mb={2}>
                <Accordion>
                  <AccordionSummary
                    expandIcon={<ExpandMoreIcon />}
                    aria-controls="advanced-options-content"
                    id="advanced-options-header"
                  >
                    <Typography>Advanced</Typography>
                  </AccordionSummary>
                  <AccordionDetails>
                    <FormControl fullWidth>
                      <InputLabel htmlFor="device-publickey">Device Public Key (Optional)</InputLabel>
                      <Input
                        id="device-publickey"
                        value={this.devicePublickey}
                        onChange={(event) => (this.devicePublickey = event.currentTarget.value)}
                        aria-describedby="device-publickey-text"
                      />
                      <FormHelperText id="device-publickey-text">Put your public key to a pre-generated private key here. Replace the private key in the config file after downloading it.</FormHelperText>
                    </FormControl>
                    <FormControlLabel 
                      control={
                        <Checkbox
                          id="device-presharedkey"
                          value={this.useDevicePresharekey}
                          onChange={(event) => (this.useDevicePresharekey = event.currentTarget.checked)}
                        />
                      } 
                      label="Use pre-shared key" 
                    />
                  </AccordionDetails>
                </Accordion>
              </Box>
              <FormHelperText id="device-error-text" error={true}>{this.error}</FormHelperText>
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
            <GetConnected configFile={this.configFile!} showMobile={this.showMobile}/>
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
});
