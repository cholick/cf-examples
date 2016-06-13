#!/usr/bin/env bash

set -e

export local_store_path=/tmp/local_blobstore

mkdir -p $local_store_path

echo "Real world would use s3 or similar solution for reproducibility."
echo "    See https://bosh.io/docs/create-release.html#config-blobstore"

# `bosh upload blobs` generated object_id
wget "https://storage.googleapis.com/golang/go1.6.2.linux-amd64.tar.gz" \
    -O "$local_store_path/17c216ff-d3a3-4bde-93e3-f9b52c5eb6e6"
