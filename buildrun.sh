set -x
rm bridge
set -ex
go build
sudo ./bridge
