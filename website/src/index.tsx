import React from 'react';
import ReactDOM from 'react-dom';
import Devices from './components/Devices';
import { view } from 'react-easy-state';
import 'typeface-roboto';
import './index.css';

const App = view(() => {
  return <Devices />;
});

ReactDOM.render(<App />, document.getElementById('root'));
