import { ButtonGroup } from '@mui/material';
import Button from '@mui/material/Button';
import Grid from '@mui/material/Grid';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemText from '@mui/material/ListItemText';
import Paper from '@mui/material/Paper';
import Tab from '@mui/material/Tab';
import Tabs from '@mui/material/Tabs';
import { GetApp } from '@mui/icons-material';
import Laptop from '@mui/icons-material/Laptop';
import PhoneIphone from '@mui/icons-material/PhoneIphone';
import React, { Component, PropsWithChildren } from 'react';
import { AppState } from '../AppState';
import { isMobile } from '../Platform';
import { download } from '../Util';
import { LinuxIcon, MacOSIcon, WindowsIcon } from './Icons';
import { QRCode } from './QRCode';
import { TabPanel } from './TabPanel';

interface Props {
  configFile: string;
  showMobile: boolean;
}

export class GetConnected extends Component<PropsWithChildren<Props>, any> {
  state = {
    currentTab: isMobile() && this.props.showMobile ? 'mobile' : 'desktop',
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
            {this.props.showMobile && <Tab icon={<PhoneIphone />} value="mobile" />}
          </Tabs>
        </Paper>

        <TabPanel for="desktop" value={this.state.currentTab}>
          <Grid container direction="row" justifyContent="space-around" alignItems="center">
            <List>
              <ListItem>
                <ListItemText style={{ width: 300 }} primary="1. Install the WireGuard App" />
                <ButtonGroup size="large" color="primary" aria-label="large outlined primary button group">
                  <Button onClick={() => this.go('https://www.WireGuard.com/install/')}>
                    <LinuxIcon />
                  </Button>
                  <Button onClick={() => this.go('https://www.wireguard.com/install/')}>
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

        {this.props.showMobile && (
          <TabPanel for="mobile" value={this.state.currentTab}>
            <Grid container direction="row" justifyContent="space-around" alignItems="center">
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
        )}
      </React.Fragment>
    );
  }
}
