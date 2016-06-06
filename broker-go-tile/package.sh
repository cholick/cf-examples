#!/usr/bin/env bash -e

echo "Archiving source"
pushd ../broker-go/
zip -r ../broker-go-tile/resources/broker-go.zip * -x manifest.yml
#tar -czv --exclude='manifest.yml' -f ../broker-go-tile/resources/broker-go.tar.gz *
popd

echo "Do build"
tile build
