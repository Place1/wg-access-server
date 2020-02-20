import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import Link from '@material-ui/core/Link';
import Button from '@material-ui/core/Button';
import { getCookie } from '../Cookies';

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
        Your Devices
        </Typography>
        {hasAuthCookie &&
          <Link color="inherit" href="/signout">
            <Button color="inherit">
              Logout
            </Button>
          </Link>
        }
      </Toolbar>
    </AppBar>
  );
}
