import React from 'react';
import ReactDOM from 'react-dom';
import CssBaseline from '@material-ui/core/CssBaseline';
import Box from '@material-ui/core/Box';
import AddDevice from './components/AddDevice';
import Devices from './components/Devices';
import Navigation from './components/Navigation';
import { view } from 'react-easy-state';
import 'typeface-roboto';

const App = view(() => {
  return (
    <React.Fragment>
      <CssBaseline />
      <Navigation />
      <Box component="div" m={3}>
        <Devices />
        <AddDevice />
      </Box>
    </React.Fragment>
  );
});

ReactDOM.render(<App />, document.getElementById('root'));
