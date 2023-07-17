#!/bin/sh
if [ -z "${COMMIT}" ]; then
    echo "-" > commit.txt
else
    echo "${COMMIT}" > commit.txt
fi
