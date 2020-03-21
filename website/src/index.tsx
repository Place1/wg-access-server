import 'typeface-roboto';
import React from 'react';
import ReactDOM from 'react-dom';
import CssBaseline from '@material-ui/core/CssBaseline';
import Box from '@material-ui/core/Box';
import Navigation from './components/Navigation';
import {
  BrowserRouter as Router,
  Switch,
  Route,
} from 'react-router-dom';
import { observer } from 'mobx-react';
import { grpc } from './Api';
import { AppState } from './AppState';
import { YourDevices } from './pages/YourDevices';
import { AllDevices } from './pages/admin/AllDevices';

@observer
class App extends React.Component {

  async componentDidMount() {
    AppState.info = await grpc.server.info({});
  }

  render() {
    if (!AppState.info) {
      return <p>loading...</p>
    }
    return (
      <Router>
        <CssBaseline />
        <Navigation />
        <Box component="div" m={2}>
          <Switch>
            <Route exact path="/" component={YourDevices} />
            {AppState.info.isAdmin &&
              <>
                <Route exact path="/admin/all-devices" component={AllDevices} />
              </>
            }
          </Switch>
        </Box>
      </Router>
    );
  }
}

ReactDOM.render(<App />, document.getElementById('root'));
