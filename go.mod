module github.com/longhorn/longhorn-share-manager

go 1.22.0

toolchain go1.22.2

replace (
	k8s.io/api => k8s.io/api v0.29.4
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.29.4
	k8s.io/apimachinery => k8s.io/apimachinery v0.29.4
	k8s.io/apiserver => k8s.io/apiserver v0.29.4
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.29.4
	k8s.io/client-go => k8s.io/client-go v0.29.4
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.29.4
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.29.4
	k8s.io/code-generator => k8s.io/code-generator v0.29.4
	k8s.io/component-base => k8s.io/component-base v0.29.4
	k8s.io/component-helpers => k8s.io/component-helpers v0.29.4
	k8s.io/controller-manager => k8s.io/controller-manager v0.29.4
	k8s.io/cri-api => k8s.io/cri-api v0.29.4
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.29.4
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.29.4
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.29.4
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.29.4
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.29.4
	k8s.io/kubectl => k8s.io/kubectl v0.29.4
	k8s.io/kubelet => k8s.io/kubelet v0.29.4
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.29.4
	k8s.io/metrics => k8s.io/metrics v0.29.4
	k8s.io/mount-utils => k8s.io/mount-utils v0.29.4
	k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.29.4
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.29.4
)

require (
	github.com/google/fscrypt v0.3.5
	github.com/longhorn/go-common-libs v0.0.0-20240420070800-82cf6b3fac64
	github.com/longhorn/types v0.0.0-20240417112740-a0d8514936b8
	github.com/mitchellh/go-ps v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.3
	github.com/urfave/cli v1.22.14
	golang.org/x/net v0.23.0
	golang.org/x/sys v0.19.0
	google.golang.org/grpc v1.63.2
	google.golang.org/protobuf v1.33.0
	k8s.io/kubernetes v1.29.4
	k8s.io/mount-utils v0.30.0
	k8s.io/utils v0.0.0-20240310230437-4693a0247e57
)

require (
	github.com/c9s/goprocinfo v0.0.0-20210130143923-c95fcf8c64a8 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/moby/sys/mountinfo v0.6.2 // indirect
	github.com/opencontainers/selinux v1.11.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/shirou/gopsutil/v3 v3.24.3 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240227224415-6ceb2ff114de // indirect
	k8s.io/apimachinery v0.0.0 // indirect
	k8s.io/klog/v2 v2.120.1 // indirect
)
