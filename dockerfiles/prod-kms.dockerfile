FROM alpine

COPY build/tmkms ${LOCAL}/bin

ENTRYPOINT [ "tmkms" ]

# Create the image
# $ docker build -f Dockerfile-ubuntu-tmkms . -t tmkms_i
# To test only 1 command
# $ docker run --rm -it tmkms_i
# To build container
# $ docker create --name tmkms -i -v $(pwd)/docker/kms-alice:/root/tmkms tmkms_i
# $ docker start tmkms
# To run server on it
# $ docker exec -it tmkms start
# In other shell, to query it
# $ docker exec -it tmkms version