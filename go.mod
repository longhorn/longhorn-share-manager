module github.com/longhorn/longhorn-share-manager

go 1.21

// install the below required grpc/protobuf versions
// https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-linux-x86_64.zip
// https://github.com/golang/protobuf.git # git checkout v1.3.2 # cd protoc-gen-go # go install

replace (
<<<<<<< HEAD
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
	github.com/golang/protobuf v1.5.4
	github.com/google/fscrypt v0.3.4
	github.com/longhorn/go-iscsi-helper v0.0.0-20220805034259-7b59e22574bb
=======
	k8s.io/api => k8s.io/api v0.29.2
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.29.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.29.2
	k8s.io/apiserver => k8s.io/apiserver v0.29.2
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.29.2
	k8s.io/client-go => k8s.io/client-go v0.29.2
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.29.2
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.29.2
	k8s.io/code-generator => k8s.io/code-generator v0.29.2
	k8s.io/component-base => k8s.io/component-base v0.29.2
	k8s.io/component-helpers => k8s.io/component-helpers v0.29.2
	k8s.io/controller-manager => k8s.io/controller-manager v0.29.2
	k8s.io/cri-api => k8s.io/cri-api v0.29.2
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.29.2
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.29.2
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.29.2
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.29.2
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.29.2
	k8s.io/kubectl => k8s.io/kubectl v0.29.2
	k8s.io/kubelet => k8s.io/kubelet v0.29.2
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.29.2
	k8s.io/metrics => k8s.io/metrics v0.29.2
	k8s.io/mount-utils => k8s.io/mount-utils v0.29.2
	k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.29.2
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.29.2
)

require (
	github.com/google/fscrypt v0.3.4
	github.com/longhorn/go-common-libs v0.0.0-20240307063052-6e77996eda29
>>>>>>> 6e921c6 (refactor(utils): move is mount read only function to common lib)
	github.com/mitchellh/go-ps v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.3
	github.com/urfave/cli v1.22.14
<<<<<<< HEAD
	golang.org/x/net v0.17.0
	golang.org/x/sys v0.13.0
	google.golang.org/grpc v1.54.0
	k8s.io/kubernetes v1.28.2
	k8s.io/mount-utils v0.28.2
=======
	golang.org/x/net v0.22.0
	golang.org/x/sys v0.18.0
	google.golang.org/grpc v1.62.1
	google.golang.org/protobuf v1.33.0
	k8s.io/kubernetes v1.29.2
	k8s.io/mount-utils v0.29.2
>>>>>>> 6e921c6 (refactor(utils): move is mount read only function to common lib)
	k8s.io/utils v0.0.0-20240102154912-e7106e64919e
)

require (
	github.com/c9s/goprocinfo v0.0.0-20170724085704-0010a05ce49f // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
<<<<<<< HEAD
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/opencontainers/selinux v1.10.0 // indirect
=======
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-logr/logr v1.3.0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/moby/sys/mountinfo v0.6.2 // indirect
	github.com/opencontainers/selinux v1.11.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
>>>>>>> 6e921c6 (refactor(utils): move is mount read only function to common lib)
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230525234030-28d5490b6b19 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	k8s.io/klog/v2 v2.100.1 // indirect
)
