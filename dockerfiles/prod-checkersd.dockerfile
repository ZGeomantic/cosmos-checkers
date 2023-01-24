FROM alpine
ARG BUILDARCH

ENV LOCAL=/usr/local

COPY build/checkersd-linux-${BUILDARCH} ${LOCAL}/bin/checkersd

ENTRYPOINT [ "checkersd" ]

# Create the image
# $ docker build -f Dockerfile-ubuntu-prod . -t checkersd_i:v1
# To test only 1 command
# $ docker run --rm -it checkersd_i:v1
# To build container
# $ docker create --name checkersd -i -v $(pwd)/docker/val-alice:/root/.checkers checkersd_i:v1
# $ docker start checkersd
# To run server on it
# $ docker exec -it checkersd start
# In other shell, to query it
# $ docker exec -it checkersd version