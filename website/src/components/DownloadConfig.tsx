import React from 'react';
import Fab from '@material-ui/core/Fab';
import VerifiedUserIcon from '@material-ui/icons/VerifiedUser';

interface Props {
  configFileUri: string;
}

export class DownloadConfig extends React.Component<Props> {
  downloadConfig = () => {
    console.log('downloading config file', this.props.configFileUri);
    const anchor = document.createElement('a');
    anchor.href = this.props.configFileUri;
    anchor.download = 'wireguard.conf';
    anchor.style.display = 'none';
    document.body.appendChild(anchor);
    anchor.click();
    document.body.removeChild(anchor);
  };

  render() {
    return (
      <Fab
        variant="extended"
        size="small"
        color="primary"
        style={{ padding: 30, borderRadius: 60 }}
        onClick={this.downloadConfig}
      >
        Download VPN Config
        <VerifiedUserIcon style={{ marginLeft: 15 }} />
      </Fab>
    );
  }
}
