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

import React from "react";
import {Card, CardContent, CardMedia, Grid, Typography} from "@material-ui/core";
import {withStyles} from "@material-ui/core/styles";

const styles = () => ({
    card: {
        height: "100%",
        display: "flex",
        flexDirection: "column",
    },
    cardMedia: {
        paddingTop: "56.25%", // 16:9
        maxWidth: "100%",
        maxHeight: "100%"
    },
    cardContent: {
        flexGrow: 1,
    }
});

class Item extends React.Component {
    render() {
        const {classes, data} = this.props;
        return (
            <Grid item sm={6} md={4} lg={3}>
                <Card className={classes.card}>
                    <CardMedia
                        className={classes.cardMedia}
                        image={"./app/assets/award.svg"}
                        title={data.name}
                    />
                    <CardContent className={classes.cardContent}>
                        <Typography gutterBottom variant="h5" component="h2">
                            {data.name}
                        </Typography>
                        <Typography>
                            {data.description}
                        </Typography>
                    </CardContent>
                </Card>
            </Grid>
        );
    }
}

export default withStyles(styles)(Item);
