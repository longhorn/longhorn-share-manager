# syntax=docker/dockerfile:1.9.0

FROM registry.suse.com/bci/bci-base:15.6 AS build

RUN zypper -n ref && \
    zypper update -y

# Based on https://github.com/kubernetes-sigs/nfs-ganesha-server-and-external-provisioner/blob/master/deploy/docker/x86_64/Dockerfile
# Build ganesha from source, install it to /usr/local and a use multi stage build to have a smaller image
# RUN dnf install -y tar gcc cmake-3.14.2-1.fc30 autoconf libtool bison flex make gcc-c++ krb5-devel dbus-devel jemalloc-devel libnfsidmap-devel libnsl2-devel userspace-rcu-devel patch libblkid-devel
RUN zypper -n addrepo --refresh https://download.opensuse.org/repositories/system:/snappy/SLE_15/ snappy && \
    zypper -n addrepo --refresh https://download.opensuse.org/repositories/devel:/tools:/building/15.5/devel:tools:building.repo && \
    zypper --gpg-auto-import-keys ref

RUN zypper -n install autoconf bison curl cmake doxygen make git gcc-c++ flex Mesa-libGL-devel libdbus-1-3 \
    nfsidmap-devel liburcu-devel libblkid-devel e2fsprogs e2fsprogs-devel xfsprogs xfsprogs-devel \
    tar gzip dbus-1-devel lsb-release graphviz-devel libnsl-devel libcurl-devel libjson-c-devel libacl-devel && \
    rm -rf /var/cache/zypp/*

RUN curl -L https://github.com/longhorn/nfs-ganesha/archive/refs/tags/v5_20240430.tar.gz | tar zx \
	  && curl -L https://github.com/nfs-ganesha/ntirpc/archive/v5.0.tar.gz | tar zx \
	  && mv nfs-ganesha-5_20240430 nfs-ganesha-5.7 \
	  && rm -r nfs-ganesha-5.7/src/libntirpc \
	  && mv ntirpc-5.0 nfs-ganesha-5.7/src/libntirpc
WORKDIR /nfs-ganesha-5.7


# build ganesha only supporting nfsv4 and vfs
# Set NFS_V4_RECOV_ROOT to /tmp we don't support recovery in this release
# we disable dbus (-DUSE_DBUS=OFF) for the single share manager since we don't use dynamic exports
RUN mkdir -p /usr/local \
    && cmake -DCMAKE_BUILD_TYPE=Release -DBUILD_CONFIG=vfs_only \
    -DUSE_DBUS=OFF -DUSE_NFS3=OFF -DUSE_NLM=OFF -DUSE_RQUOTA=OFF -DUSE_9P=OFF -D_MSPAC_SUPPORT=OFF -DRPCBIND=OFF \
    -DUSE_RADOS_RECOV=OFF -DRADOS_URLS=OFF -DUSE_FSAL_VFS=ON -DUSE_FSAL_XFS=OFF \
    -DUSE_FSAL_PROXY_V4=OFF -DUSE_FSAL_PROXY_V3=OFF -DUSE_FSAL_LUSTRE=OFF -DUSE_FSAL_LIZARDFS=OFF \
    -DUSE_FSAL_KVSFS=OFF -DUSE_FSAL_CEPH=OFF -DUSE_FSAL_GPFS=OFF -DUSE_FSAL_PANFS=OFF -DUSE_FSAL_GLUSTER=OFF \
    -DUSE_GSS=NO -DHAVE_ACL_GET_FD_NP=ON -DHAVE_ACL_SET_FD_NP=ON \
    -DCMAKE_INSTALL_PREFIX=/usr/local src/ \
	  && make \
	  && make install
RUN mkdir -p /ganesha-extra \
    && mkdir -p /ganesha-extra/etc/dbus-1/system.d \
    && cp src/scripts/ganeshactl/org.ganesha.nfsd.conf /ganesha-extra/etc/dbus-1/system.d/

FROM registry.suse.com/bci/bci-base:15.6

ARG TARGETPLATFORM
RUN if [ "$TARGETPLATFORM" != "linux/amd64" ] && [ "$TARGETPLATFORM" != "linux/arm64" ]; then \
    echo "Error: Unsupported TARGETPLATFORM: $TARGETPLATFORM" && \
    exit 1; \
    fi
ENV ARCH ${TARGETPLATFORM#linux/}

RUN zypper addrepo --refresh https://download.opensuse.org/repositories/system:/snappy/SLE_15/ snappy && \
    zypper --gpg-auto-import-keys ref

# RUN microdnf install -y nano tar lsof e2fsprogs fuse-libs libss libblkid userspace-rcu dbus-x11 rpcbind hostname nfs-utils xfsprogs jemalloc libnfsidmap && microdnf clean all
RUN zypper -n install rpcbind hostname libblkid1 liburcu6 libjson-c* dbus-1-x11 dbus-1 nfsidmap-devel nfs-kernel-server nfs-client nfs4-acl-tools xfsprogs e2fsprogs awk && \
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
