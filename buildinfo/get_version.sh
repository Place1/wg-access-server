#!/bin/sh
if  [ -z "$VERSION" ] ; then
    echo "development" > version.txt
else
    echo "${VERSION}" > version.txt
fi
