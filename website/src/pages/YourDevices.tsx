import React from 'react';
import { view } from 'react-easy-state';
import AddDevice from '../components/AddDevice';
import Devices from '../components/Devices';

class YourDevices extends React.Component {
  render() {
    return (
      <>
        <Devices />
        <AddDevice />
      </>
    );
  }
}

export default view(YourDevices);
