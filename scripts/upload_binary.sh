#!/bin/bash

echo "Enter host"
read HOST
echo "Enter username"
read USER 
scp build/reputation-adapter-linux-amd64 $USER@$HOST:reputation-adapter-linux-amd64 
echo "Upload succesfull!"