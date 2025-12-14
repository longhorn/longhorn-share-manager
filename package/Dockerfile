# syntax=docker/dockerfile:1.20.0
FROM registry.suse.com/bci/golang:1.25 AS app_builder

WORKDIR /app

# Copy the build script and source code
COPY . /app

# Make the build script executable
RUN chmod +x /app/scripts/build

# Run the build script
RUN /app/scripts/build

FROM registry.suse.com/bci/bci-base:16.0 AS lib_builder

ARG TARGETPLATFORM
ARG SRC_BRANCH=master
ARG SRC_TAG

RUN zypper -n ref && \
    zypper update -y

# Based on https://github.com/kubernetes-sigs/nfs-ganesha-server-and-external-provisioner/blob/master/deploy/docker/x86_64/Dockerfile
# Build ganesha from source, install it to /usr/local and a use multi stage build to have a smaller image
RUN zypper -n install autoconf bison curl git tar gzip jq libcurl-devel libjson-c-devel libacl-devel lsb-release \
                      libdbus-1-3 liburcu-devel libblkid-devel e2fsprogs e2fsprogs-devel xfsprogs xfsprogs-devel \
                      dbus-1-devel flex Mesa-libGL-devel nfsidmap-devel

RUN for i in {1..10}; do \
        # Avoid using mirrors in Russia and Belarus
        zypper -n addrepo --refresh "https://download.opensuse.org/repositories/devel:libraries:c_c++/openSUSE_Factory/?AVOID_COUNTRY=ru,by" devel:libraries:c_c++.repo; \
        if [ "$TARGETPLATFORM" = "linux/amd64" ]; then \
            zypper -n addrepo --refresh "https://download.opensuse.org/repositories/devel:tools:building/openSUSE_Factory/?AVOID_COUNTRY=ru,by" devel:tools:building.repo; \
        elif [ "$TARGETPLATFORM" = "linux/arm64" ]; then \
            zypper -n addrepo --refresh "https://download.opensuse.org/repositories/devel:tools:building/openSUSE_Factory_ARM/?AVOID_COUNTRY=ru,by" devel:tools:building.repo; \
        else \
            exit 1; \
        fi && \
        zypper --gpg-auto-import-keys ref && break || sleep 1; \
    done

RUN zypper -n install cmake
RUN zypper -n install make

RUN for i in {1..10}; do \
        if [ "$TARGETPLATFORM" = "linux/amd64" ]; then \
            zypper -n addrepo --refresh "https://download.opensuse.org/repositories/graphics/openSUSE_Factory/?AVOID_COUNTRY=ru,by" graphics.repo; \
        elif [ "$TARGETPLATFORM" = "linux/arm64" ]; then \
            zypper -n addrepo --refresh "https://download.opensuse.org/repositories/graphics/openSUSE_Factory_ARM/?AVOID_COUNTRY=ru,by" graphics.repo; \
        else \
            exit 1; \
        fi && \
        zypper --gpg-auto-import-keys ref && break || sleep 1; \
    done

RUN zypper -n install graphviz-devel

RUN for i in {1..10}; do \
        if [ "$TARGETPLATFORM" = "linux/amd64" ]; then \
            zypper -n addrepo --refresh "https://download.opensuse.org/repositories/devel:gcc/openSUSE_Factory/?AVOID_COUNTRY=ru,by" devel:gcc.repo; \
        elif [ "$TARGETPLATFORM" = "linux/arm64" ]; then \
            zypper -n addrepo --refresh "https://download.opensuse.org/repositories/devel:gcc/openSUSE_Factory_ARM/?AVOID_COUNTRY=ru,by" devel:gcc.repo; \
        else \
            exit 1; \
        fi && \
        zypper --gpg-auto-import-keys ref && break || sleep 1; \
    done

RUN zypper -n install gcc14
RUN zypper -n install gcc14-c++

RUN for i in {1..10}; do \
        zypper -n addrepo --refresh "https://download.opensuse.org/repositories/devel:tools/16.0/?AVOID_COUNTRY=ru,by" devel:tools.repo && \
        zypper --gpg-auto-import-keys ref && break || sleep 1; \
    done
RUN zypper -n install doxygen

