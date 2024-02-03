import React from 'react';
import Button from '@mui/material/Button';
import Popover from '@mui/material/Popover';

interface Props {
  label: string;
  children: React.ReactNode;
}

export class PopoverDisplay extends React.Component<Props> {
  state = {
    anchorEl: undefined as any,
  };

  render() {
    return (
      <React.Fragment>
        <Button
          size="small"
          variant="outlined"
          color="secondary"
          style={{ padding: 0 }}
          onClick={(event) => this.setState({ anchorEl: event.currentTarget })}
        >
          {this.props.label}
        </Button>
        <Popover
          open={Boolean(this.state.anchorEl)}
          anchorEl={this.state.anchorEl}
          onClose={() => this.setState({ anchorEl: undefined })}
          anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
          transformOrigin={{ vertical: 'top', horizontal: 'center' }}
        >
          <div style={{ padding: '2rem' }}>{this.props.children}</div>
        </Popover>
      </React.Fragment>
    );
  }
}
