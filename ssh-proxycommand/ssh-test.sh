#!/bin/bash

echo "###########################################################"
docker ps
echo "###########################################################"
echo "## SSH into ProxyServer"
ssh -q -F ssh-config proxyServer 'uptime && echo "Hostname: $HOSTNAME"'
echo "###########################################################"
echo "## SSH into Backend Server"
ssh -q -F ssh-config backendServer 'uptime && echo "Hostname: $HOSTNAME"'
