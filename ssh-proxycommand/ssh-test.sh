#!/bin/bash

echo "###########################################################"
docker ps
echo "###########################################################"
echo "## SSH into ProxyServer"
ssh -F ssh-config proxyServer 'uptime && echo "Hostname: $HOSTNAME"'
echo "###########################################################"
echo "## SSH into Backend Server"
ssh -F ssh-config backendServer 'uptime && echo "Hostname: $HOSTNAME"'