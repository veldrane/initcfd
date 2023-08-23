#### Build

```cd src
go build cmd/initcfd.go
cd ../build
mv ../src/initcfd
docker build . -f Dockerfile -t czdcm-quay.lx.ifortuna.cz/shared-images/initcfd:<VERSION> --squash
```