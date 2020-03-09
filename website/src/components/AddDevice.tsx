import React from 'react';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormControl from '@material-ui/core/FormControl';
import FormHelperText from '@material-ui/core/FormHelperText';
import Grid from '@material-ui/core/Grid';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import Paper from '@material-ui/core/Paper';
import AddIcon from '@material-ui/icons/Add';
import Typography from '@material-ui/core/Typography';
import qrcode from 'qrcode';
import { makeStyles } from '@material-ui/core/styles';
import { codeBlock } from 'common-tags';
import { box_keyPair } from 'tweetnacl-ts';
import { AppState } from '../Store';
import { GetConnected } from './GetConnected';
import { grpc } from '../Api';


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
  const [dialogOpen, setDialogOpen] = React.useState(false);
  const [error, setError] = React.useState('');
  const [name, setName] = React.useState('');
  const [qrCodeUri, setQrCodeUri] = React.useState('');
  const [configFileUri, setConfigFileUri] = React.useState('');

  const reset = () => {
    setName('');
  };

  const addDevice = async (event: React.FormEvent) => {
    event.preventDefault();

    const keypair = box_keyPair();
    const publicKey = window.btoa(String.fromCharCode(...(new Uint8Array(keypair.publicKey) as any)));
    const privateKey = window.btoa(String.fromCharCode(...(new Uint8Array(keypair.secretKey) as any)));

    try {
      const device = await grpc.devices.addDevice({ name, publicKey });
      AppState.devices.push(device);
      const configFile = codeBlock`
        [Interface]
        PrivateKey = ${privateKey}
        Address = ${device.address}
        DNS = ${info.hostVpnIp}

        [Peer]
        PublicKey = ${info.publicKey}
        AllowedIPs = 0.0.0.0/1, 128.0.0.0/1, ::/0
        Endpoint = ${`${info.host?.value || window.location.hostname}:${info.port || '51820'}`}
      `;
      setQrCodeUri(await qrcode.toDataURL(configFile));
      setConfigFileUri(URL.createObjectURL(new Blob([configFile])));
      reset();
      setDialogOpen(true);
    } catch (error) {
      console.log(error);
      setError('failed');
    }
  };

  return (
    <React.Fragment>
      <Grid container spacing={3}>
        <Grid item xs></Grid>
        <Grid item xs={12} md={4} lg={6}>
          <Paper className={classes.paper}>
            <h2>Add A Device</h2>
            <form onSubmit={addDevice}>
              <FormControl error={error !== ''} fullWidth>
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
                  onClick={reset}
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
                  Add
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
