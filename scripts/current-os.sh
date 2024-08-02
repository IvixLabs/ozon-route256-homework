#!/bin/sh

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
       echo "linux"
       exit
fi

if [[ "$OSTYPE" == "darwin"* ]]; then
        echo "darwin"
        exit
fi

#if [[ "$OSTYPE" == "cygwin" ]]; then
        # POSIX compatibility layer and Linux environment emulation for Windows
#fi

#if [[ "$OSTYPE" == "msys" ]]; then
        # Lightweight shell and GNU utilities compiled for Windows (part of MinGW)
#fi

#if [[ "$OSTYPE" == "win32" ]]; then
        # I'm not sure this can happen.
#fi

if [[ "$OSTYPE" == "freebsd"* ]]; then
        echo "freebsd"
        exit
fi

echo $OSTYPE
