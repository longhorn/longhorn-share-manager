ARG ARCH=amd64

# Based on https://github.com/kubernetes-sigs/nfs-ganesha-server-and-external-provisioner/blob/master/deploy/docker/x86_64/Dockerfile
FROM fedora:30 AS build

# TODO: considering switching this to ubuntu, we lose the multi stage build and the image size will increase
# Build ganesha from source, install it to /usr/local and a use multi stage build to have a smaller image
# Set NFS_V4_RECOV_ROOT to /export

RUN dnf install -y tar gcc cmake-3.14.2-1.fc30 autoconf libtool bison flex make gcc-c++ krb5-devel dbus-devel jemalloc-devel libnfsidmap-devel libnsl2-devel userspace-rcu-devel patch libblkid-devel
RUN curl -L https://github.com/nfs-ganesha/nfs-ganesha/archive/V3.3.tar.gz | tar zx \
	  && curl -L https://github.com/nfs-ganesha/ntirpc/archive/v3.3.tar.gz | tar zx \
	  && rm -r nfs-ganesha-3.3/src/libntirpc \
	  && mv ntirpc-3.3 nfs-ganesha-3.3/src/libntirpc
WORKDIR /nfs-ganesha-3.3

# build ganesha only supporting nfsv4 and vfs
RUN mkdir -p /usr/local \
    && cmake -DCMAKE_BUILD_TYPE=Release -DBUILD_CONFIG=vfs_only \
    -DUSE_DBUS=ON -DUSE_NFS3=OFF -DUSE_NLM=OFF -DUSE_RQUOTA=OFF -DUSE_9P=OFF -D_MSPAC_SUPPORT=OFF -DRPCBIND=OFF \
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

FROM registry.fedoraproject.org/fedora-minimal:30
RUN microdnf install -y nano tar lsof e2fsprogs fuse-libs libss libblkid userspace-rcu dbus-x11 rpcbind hostname nfs-utils xfsprogs jemalloc libnfsidmap && microdnf clean all
# ext4 filesystem tools are not present on the minimal disk, required packages: e2fsprogs fuse-libs libss
# add lsof for debugging

RUN mkdir -p /var/run/dbus \
    && mkdir -p /export

# add libs from /usr/local/lib64
RUN echo /usr/local/lib64 > /etc/ld.so.conf.d/local_libs.conf

# do not ask systemd for user IDs or groups (slows down dbus-daemon start)
RUN sed -i s/systemd// /etc/nsswitch.conf

COPY --from=build /usr/local /usr/local/
COPY --from=build /ganesha-extra /
COPY bin/longhorn-share-manager /longhorn-share-manager

# run ldconfig after libs have been copied
RUN ldconfig

# onyl expose the nfsd since for v4 only that is necessary
EXPOSE 2049/tcp

# TODO: finish the docker file to run the share manager
ENTRYPOINT ["/longhorn-share-manager"]
