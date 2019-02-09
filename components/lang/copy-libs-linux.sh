#!/usr/bin/env bash

BALLERINA_VERSION=$(ballerina version | awk '{print $2}')

mvn clean install

sudo cp target/cellery-0.1.0-SNAPSHOT.jar /usr/lib/ballerina/ballerina-${BALLERINA_VERSION}/bre/lib
sudo cp -r target/generated-balo/repo/celleryio /usr/lib/ballerina/ballerina-${BALLERINA_VERSION}/lib/repo
cp -r target/generated-balo/repo/celleryio ~/.ballerina/repo
