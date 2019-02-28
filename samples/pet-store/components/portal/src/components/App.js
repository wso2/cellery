/*
 * Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import Catalog from "./catalog";
import Orders from "./orders";
import React from "react";
import {withStyles} from "@material-ui/core/styles";
import {AppBar, Toolbar, Typography} from "@material-ui/core";
import {Pets} from "@material-ui/icons";
import {Switch, Redirect, Route, withRouter} from "react-router-dom";

const styles = (theme) => ({
    appBar: {
        position: 'relative',
    },
    icon: {
        marginRight: theme.spacing.unit * 2,
    },
    heroContent: {
        maxWidth: 600,
        margin: '0 auto',
        padding: `${theme.spacing.unit * 8}px 0 ${theme.spacing.unit * 6}px`,
    },
    heroButtons: {
        marginTop: theme.spacing.unit * 4,
    },
    layout: {
        width: 'auto',
        marginLeft: theme.spacing.unit * 3,
        marginRight: theme.spacing.unit * 3,
        [theme.breakpoints.up(1100 + theme.spacing.unit * 3 * 2)]: {
            width: 1100,
            marginLeft: 'auto',
            marginRight: 'auto',
        },
    },
    cardGrid: {
        padding: `${theme.spacing.unit * 8}px 0`,
    },
    card: {
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
    },
    cardMedia: {
        paddingTop: '56.25%', // 16:9
    },
    cardContent: {
        flexGrow: 1,
    },
    footer: {
        backgroundColor: theme.palette.background.paper,
        padding: theme.spacing.unit * 6,
    },
});

class App extends React.Component {
    state = {
        open: false,
    };

    render() {
        const {classes, location} = this.props;

        const pages = [
            "/",
            "/orders"
        ];
        let selectedIndex = 0;
        for (let i = 0; i < pages.length; i++) {
            if (location.pathname.startsWith(pages[i])) {
                selectedIndex = i;
            }
        }

        return (
            <div className={classes.root}>
                <AppBar position="static" className={classes.appBar}>
                    <Toolbar>
                        <Pets className={classes.icon} />
                        <Typography variant="h6" color="inherit" noWrap>
                            Pet Store
                        </Typography>
                    </Toolbar>
                </AppBar>
                <main>
                    <Switch>
                        <Route exact path={pages[0]} component={Catalog}/>
                        <Route exact path={pages[0]} component={Orders}/>
                        <Redirect from={"*"} to={"/"}/>
                    </Switch>
                </main>
            </div>
        );
    }
}

export default withStyles(styles, {withTheme: true})(withRouter(App));
