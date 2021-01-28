#!/bin/bash

# Ensure keys have correct permissions
chmod 600 container-ssh-setup/*-privatekey

echo "###########################################################"
docker ps
echo "###########################################################"
echo "## SSH into proxyServer"
ssh -q -F ssh-config proxyServer 'uptime && echo "Hostname: $HOSTNAME"'
echo "###########################################################"
echo "## SSH into backendServer"
ssh -q -F ssh-config backendServer 'uptime && echo "Hostname: $HOSTNAME"'
