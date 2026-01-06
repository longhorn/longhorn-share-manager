module github.com/longhorn/longhorn-share-manager

go 1.24.0

toolchain go1.25.5

// Replace directives are required for dependencies in this section because:
// - This module imports k8s.io/kubernetes.
// - The development for all of these dependencies is done at kubernetes/staging and then synced to other repos.
// - The go.mod file for k8s.io/kubernetes imports these dependencies with version v0.0.0 (which does not exist) and \
//   uses its own replace directives to load the appropriate code from kubernetes/staging.
// - Go is not able to find a version v0.0.0 for these dependencies and cannot meaningfully follow replace directives in
//   another go.mod file.
//
// The solution (which is used by all projects that import k8s.io/kubernetes) is to add replace directives for all
// k8s.io dependencies of k8s.io/kubernetes that k8s.io/kubernetes itself replaces in its go.mod file. The replace
// directives should pin the version of each dependency to the version of k8s.io/kubernetes that is imported. For
// example, if we import k8s.io/kubernetes v1.28.5, we should use v0.28.5 of all the replace directives. Depending on
// the portions of k8s.io/kubernetes code this module actually uses, not all of the replace directives may strictly be
// necessary. However, it is better to include all of them for consistency.
replace (
	k8s.io/api => k8s.io/api v0.34.3
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.34.3
	k8s.io/apimachinery => k8s.io/apimachinery v0.34.3
	k8s.io/apiserver => k8s.io/apiserver v0.34.3
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.34.3
	k8s.io/client-go => k8s.io/client-go v0.34.3
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.34.3
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.34.3
	k8s.io/code-generator => k8s.io/code-generator v0.34.3
	k8s.io/component-base => k8s.io/component-base v0.34.3
	k8s.io/component-helpers => k8s.io/component-helpers v0.34.3
	k8s.io/controller-manager => k8s.io/controller-manager v0.34.3
	k8s.io/cri-api => k8s.io/cri-api v0.34.3
	k8s.io/cri-client => k8s.io/cri-client v0.34.3
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.34.3
	k8s.io/dynamic-resource-allocation => k8s.io/dynamic-resource-allocation v0.34.3
	k8s.io/endpointslice => k8s.io/endpointslice v0.34.3
	k8s.io/kms => k8s.io/kms v0.34.3
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.34.3
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.34.3
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.34.3
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.34.3
	k8s.io/kubectl => k8s.io/kubectl v0.34.3
	k8s.io/kubelet => k8s.io/kubelet v0.34.3
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.30.14
	k8s.io/metrics => k8s.io/metrics v0.34.3
	k8s.io/mount-utils => k8s.io/mount-utils v0.34.3
	k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.34.3
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.34.3
)

require (
	github.com/google/fscrypt v0.3.6
	github.com/longhorn/go-common-libs v0.0.0-20260103034008-119bdcf1b2d6
	github.com/longhorn/types v0.0.0-20251228142423-336840fb2fd6
	github.com/mitchellh/go-ps v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.3
	github.com/urfave/cli v1.22.17
	golang.org/x/net v0.47.0
	golang.org/x/sys v0.38.0
	google.golang.org/grpc v1.78.0
	google.golang.org/protobuf v1.36.11
	k8s.io/api v0.34.3
	k8s.io/apimachinery v0.34.3
	k8s.io/client-go v0.34.3
	k8s.io/kubernetes v1.34.3
	k8s.io/mount-utils v0.34.3
	k8s.io/utils v0.0.0-20260106112306-0fe9cd71b2f8
)

require (
	github.com/c9s/goprocinfo v0.0.0-20210130143923-c95fcf8c64a8 // indirect
	github.com/cockroachdb/errors v1.12.0 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/redact v1.1.5 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/emicklei/go-restful/v3 v3.12.2 // indirect
	github.com/fxamacker/cbor/v2 v2.9.0 // indirect
	github.com/getsentry/sentry-go v0.27.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/gnostic-models v0.7.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/moby/sys/mountinfo v0.7.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.3-0.20250322232337-35a7c28c31ee // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/opencontainers/selinux v1.11.1 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/shirou/gopsutil/v3 v3.24.5 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/exp v0.0.0-20251219203646-944ab1f22d93 // indirect
	golang.org/x/oauth2 v0.32.0 // indirect
	golang.org/x/term v0.37.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	golang.org/x/time v0.9.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251029180050-ab9386a59fda // indirect
	gopkg.in/evanphx/json-patch.v4 v4.12.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/kube-openapi v0.0.0-20250710124328-f3f2b991d03b // indirect
	sigs.k8s.io/json v0.0.0-20241014173422-cfa47c3a1cc8 // indirect
	sigs.k8s.io/randfill v1.0.0 // indirect
	sigs.k8s.io/structured-merge-diff/v6 v6.3.0 // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
)
