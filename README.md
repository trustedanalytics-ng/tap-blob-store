# tap-blob-store

Catalog is a microservice developed to be a part of TAP platform.
It is a temporary staging area for arbitrary binary objects.

## REQUIREMENTS

### Binary
Blob store needs minio instance running.  It will be automatically downloaded for "make build_anywhere" and "make run" commands.

### Compilation
* git (for pulling repository)
* go >= 1.6

## Compilation
To build project:
```
  git clone https://github.com/intel-data/tap-blob-store
  cd tap-blob-store
  make build_anywhere
```
To build and run project:
```
  make run
```
Binaries are available in ./build directory.

## USAGE

To provide IP and port for the application, you have to setup system environment variables:
```
export BIND_ADDRESS=127.0.0.1
export PORT=80
```

Blob Store endpoints are documented in swagger.yaml file.
Below you can find sample Catalog usage.

#### Creating blob
```
curl -F blob_id=1 -F uploadfile=@file.txt http://127.0.0.1/api/v1/blobs --user admin:password
The blob has been successfully stored
```

This stores file "file.txt" in Blob Store.

#### Retrieving blob
```
curl -o retrieved_file.txt http://127.0.0.1/api/v1/blobs/1 --user admin:password
```
This retrieves blob and stores to the file "retrieved_file.txt".

#### Removing blob
```
curl -XDELETE http://127.0.0.1/api/v1/blobs/1 --user admin:password
```
