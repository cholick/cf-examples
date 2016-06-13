### An sample bosh release

#### Setup

Initialize `config/private` file to

```yaml

---
blobstore_secret: 'does-not-matter'
blobstore:
  local:
    blobstore_path: /tmp/local_blobstore
```

Run `./scripts/blobs.sh` to download the proper packages to the local store. A real world release would
use s3 or a similar solution. See [configuring a blobstore](https://bosh.io/docs/create-release.html#config-blobstore)

#### References

See https://bosh.io/docs/create-release.html
