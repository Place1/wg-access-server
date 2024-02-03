import React from 'react';
import makeStyles from '@mui/styles/makeStyles';
import { getCookie } from '../Cookies';
import { AppState } from '../AppState';
import { NavLink } from 'react-router-dom';
import AppBar from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import Link from '@mui/material/Link';
import Button from '@mui/material/Button';
import Chip from '@mui/material/Chip';
import VpnKey from "@mui/icons-material/VpnKey";

const useStyles = makeStyles((theme) => ({
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
          <Link to="/" color="inherit" component={NavLink}>
            <VpnKey /> wg-access-server
          </Link>
          {AppState.info?.isAdmin && (
            <Chip
              label="admin"
              color="secondary"
              variant="outlined"
              size="small"
              style={{ marginLeft: 20, background: 'white' }}
            />
          )}
        </Typography>

        {AppState.info?.isAdmin && (
          <Link to="/admin/all-devices" color="inherit" component={NavLink}>
            <Button color="inherit">All Devices</Button>
          </Link>
        )}

        {hasAuthCookie && (
          <Link href="/signout" color="inherit">
            <Button color="inherit">Logout</Button>
          </Link>
        )}
      </Toolbar>
    </AppBar>
  );
}
