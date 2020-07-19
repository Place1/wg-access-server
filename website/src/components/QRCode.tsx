import React from 'react';
import qrcode from 'qrcode';
import { lazy } from '../Util';
import { CircularProgress } from '@material-ui/core';

interface Props {
  content: string;
}

export class QRCode extends React.Component<Props> {
  uri = lazy(async () => {
    return await qrcode.toDataURL(this.props.content);
  });

  render() {
    if (!this.uri.current) {
      return <CircularProgress color="secondary" />;
    }
    return <img alt="WireGuard QR code" src={this.uri.current} />;
  }
}
