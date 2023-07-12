echo $KS_LABORATORY_BACKEND_WRITE_PACKAGES | docker login ghcr.io -u USERNAME --password-stdin
docker build -f Dockerfile.lab -t "ghcr.io/koheiyamayama/ks-laboratory-backend:$(git rev-parse HEAD)" .
docker push "ghcr.io/koheiyamayama/ks-laboratory-backend:$(git rev-parse HEAD)"
