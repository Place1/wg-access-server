import React from 'react';

import AddIcon from '@material-ui/icons/Add';
import Button from '@material-ui/core/Button';
import Container from '@material-ui/core/Container';
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';
import DialogActions from '@material-ui/core/DialogActions';
import Fab from '@material-ui/core/Fab';
import FormControl from '@material-ui/core/FormControl';
import FormHelperText from '@material-ui/core/FormHelperText';
import Grid from '@material-ui/core/Grid';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import Paper from '@material-ui/core/Paper';
import Typography from '@material-ui/core/Typography';

import qrcode from 'qrcode';
import { box_keyPair } from 'tweetnacl-ts';
import { codeBlock } from 'common-tags';
import { makeStyles } from '@material-ui/core/styles';

import { GetConnected } from './GetConnected';
import { IDevice, AppState } from '../Store';

const useStyles = makeStyles(theme => ({
  hidden: {
    display: 'none',
  },
  button: {
    margin: theme.spacing(1),
  },
  fabButton: {
    position: 'absolute',
    margin: '0 auto',
    left: 0,
    right: 0,
  },
  paper: {
    padding: theme.spacing(2),
  },
}));

export default function AddDevice() {
  const classes = useStyles();
  const [formOpen, setFormOpen] = React.useState(false);
  const [dialogOpen, setDialogOpen] = React.useState(false);
  const [error, setError] = React.useState('');
  const [name, setName] = React.useState('');
  const [qrCodeUri, setQrCodeUri] = React.useState('');
  const [configFileUri, setConfigFileUri] = React.useState('');
  
  var closeForm = () => {
    setFormOpen(false);
    setName('');
  };

  var addDevice = async (event: React.FormEvent) => {
    event.preventDefault();

    const keypair = box_keyPair();
    const b64PublicKey = window.btoa(String.fromCharCode(...(new Uint8Array(keypair.publicKey) as any)));
    const b64PrivateKey = window.btoa(String.fromCharCode(...(new Uint8Array(keypair.secretKey) as any)));

    const res = await fetch('/api/devices', {
      method: 'POST',
      body: JSON.stringify({
        name: name,
        publicKey: b64PublicKey,
      }),
    });
    if (res.status >= 400) {
      setError(await res.text());
      return;
    }
    const { device } = (await res.json()) as { device: IDevice };

    AppState.devices.push(device);

    const configFile = codeBlock`
      [Interface]
      PrivateKey = ${b64PrivateKey}
      Address = ${device.address}
      DNS = ${'1.1.1.1, 8.8.8.8'}

      [Peer]
      PublicKey = ${device.serverPublicKey}
      AllowedIPs = 0.0.0.0/1, 128.0.0.0/1, ::/0
      Endpoint = ${device.endpoint || `${window.location.hostname}:51820`}
    `;

    setQrCodeUri(await qrcode.toDataURL(configFile));
    setConfigFileUri(URL.createObjectURL(new Blob([configFile])));

    closeForm();
    setDialogOpen(true);
  };

  return (
    <React.Fragment>
      <Grid container spacing={3}>
        <Grid item xs></Grid>
        <Grid item xs={12} md={4} lg={6}>
          <Container hidden={formOpen}>
            <Fab color="secondary" aria-label="add" className={classes.fabButton} onClick={() => setFormOpen(true)}>
              <AddIcon />
            </Fab>
          </Container>
          <Paper hidden={!formOpen} className={classes.paper}>
              <form onSubmit={addDevice}>
                <FormControl error={error != ''} fullWidth>
                  <InputLabel htmlFor="device-name">Device Name</InputLabel>
                  <Input
                    id="device-name"
                    value={name}
                    onChange={(event) => setName(event.currentTarget.value)}
                    aria-describedby="device-name-text"
                  />
                  <FormHelperText id="device-name-text">{error}</FormHelperText>
                </FormControl>
                <Typography component="div" align="right">
                  <Button
                      color="secondary"
                      type="button"
                      onClick={closeForm}
                      className={classes.button}
                      >
                      Cancel
                  </Button>
                  <Button
                      color="primary"
                      variant="contained"
                      endIcon={<AddIcon />}
                      type="submit"
                      className={classes.button}
                      >
                      Next
                  </Button>
                </Typography>
              </form>
          </Paper>
        </Grid>
        <Grid item xs></Grid>
      </Grid>
      <Dialog
          disableBackdropClick
          disableEscapeKeyDown
          maxWidth="xl"
          open={dialogOpen}
      >
        <DialogTitle>Get Connected</DialogTitle>
        <DialogContent>
        <GetConnected
            qrCodeUri={qrCodeUri}
            configFileUri={configFileUri}
        />
        </DialogContent>
        <DialogActions>
        <Button color="secondary" variant="outlined" onClick={() => setDialogOpen(false)}>
            Done
        </Button>
        </DialogActions>
      </Dialog>
    </React.Fragment>
  );
}