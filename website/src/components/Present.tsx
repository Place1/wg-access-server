import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import React from 'react';
import { render, unmountComponentAtNode } from 'react-dom';

export function present<T>(content: (close: (result: T) => void) => React.ReactNode) {
  const root = document.createElement('div');
  document.body.appendChild(root);
  return new Promise<T>((resolve) => {
    const close = (result: T) => {
      unmountComponentAtNode(root);
      resolve(result);
    };
    render(<>{content(close)}</>, root);
  });
}

export function confirm(msg: string): Promise<boolean> {
  return present<boolean>((close) => (
    <Dialog open={true} onClose={() => close(false)}>
      <DialogTitle>Confirm</DialogTitle>
      <DialogContent>
        <DialogContentText>{msg}</DialogContentText>
      </DialogContent>
      <DialogActions>
        <Button onClick={() => close(false)} variant="contained" color="primary" autoFocus>
          Cancel
        </Button>
        <Button onClick={() => close(true)} variant="outlined" color="secondary">
          Ok
        </Button>
      </DialogActions>
    </Dialog>
  ));
}
