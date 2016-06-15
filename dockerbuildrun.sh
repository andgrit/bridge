set -ex
docker build -t bridge .
docker run -it -p 80:80 bridge
