import React from 'react';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import Paper from '@material-ui/core/Paper';
import Grid from '@material-ui/core/Grid';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import Typography from '@material-ui/core/Typography';
import { MacOSIcon, IosIcon, WindowsIcon, LinuxIcon, AndroidIcon } from './Icons';
import { TabPanel } from './TabPanel';
import { Platform, getPlatform } from '../Platform';
import { DownloadConfig } from './DownloadConfig';
import { DownloadLink } from './DownloadLink';

interface Props {
  qrCodeUri: string;
  configFileUri: string;
}

export class GetConnected extends React.Component<Props> {
  state = {
    platform: getPlatform(),
  };

  render() {
    return (
      <React.Fragment>
        <Paper>
          <Tabs
            value={this.state.platform}
            onChange={(_, platform) => this.setState({ platform })}
            indicatorColor="primary"
            textColor="primary"
            variant="scrollable"
            scrollButtons="auto"
          >
            <Tab icon={<LinuxIcon />} value={Platform.Linux} />
            <Tab icon={<MacOSIcon />} value={Platform.Mac} />
            <Tab icon={<WindowsIcon />} value={Platform.Windows} />
            <Tab icon={<IosIcon />} value={Platform.Ios} />
            <Tab icon={<AndroidIcon />} value={Platform.Android} />
          </Tabs>
        </Paper>
        <TabPanel for={Platform.Linux} value={this.state.platform}>
          <Grid container direction="row" justify="space-around" alignItems="center">
            <Grid item xs={12} sm={6}>
              <List>
                <ListItem>
                  <ListItemText primary="1. Install WireGuard for Linux" />
                </ListItem>
                <ListItem>
                  <ListItemText primary="2. Download your connection file" />
                </ListItem>
                <ListItem>
                  <ListItemText
                    primary="3. Copy it to /etc/wireguard/wg0.conf"
                    secondary="This will allow you to use wg-quick to bring the interface up and down easily."
                  />
                </ListItem>
              </List>
            </Grid>
            <Grid item direction="column" container spacing={3} xs={12} sm={6}>
              <Grid item>
                <DownloadConfig configFileUri={this.props.configFileUri} />
              </Grid>
              <Grid item>
                <DownloadLink
                  label="Download WireGuard"
                  href="https://www.wireguard.com/install/"
                  icon={<LinuxIcon />}
                />
              </Grid>
            </Grid>
          </Grid>
        </TabPanel>
        <TabPanel for={Platform.Mac} value={this.state.platform}>
          <Grid container direction="row" justify="space-around" alignItems="center">
            <Grid item xs={12} sm={6}>
              <List>
                <ListItem>
                  <ListItemText primary="1. Install WireGuard for MacOS" />
                </ListItem>
                <ListItem>
                  <ListItemText primary="2. Download your connection file" />
                </ListItem>
                <ListItem>
                  <ListItemText primary="3. Add tunnel from file" />
                </ListItem>
              </List>
            </Grid>
            <Grid item direction="column" container spacing={3} xs={12} sm={6}>
              <Grid item>
                <DownloadConfig configFileUri={this.props.configFileUri} />
              </Grid>
              <Grid item>
                <DownloadLink
                  label="Download WireGuard"
                  href="https://itunes.apple.com/us/app/wireguard/id1451685025?ls=1&mt=12"
                  icon={<MacOSIcon />}
                />
              </Grid>
            </Grid>
          </Grid>
        </TabPanel>
        <TabPanel for={Platform.Ios} value={this.state.platform}>
          <Grid container direction="row" justify="space-around" alignItems="center">
            <Grid item>
              <List>
                <ListItem>
                  <ListItemText primary="1. Install the WireGuard app" />
                </ListItem>
                <ListItem>
                  <ListItemText primary="2. Add a tunnel" />
                </ListItem>
                <ListItem>
                  <ListItemText primary="3. Create from QR code" />
                </ListItem>
              </List>
            </Grid>
            <Grid item>
              <img alt="wireguard qr code" src={this.props.qrCodeUri} />
            </Grid>
          </Grid>
        </TabPanel>
        <TabPanel for={Platform.Android} value={this.state.platform}>
          <Grid container direction="row" justify="space-around" alignItems="center">
            <Grid item>
              <List>
                <ListItem>
                  <ListItemText primary="1. Install the WireGuard app" />
                </ListItem>
                <ListItem>
                  <ListItemText primary="2. Add a tunnel" />
                </ListItem>
                <ListItem>
                  <ListItemText primary="3. Create from QR code" />
                </ListItem>
              </List>
            </Grid>
            <Grid item>
              <img alt="wireguard qr code" src={this.props.qrCodeUri} />
            </Grid>
          </Grid>
        </TabPanel>
        <TabPanel for={Platform.Windows} value={this.state.platform}>
          <Grid container direction="row" justify="space-around" alignItems="center">
            <Grid item xs={12} sm={6}>
              <List>
                <ListItem>
                  <ListItemText primary="1. Install WireGuard for Windows" />
                </ListItem>
                <ListItem>
                  <ListItemText primary="2. Download your connection file" />
                </ListItem>
                <ListItem>
                  <ListItemText primary="3. Add tunnel from file" />
                </ListItem>
              </List>
            </Grid>
            <Grid item direction="column" container spacing={3} xs={12} sm={6}>
              <Grid item>
                <DownloadConfig configFileUri={this.props.configFileUri} />
              </Grid>
              <Grid item>
                <DownloadLink
                  label="Download WireGuard"
                  href="https://download.wireguard.com/windows-client/wireguard-amd64-0.0.32.msi"
                  icon={<WindowsIcon />}
                />
              </Grid>
            </Grid>
          </Grid>
        </TabPanel>
        <Grid container justify="center">
          <Typography style={{ fontStyle: 'italic', maxWidth: 600 }}>
            The VPN configuration file or QR code will not be available again.
            <br />
            If you lose your connection settings or reset your device, you can remove and re-add it to generate a new
            connection file or QR code.
            <br />
            They contain your WireGuard Private Key and should never be shared.
          </Typography>
        </Grid>
      </React.Fragment>
    );
  }
}
