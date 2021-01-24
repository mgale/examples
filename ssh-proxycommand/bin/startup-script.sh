#!/bin/bash

apt-get update
apt-get -y install rsyslog openssh-server sudo
service rsyslog restart
ssh-keygen -A
useradd -rm -d /home/ubuntu -s /bin/bash -G sudo -U -u 1000 ubuntu
service ssh restart
service ssh restart
touch /var/log/auth
tail -F /var/log/auth