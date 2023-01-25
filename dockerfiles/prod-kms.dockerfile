FROM alpine

# RUN apk add libc6-compat
# ENV GLIBC_REPO=https://github.com/sgerrand/alpine-pkg-glibc
# ENV GLIBC_VERSION=2.30-r0

# RUN set -ex && \
#     apk --update add libstdc++ curl ca-certificates && \
#     for pkg in glibc-${GLIBC_VERSION} glibc-bin-${GLIBC_VERSION}; \
#         do curl -sSL ${GLIBC_REPO}/releases/download/${GLIBC_VERSION}/${pkg}.apk -o /tmp/${pkg}.apk; done && \
#     apk add --allow-untrusted /tmp/*.apk && \
#     rm -v /tmp/*.apk && \
#     /usr/glibc-compat/sbin/ldconfig /lib /usr/glibc-compat/lib

ENV LOCAL=/usr/local
COPY build/tmkms ${LOCAL}/bin

ENTRYPOINT [ "tmkms" ]

# Create the image
# $ docker build -f Dockerfile-ubuntu-tmkms . -t tmkms_i
# To test only 1 command
# $ docker run --rm -it tmkms_i:v0.12.2
# To build container
# $ docker create --name tmkms -i -v $(pwd)/docker/kms-alice:/root/tmkms tmkms_i
# $ docker start tmkms
# To run server on it
# $ docker exec -it tmkms start
# In other shell, to query it
# $ docker exec -it tmkms version