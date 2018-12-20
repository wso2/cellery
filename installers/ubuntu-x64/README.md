# <img src="https://cdn1.iconfinder.com/data/icons/Futurosoft%20Icons%200.5.2/32x32/apps/ubuntu-logo.png"> Ubuntu-x64 Cellery Installer

You can generate Ubuntu OS product installer Cellery. Download the archive file and extract the archive file or clone the repository to a dedicated directory.

## Prerequisites
1. Apache Maven.
   * Download the latest version of [Maven for Ubuntu](https://maven.apache.org/download.cgi) and install.
2. Ballerina
   * Download latest [Ballerina](https://ballerina.io/downloads/) that need to build cellery natives. 

## Installer Generate Instructions
- Run `bash build-ubuntu-x64.sh <CELLERY-INSTALLER-VERSION>`
  ```
  eg : $ bash build-ubuntu-x64.sh 1.0.0
  ```
- This will generate the cellery installer `cellery-ubuntu-installer-x64-<CELLERY-INSTALLER-VERSION>.deb` file in target folder.
  ```
  eg : target/cellery-ubuntu-x64-1.0.0.deb
  ```
## Installion Instructions

- Installer will save cellery files in `/home/usr/local/bin/`

#### Using Ubuntu Software Center
- Double click and open build debian package `cellery-ubuntu-x64-<CELLERY-INSTALLER-VERSION>.deb` with Ubuntu Software Center. Then install Cellery using user interface.
#### Using Terminal
- Run `sudo dpkg -i cellery-ubuntu-installer-x64-<CELLERY-INSTALLER-VERSION>.deb`
 ```
  eg : $ sudo dpkg -i cellery-ubuntu-installer-x64-1.0.0.deb
 ```
## Test your Installation 
- Open terminal and run `cellery` to start with cellery commads. 

## Uninstall Cellery
- Run following command to uninstall cellery
  ```
  $ sudo apt-get purge cellery
  ```
