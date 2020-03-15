import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { getCookie } from '../Cookies';
import { AppState } from '../AppState';
import { NavLink } from 'react-router-dom';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import Link from '@material-ui/core/Link';
import Button from '@material-ui/core/Button';
import Chip from '@material-ui/core/Chip';

const useStyles = makeStyles(theme => ({
  title: {
    flexGrow: 1,
  },
}));

export default function Navigation() {
  const classes = useStyles();
  const hasAuthCookie = !!getCookie('auth-session');
  return (
    <AppBar position="static">
      <Toolbar>
        <Typography variant="h6" className={classes.title}>
          <span>Your Devices</span>
          {AppState.info?.isAdmin &&
            <Chip label="admin" color="secondary" variant="outlined" size="small" style={{ marginLeft: 20, background: 'white' }} />
          }
        </Typography>

        {AppState.info?.isAdmin &&
          <Link to="/admin/all-devices" color="inherit" component={NavLink}>
            <Button color="inherit">
              All Devices
            </Button>
          </Link>
        }

        {hasAuthCookie &&
          <Link href="/signout" color="inherit">
            <Button color="inherit">
              Logout
            </Button>
          </Link>
        }
      </Toolbar>
    </AppBar>
  );
}
