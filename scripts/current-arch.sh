#!/bin/sh

ARCH=$(arch)

if [ $ARCH = 'x86_64' ]
then
  echo 'amd64'
  exit
fi

echo $ARCH
