module github.com/longhorn/longhorn-share-manager

go 1.13

// install the below required grpc/protobuf versions
// https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-linux-x86_64.zip
// https://github.com/golang/protobuf.git # git checkout v1.3.2 # cd protoc-gen-go # go install

replace (
	k8s.io/api => k8s.io/api v0.16.15
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.16.15
	k8s.io/apimachinery => k8s.io/apimachinery v0.16.15
	k8s.io/apiserver => k8s.io/apiserver v0.16.15
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.16.15
	k8s.io/client-go => k8s.io/client-go v0.16.15
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.16.15
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.16.15
	k8s.io/code-generator => k8s.io/code-generator v0.16.16-rc.0
	k8s.io/component-base => k8s.io/component-base v0.16.15
	k8s.io/cri-api => k8s.io/cri-api v0.16.16-rc.0
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.16.15
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.16.15
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.16.15
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.16.15
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.16.15
	k8s.io/kubectl => k8s.io/kubectl v0.16.15
	k8s.io/kubelet => k8s.io/kubelet v0.16.15
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.16.15
	k8s.io/metrics => k8s.io/metrics v0.16.15
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.16.15
)

require (
	github.com/golang/protobuf v1.3.3-0.20190920234318-1680a479a2cf
	github.com/guelfey/go.dbus v0.0.0-20131113121618-f6a3a2366cc3
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.4.2
	github.com/urfave/cli v1.22.1
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	google.golang.org/grpc v1.23.0 // pinned
	k8s.io/kubernetes v1.16.15
)
