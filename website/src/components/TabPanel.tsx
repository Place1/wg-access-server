import React, { Component, PropsWithChildren } from 'react';

interface Props {
  for: any;
  value: any;
}

export class TabPanel extends Component<PropsWithChildren<Props>, any> {
  render() {
    return (
      <div style={{ padding: '1.5rem 1rem' }} hidden={this.props.for !== this.props.value}>
        {this.props.children}
      </div>
    );
  }
}
