import React from 'react';
import { render, unmountComponentAtNode } from 'react-dom';
import Snackbar from '@mui/material/Snackbar';
import Alert from '@mui/material/Alert';

interface Props {
  intent: 'success' | 'info' | 'warning' | 'error';
  text: string;
}

export function toast(props: Props) {
  const root = document.createElement('div');
  document.body.appendChild(root);

  const onClose = () => {
    unmountComponentAtNode(root);
    document.body.removeChild(root);
  };

  render(
    <Snackbar
      open={true}
      autoHideDuration={3000}
      onClose={onClose}
      anchorOrigin={{ vertical: 'top', horizontal: 'right' }}
    >
      <Alert severity={props.intent} elevation={6} variant="filled" onClose={onClose}>
        {props.text}
      </Alert>
    </Snackbar>,
    root,
  );
}
