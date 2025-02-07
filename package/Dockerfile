# syntax=docker/dockerfile:1.12.1

FROM registry.suse.com/bci/bci-base:15.6 AS build

ARG BRANCH=main

RUN zypper -n ref && \
    zypper update -y

# Based on https://github.com/kubernetes-sigs/nfs-ganesha-server-and-external-provisioner/blob/master/deploy/docker/x86_64/Dockerfile
# Build ganesha from source, install it to /usr/local and a use multi stage build to have a smaller image
# RUN dnf install -y tar gcc cmake-3.14.2-1.fc30 autoconf libtool bison flex make gcc-c++ krb5-devel dbus-devel jemalloc-devel libnfsidmap-devel libnsl2-devel userspace-rcu-devel patch libblkid-devel
RUN zypper -n addrepo --refresh https://download.opensuse.org/repositories/system:/snappy/SLE_15/ snappy && \
    zypper -n addrepo --refresh https://download.opensuse.org/repositories/devel:/tools:/building/15.6/devel:tools:building.repo && \
    zypper --gpg-auto-import-keys ref

RUN zypper -n install autoconf bison curl cmake doxygen make git gcc14 gcc14-c++ flex Mesa-libGL-devel libdbus-1-3 \
    nfsidmap-devel liburcu-devel libblkid-devel e2fsprogs e2fsprogs-devel xfsprogs xfsprogs-devel \
    tar gzip dbus-1-devel lsb-release graphviz-devel libnsl-devel libcurl-devel libjson-c-devel libacl-devel jq && \
    rm -rf /var/cache/zypp/*

RUN git clone https://github.com/longhorn/dep-versions.git -b ${BRANCH} /usr/src/dep-versions

# Build nfs-ganesha
RUN export REPO_OVERRIDE="" && \
    export COMMIT_ID_OVERRIDE="" && \
    bash /usr/src/dep-versions/scripts/build-nfs-ganesha.sh "${REPO_OVERRIDE}" "${COMMIT_ID_OVERRIDE}"


FROM registry.suse.com/bci/bci-base:15.6

ARG TARGETPLATFORM
RUN if [ "$TARGETPLATFORM" != "linux/amd64" ] && [ "$TARGETPLATFORM" != "linux/arm64" ]; then \
    echo "Error: Unsupported TARGETPLATFORM: $TARGETPLATFORM" && \
    exit 1; \
    fi
ENV ARCH ${TARGETPLATFORM#linux/}

RUN zypper -n ref && \
    zypper update -y

RUN zypper addrepo --refresh https://download.opensuse.org/repositories/system:/snappy/SLE_15/ snappy && \
    zypper --gpg-auto-import-keys ref

# RUN microdnf install -y nano tar lsof e2fsprogs fuse-libs libss libblkid userspace-rcu dbus-x11 rpcbind hostname nfs-utils xfsprogs jemalloc libnfsidmap && microdnf clean all
RUN zypper -n install rpcbind hostname libblkid1 liburcu6 libjson-c* dbus-1-x11 dbus-1 nfsidmap-devel \
    nfs-kernel-server nfs-client nfs4-acl-tools xfsprogs e2fsprogs awk && \
    rm -rf /var/cache/zypp/*

RUN mkdir -p /var/run/dbus && mkdir -p /export

# add libs from /usr/local/lib64
RUN echo /usr/local/lib64 > /etc/ld.so.conf.d/local_libs.conf

# do not ask systemd for user IDs or groups (slows down dbus-daemon start)
RUN sed -i s/systemd// /etc/nsswitch.conf

# ganesha reads /etc/mtab for mounted volumes
RUN ln -sf /proc/self/mounts /etc/mtab

COPY --from=build /usr/local /usr/local/
COPY --from=build /ganesha-extra /
COPY package/bin/longhorn-share-manager-${ARCH} /longhorn-share-manager

# run ldconfig after libs have been copied
RUN ldconfig

# only expose the nfsd since for v4 only that is necessary
EXPOSE 2049/tcp

ENTRYPOINT ["/longhorn-share-manager"]
