ARG ARCH=amd64

# Based on https://github.com/kubernetes-sigs/nfs-ganesha-server-and-external-provisioner/blob/master/deploy/docker/x86_64/Dockerfile
FROM ubuntu:18.04 AS build

# Build ganesha from source, install it to /usr/local and a use multi stage build to have a smaller image
# RUN dnf install -y tar gcc cmake-3.14.2-1.fc30 autoconf libtool bison flex make gcc-c++ krb5-devel dbus-devel jemalloc-devel libnfsidmap-devel libnsl2-devel userspace-rcu-devel patch libblkid-devel
RUN apt-get update && apt-get install -y \
  autoconf \
  build-essential \
  bison \
  curl \
  cmake \
  doxygen \
  make \
  git \
  g++ \
  flex \
  libboost-dev \
  libglu1-mesa-dev \
  libc6-dev \
  libkrb5-dev \
  libdbus-1-dev \
  libnfsidmap-dev \
  liburcu-dev \
  libblkid-dev \
  && rm -rf /var/lib/apt/lists/*

RUN curl -L https://github.com/longhorn/nfs-ganesha/archive/v3_20210302.tar.gz | tar zx \
	  && curl -L https://github.com/nfs-ganesha/ntirpc/archive/v3.4.tar.gz | tar zx \
	  && mv nfs-ganesha-3_20210302 nfs-ganesha-3.4 \
	  && rm -r nfs-ganesha-3.4/src/libntirpc \
	  && mv ntirpc-3.4 nfs-ganesha-3.4/src/libntirpc
WORKDIR /nfs-ganesha-3.4

# build ganesha only supporting nfsv4 and vfs
# Set NFS_V4_RECOV_ROOT to /tmp we don't support recovery in this release
# we disable dbus (-DUSE_DBUS=OFF) for the single share manager since we don't use dynamic exports
RUN mkdir -p /usr/local \
    && cmake -DCMAKE_BUILD_TYPE=Release -DBUILD_CONFIG=vfs_only \
    -DUSE_DBUS=OFF -DUSE_NFS3=OFF -DUSE_NLM=OFF -DUSE_RQUOTA=OFF -DUSE_9P=OFF -D_MSPAC_SUPPORT=OFF -DRPCBIND=OFF \
    -DUSE_RADOS_RECOV=OFF -DRADOS_URLS=OFF -DUSE_FSAL_VFS=ON -DUSE_FSAL_XFS=OFF \
    -DUSE_FSAL_PROXY_V4=OFF -DUSE_FSAL_PROXY_V3=OFF -DUSE_FSAL_LUSTRE=OFF -DUSE_FSAL_LIZARDFS=OFF \
    -DUSE_FSAL_KVSFS=OFF -DUSE_FSAL_CEPH=OFF -DUSE_FSAL_GPFS=OFF -DUSE_FSAL_PANFS=OFF -DUSE_FSAL_GLUSTER=OFF \
    -DCMAKE_INSTALL_PREFIX=/usr/local src/ \
    && sed -i 's|@SYSSTATEDIR@/lib/nfs/ganesha|/tmp|' src/include/config-h.in.cmake \
	  && make \
	  && make install
RUN mkdir -p /ganesha-extra \
    && mkdir -p /ganesha-extra/etc/dbus-1/system.d \
    && cp src/scripts/ganeshactl/org.ganesha.nfsd.conf /ganesha-extra/etc/dbus-1/system.d/

FROM ubuntu:18.04
# RUN microdnf install -y nano tar lsof e2fsprogs fuse-libs libss libblkid userspace-rcu dbus-x11 rpcbind hostname nfs-utils xfsprogs jemalloc libnfsidmap && microdnf clean all
RUN apt-get update && apt-get install -y \
  rpcbind \
  hostname \
  libblkid1 \
  liburcu6 \
  dbus-x11 \
  dbus \
  libnfsidmap2 \
  nfs-common \
  && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /var/run/dbus \
    && mkdir -p /export

# add libs from /usr/local/lib64
RUN echo /usr/local/lib64 > /etc/ld.so.conf.d/local_libs.conf

# do not ask systemd for user IDs or groups (slows down dbus-daemon start)
RUN sed -i s/systemd// /etc/nsswitch.conf

# ganesha reads /etc/mtab for mounted volumes
RUN ln -sf /proc/self/mounts /etc/mtab

COPY --from=build /usr/local /usr/local/
COPY --from=build /ganesha-extra /
COPY bin/longhorn-share-manager /longhorn-share-manager

# run ldconfig after libs have been copied
RUN ldconfig

# onyl expose the nfsd since for v4 only that is necessary
EXPOSE 2049/tcp

ENTRYPOINT ["/longhorn-share-manager"]
