module github.com/longhorn/longhorn-share-manager

go 1.20

// install the below required grpc/protobuf versions
// https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-linux-x86_64.zip
// https://github.com/golang/protobuf.git # git checkout v1.3.2 # cd protoc-gen-go # go install

replace (
	k8s.io/api => k8s.io/api v0.23.6
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.23.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.23.6
	k8s.io/apiserver => k8s.io/apiserver v0.23.6
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.23.6
	k8s.io/client-go => k8s.io/client-go v0.23.6
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.23.6
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.23.6
	k8s.io/code-generator => k8s.io/code-generator v0.23.6
	k8s.io/component-base => k8s.io/component-base v0.23.6
	k8s.io/component-helpers => k8s.io/component-helpers v0.23.6
	k8s.io/controller-manager => k8s.io/controller-manager v0.23.6
	k8s.io/cri-api => k8s.io/cri-api v0.23.6
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.23.6
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.23.6
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.23.6
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.23.6
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.23.6
	k8s.io/kubectl => k8s.io/kubectl v0.23.6
	k8s.io/kubelet => k8s.io/kubelet v0.23.6
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.23.6
	k8s.io/metrics => k8s.io/metrics v0.23.6
	k8s.io/mount-utils => k8s.io/mount-utils v0.23.6
	k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.23.6
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.23.6

)

require (
	github.com/golang/protobuf v1.5.3
	github.com/google/fscrypt v0.3.4
	github.com/longhorn/go-iscsi-helper v0.0.0-20220805034259-7b59e22574bb
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.0
	github.com/urfave/cli v1.22.12
	golang.org/x/net v0.0.0-20211209124913-491a49abca63
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8
	google.golang.org/grpc v1.40.0
	k8s.io/kubernetes v1.23.6
	k8s.io/mount-utils v0.23.6
	k8s.io/utils v0.0.0-20211116205334-6203023598ed
)

require (
	github.com/bits-and-blooms/bitset v1.2.0 // indirect
	github.com/c9s/goprocinfo v0.0.0-20170724085704-0010a05ce49f // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/go-logr/logr v1.2.0 // indirect
	github.com/opencontainers/selinux v1.8.2 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20210831024726-fe130286e0e2 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	k8s.io/klog/v2 v2.30.0 // indirect
)
