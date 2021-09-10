module github.com/longhorn/longhorn-share-manager

go 1.13

// install the below required grpc/protobuf versions
// https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-linux-x86_64.zip
// https://github.com/golang/protobuf.git # git checkout v1.3.2 # cd protoc-gen-go # go install

replace (
	k8s.io/api => k8s.io/api v0.18.19
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.18.19
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.19
	k8s.io/apiserver => k8s.io/apiserver v0.18.19
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.18.19
	k8s.io/client-go => k8s.io/client-go v0.18.19
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.18.19
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.18.19
	k8s.io/code-generator => k8s.io/code-generator v0.18.19
	k8s.io/component-base => k8s.io/component-base v0.18.19
	k8s.io/cri-api => k8s.io/cri-api v0.18.19
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.18.19
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.18.19
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.18.19
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.18.19
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.18.19
	k8s.io/kubectl => k8s.io/kubectl v0.18.19
	k8s.io/kubelet => k8s.io/kubelet v0.18.19
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.18.19
	k8s.io/metrics => k8s.io/metrics v0.18.19
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.18.19
)

require (
	github.com/longhorn/go-iscsi-helper v0.0.0-20210330030558-49a327fb024e
	github.com/sirupsen/logrus v1.4.2
	github.com/urfave/cli v1.22.1
	golang.org/x/net v0.0.0-20201201195509-5d6afe98e0b7 // indirect
	k8s.io/api v0.18.19 // indirect
	k8s.io/apiextensions-apiserver v0.0.0 // indirect
	k8s.io/apimachinery v0.18.19 // indirect
	k8s.io/client-go v0.18.19 // indirect
	k8s.io/kubernetes v1.18.19
	k8s.io/metrics v0.18.19 // indirect
	k8s.io/utils v0.0.0-20200324210504-a9aa75ae1b89
	sigs.k8s.io/yaml v1.2.0 // indirect
)
