import { ButtonGroup } from '@material-ui/core';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import Paper from '@material-ui/core/Paper';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import { GetApp } from '@material-ui/icons';
import Laptop from '@material-ui/icons/Laptop';
import PhoneIphone from '@material-ui/icons/PhoneIphone';
import React from 'react';
import { AppState } from '../AppState';
import { isMobile } from '../Platform';
import { download } from '../Util';
import { LinuxIcon, MacOSIcon, WindowsIcon } from './Icons';
import { QRCode } from './QRCode';
import { TabPanel } from './TabPanel';

interface Props {
  configFile: string;
}

export class GetConnected extends React.Component<Props> {
  state = {
    currentTab: isMobile() ? 'mobile' : 'desktop',
  };

  go = (href: string) => {
    window.open(href, '__blank', 'noopener noreferrer');
  };

  download = () => {
    const info = AppState.info!;
    download({
      filename: info.filename.length > 0 ? info.filename + '.conf' : 'WireGuard.conf',
      content: this.props.configFile,
    });
  };

  getqr = async () => {
    return;
  };

  render() {
    return (
      <React.Fragment>
        <Paper>
          <Tabs
            value={this.state.currentTab}
            onChange={(_, currentTab) => this.setState({ currentTab })}
            indicatorColor="primary"
            textColor="primary"
            variant="fullWidth"
          >
            <Tab icon={<Laptop />} value="desktop" />
            <Tab icon={<PhoneIphone />} value="mobile" />
          </Tabs>
        </Paper>

        <TabPanel for="desktop" value={this.state.currentTab}>
          <Grid container direction="row" justify="space-around" alignItems="center">
            <List>
              <ListItem>
                <ListItemText style={{ width: 300 }} primary="1. Install the WireGuard App" />
                <ButtonGroup size="large" color="primary" aria-label="large outlined primary button group">
                  <Button onClick={() => this.go('https://www.WireGuard.com/install/')}>
                    <LinuxIcon />
                  </Button>
                  <Button
                    onClick={() => this.go('https://www.wireguard.com/install/')}
                  >
                    <WindowsIcon />
                  </Button>
                  <Button onClick={() => this.go('https://www.wireguard.com/install/#macos-app-store')}>
                    <MacOSIcon />
                  </Button>
                </ButtonGroup>
              </ListItem>
              <ListItem>
                <ListItemText style={{ width: 300 }} primary="2. Download your connection file" />
                <Button variant="outlined" color="primary" onClick={this.download}>
                  <GetApp /> Connection File
                </Button>
              </ListItem>
              <ListItem>
                <ListItemText style={{ width: 300 }} primary="3. Import your connection file in the App" />
              </ListItem>
            </List>
          </Grid>
        </TabPanel>

        <TabPanel for="mobile" value={this.state.currentTab}>
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
              <QRCode content={this.props.configFile} />
            </Grid>
          </Grid>
        </TabPanel>
      </React.Fragment>
    );
  }
}
