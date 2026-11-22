package volume

import (
	"context"
	"fmt"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/longhorn/go-common-libs/crypto/keystore"
	lhtypes "github.com/longhorn/go-common-libs/types"

	"github.com/longhorn/longhorn-share-manager/pkg/types"
)

// secretReadTimeout keeps a stalled API server from holding the server mutex.
const secretReadTimeout = 30 * time.Second

// readEncryptionSecret reads a Secret from the pod's namespace with the pod's
// ServiceAccount, which already reads Secrets. The client is built per call so
// a pod started before its token was mounted recovers.
func readEncryptionSecret(ctx context.Context, name string) (map[string]string, error) {
	namespace := os.Getenv(types.EnvPodNamespace)
	if namespace == "" {
		return nil, fmt.Errorf("cannot resolve an encryption secret: %v is not set", types.EnvPodNamespace)
	}
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, secretReadTimeout)
	defer cancel()
	secret, err := client.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return keystore.SecretData(secret.Data), nil
}

// NewKeyProvider builds the volume's key store: the secret provider from the
// passphrase the pod was given, any other from the Secret the volume names.
func (v Volume) NewKeyProvider(ctx context.Context) (keystore.KeyProvider, error) {
	if v.KeyProviderName() == lhtypes.EncryptionKeyProviderSecret {
		return keystore.NewProvider(map[string]string{lhtypes.CryptoKeyValue: v.Passphrase})
	}
	if v.EncryptionSecretName == "" {
		return nil, fmt.Errorf("volume %v names no encryption secret", v.Name)
	}
	secret, err := readEncryptionSecret(ctx, v.EncryptionSecretName)
	if err != nil {
		return nil, fmt.Errorf("failed to read the encryption secret %v of volume %v: %v", v.EncryptionSecretName, v.Name, err)
	}
	return keystore.NewProvider(secret)
}
