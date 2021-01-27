#!/bin/bash

apt-get update
apt-get -y install rsyslog openssh-server sudo
service rsyslog restart
ssh-keygen -A
useradd -rm -d /home/ubuntu -s /bin/bash -G sudo -U -u 1000 ubuntu
chmod 600 /home/ubuntu/.ssh/*
chmod 700 /home/ubuntu/.ssh
service ssh restart
service ssh restart
touch /var/log/auth
echo "###############################################################"
echo "## Containers are ready !!!!"
echo "###############################################################"
tail -F /var/log/auth.log