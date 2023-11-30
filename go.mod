module github.com/longhorn/longhorn-share-manager

go 1.21

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
	github.com/google/fscrypt v0.3.4
	github.com/longhorn/go-iscsi-helper v0.0.0-20220805034259-7b59e22574bb
	github.com/mitchellh/go-ps v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.0
	github.com/urfave/cli v1.22.12
	golang.org/x/net v0.17.0
	golang.org/x/sys v0.13.0
	google.golang.org/grpc v1.58.3
	google.golang.org/protobuf v1.31.0
	k8s.io/kubernetes v1.28.2
	k8s.io/mount-utils v0.28.2
	k8s.io/utils v0.0.0-20230406110748-d93618cff8a2
)

require (
	github.com/c9s/goprocinfo v0.0.0-20170724085704-0010a05ce49f // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/opencontainers/selinux v1.10.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230711160842-782d3b101e98 // indirect
	k8s.io/klog/v2 v2.100.1 // indirect
)
