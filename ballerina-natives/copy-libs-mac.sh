#!/usr/bin/env bash

BALLERINA_VERSION=$(ballerina version | awk '{print $2}')

mvn clean install

sudo cp target/cellery-1.0.0-SNAPSHOT.jar /Library/Ballerina/ballerina-${BALLERINA_VERSION}/bre/lib
sudo cp -r target/generated-balo/repo/celleryio /Library/Ballerina/ballerina-${BALLERINA_VERSION}/lib/repo
cp -r target/generated-balo/repo/celleryio ~/.ballerina/repo