RUN for i in {1..10}; do \
        zypper -n addrepo --refresh "https://download.opensuse.org/repositories/network:NIS/16.0/?AVOID_COUNTRY=ru,by" network:NIS.repo && \
        zypper --gpg-auto-import-keys ref && break || sleep 1; \
    done
RUN zypper -n install libnsl-devel

RUN zypper clean -a

RUN git clone https://github.com/longhorn/dep-versions.git -b ${SRC_BRANCH} /usr/src/dep-versions && \
    cd /usr/src/dep-versions && \
    if [ -n "${SRC_TAG}" ] && git show-ref --tags ${SRC_TAG} > /dev/null 2>&1; then \
        echo "Checking out tag ${SRC_TAG}"; \
        cd /usr/src/dep-versions && git checkout tags/${SRC_TAG}; \
    fi

# Build nfs-ganesha
RUN export REPO_OVERRIDE="" && \
    export COMMIT_ID_OVERRIDE="" && \
    bash /usr/src/dep-versions/scripts/build-nfs-ganesha.sh "${REPO_OVERRIDE}" "${COMMIT_ID_OVERRIDE}"


FROM registry.suse.com/bci/bci-base:16.0 AS release

ARG TARGETPLATFORM
RUN if [ "$TARGETPLATFORM" != "linux/amd64" ] && [ "$TARGETPLATFORM" != "linux/arm64" ]; then \
    echo "Error: Unsupported TARGETPLATFORM: $TARGETPLATFORM" && \
    exit 1; \
    fi
ENV ARCH=${TARGETPLATFORM#linux/}

RUN zypper -n ref && \
    zypper update -y



RUN for i in {1..10}; do \
        zypper -n addrepo --refresh "https://download.opensuse.org/repositories/network:NIS/16.0/?AVOID_COUNTRY=ru,by" network:NIS.repo && \
        zypper -n addrepo --refresh "https://download.opensuse.org/repositories/graphics/openSUSE_Factory/?AVOID_COUNTRY=ru,by" graphics.repo && \
        zypper -n addrepo --refresh "https://download.opensuse.org/repositories/devel:tools/16.0/?AVOID_COUNTRY=ru,by" devel:tools.repo && \
        zypper -n addrepo --refresh "https://download.opensuse.org/repositories/devel:tools:building/16.0/?AVOID_COUNTRY=ru,by" devel:tools:building.repo && \
        zypper -n addrepo --refresh "https://download.opensuse.org/repositories/network:ldap/16.0/?AVOID_COUNTRY=ru,by" network:ldap.repo && \
        zypper --gpg-auto-import-keys ref && break || sleep 1; \
    done

# RUN microdnf install -y nano tar lsof e2fsprogs fuse-libs libss libblkid userspace-rcu dbus-x11 rpcbind hostname nfs-utils xfsprogs jemalloc libnfsidmap && microdnf clean all
RUN zypper -n install rpcbind hostname libblkid1 libjson-c* dbus-1-x11 dbus-1 nfsidmap-devel libnsl-devel nfs-kernel-server xfsprogs e2fsprogs

RUN for i in {1..10}; do \
        zypper -n addrepo --refresh "https://download.opensuse.org/repositories/Base:System/openSUSE_Factory/?AVOID_COUNTRY=ru,by" Base:System.repo && \
        zypper --gpg-auto-import-keys ref && break || sleep 1; \
    done

RUN zypper -n install nfs-client nfs4-acl-tools
RUN zypper clean -a

RUN mkdir -p /var/run/dbus && mkdir -p /export

# add libs from /usr/local/lib64
RUN echo /usr/local/lib64 > /etc/ld.so.conf.d/local_libs.conf

# do not ask systemd for user IDs or groups (slows down dbus-daemon start)
RUN [ -f /etc/nsswitch.conf ] && sed -i 's/systemd//g' /etc/nsswitch.conf || true

# ganesha reads /etc/mtab for mounted volumes
RUN ln -sf /proc/self/mounts /etc/mtab

COPY --from=lib_builder /usr/local /usr/local/
COPY --from=lib_builder /ganesha-extra /
COPY --from=app_builder /app/bin/longhorn-share-manager-${ARCH} /longhorn-share-manager

# run ldconfig after libs have been copied
RUN ldconfig

# only expose the nfsd since for v4 only that is necessary
EXPOSE 2049/tcp

ENTRYPOINT ["/longhorn-share-manager"]