#!/bin/bash

printf "Note you must be running this on wsl on the Ubuntu distro if you are on a windows machine"
printf "\nTo install this for windows run: wsl --install -d Ubuntu"
printf "\nStarting Installation\n\n"

UserInfo=$(<PostgresqlUserInfo.txt)
IFS=' '
read -ra UserInfoArr <<< "$UserInfo"

sudo mkdir /usr/local/postgresSetups
cd /usr/local/postgresSetups

pwd

sudo apt-get update
sudo apt-get install g++ -y
sudo apt-get install gcc -y --fix-missing
sudo apt-get install make -y
sudo apt-get install build-essential -y
sudo apt-get install tar -y
sudo apt-get install gzip -y
sudo apt-get install flex -y
sudo apt-get install m4 -y
sudo apt-get install bison -y
sudo apt-get install perl -y
sudo apt-get install zlib1g -y
sudo apt-get install libicu-dev -y
sudo apt-get install zlib1g-dev -y
sudo apt-get install pkgconf -y
sudo apt-get install git -y
sudo apt-get install net-tools

export GIT_TRACE_PACKET=1
export GIT_TRACE=1
export GIT_CURL_VERBOSE=1
git config --global http.postBuffer 157286400

sudo chown -R $USER /usr/local/postgresSetup
git clone https://git.postgresql.org/git/postgresql.git --depth 1;
cd postgresql/
git fetch --unshallow
./configure --without-readline
make
sudo su
make install
printf "${UserInfoArr[1]}\n${UserInfoArr[1]}\n${UserInfoArr[0]\n"0"\n"0"\n"0"\n"0"\n"Y"\n}" | sudo adduser postgres