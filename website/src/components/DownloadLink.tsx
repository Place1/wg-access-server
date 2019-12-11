import React from 'react';
import Fab from '@material-ui/core/Fab';

interface Props {
  label: string;
  icon: React.ReactNode;
  href: string;
}

export class DownloadLink extends React.Component<Props> {
  render() {
    return (
      <a href={this.props.href} target="__blank" rel="noopener noreferrer">
        <Fab variant="extended" size="small" color="primary" style={{ padding: 30, borderRadius: 60 }}>
          <span>{this.props.label}</span>
          <span style={{ marginLeft: 15 }}>{this.props.icon}</span>
        </Fab>
      </a>
    );
  }
}
