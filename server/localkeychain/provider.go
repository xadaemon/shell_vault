package localkeychain

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/spf13/afero"
	"github.com/stateprism/libprisma/memkv"
	"github.com/stateprism/prisma_ca/server/providers"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"time"
)

// LocalKeychain must implement providers.KeychainProvider
type LocalKeychain struct {
	fs           afero.Fs
	expHook      providers.ExpKeyHook
	logger       *zap.Logger
	store        *memkv.MemKV
	activeKey    *providers.PrivateKey
	ticker       *time.Ticker
	tickInterval time.Duration
	tickStop     chan struct{}
	isLocal      bool
	localPath    string
}

type LKParams struct {
	fx.In
	Lc     fx.Lifecycle
	Config providers.ConfigurationProvider
	Logger *zap.Logger
}

type onDiskKey struct {
	Name string
	Ttl  uint64
	Key  crypto.PrivateKey
}

func NewLocalKeychain(par LKParams) (providers.KeychainProvider, error) {
	kcPath, err := par.Config.GetString("providers.local_keychain_provider.path")
	if err != nil {
		return nil, err
	}
	kcFsType, err := par.Config.GetString("providers.local_keychain_provider.fs")
	if err != nil {
		kcFsType = "local"
	}

	var fs afero.Fs
	switch kcFsType {
	case "local", "":
		fs = afero.NewOsFs()
		err := os.Chdir(par.Config.GetLocalStore())
		if err != nil {
			return nil, err
		}
		kcPath, err = filepath.Abs(kcPath)
	case "memory":
		fs = afero.NewMemMapFs()
		_ = fs.MkdirAll(kcPath, 0700)
	default:
		return nil, fmt.Errorf("unknown fs type: %s for provider", kcFsType)
	}

	if err != nil {
		return nil, err
	}

	stat, err := fs.Stat(kcPath)
	if os.IsNotExist(err) {
		errDir := fs.MkdirAll(kcPath, 0755)
		if errDir != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("path %s is a file, this provider requires a directory", kcPath)
	}

	store := memkv.NewMemKV("/", &memkv.Opts{CaseInsensitive: true})

	lk := &LocalKeychain{
		fs:     fs,
		store:  store,
		logger: par.Logger,
	}

	// Start the ttl ticker, we check for expired every minute
	if tr, err := par.Config.GetInt64("providers.local_keychain_provider.ttl_tick"); err != nil {
		lk.logger.Warn("TTL tick is not configured setting to 60s default")
		lk.tickInterval = 60 * time.Second
	} else {
		lk.logger.Info("TTL tick will be configured to", zap.Duration("tickRate", time.Duration(tr)*time.Second))
		lk.tickInterval = time.Duration(tr) * time.Second
	}

	par.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lk.ticker = time.NewTicker(lk.tickInterval)
			lk.tickStop = make(chan struct{})
			go ttlTick(lk)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			close(lk.tickStop)
			return nil
		},
	})

	return lk, nil
}

func ttlTick(l *LocalKeychain) {
	for {
		select {
		case <-l.ticker.C:
			now := time.Now()
			l.logger.Debug(
				"Checking ttl of stored keys:",
				zap.Time("event_at", now),
				zap.Time("next_check", now.Add(60*time.Second)),
			)
			if l.expHook == nil {
				l.logger.Warn("No hook set on ttl provider to notify ca")
			}
			l.expHook(l, nil)
		case <-l.tickStop:
			l.logger.Info("stopping ttl ticks")
			l.ticker.Stop()
			return
		}
	}
}

func (l *LocalKeychain) Unseal(key []byte) bool {
	return true
}

func (l *LocalKeychain) Seal() bool {
	return true
}

func makeNewEcdsaKey(kt providers.KeyType) (ecdsa.PrivateKey, error) {
	var c elliptic.Curve
	switch kt {
	case providers.KEY_TYPE_ECDSA_256:
		c = elliptic.P256()
	case providers.KEY_TYPE_ECDSA_384:
		c = elliptic.P384()
	case providers.KEY_TYPE_ECDSA_521:
		c = elliptic.P521()
	default:
		panic("unhandled default case")
	}
	t, err := ecdsa.GenerateKey(c, cryptorand.Reader)
	if err != nil {
		return ecdsa.PrivateKey{}, err
	}
	return *t, nil
}

func makeRsaKey(kt providers.KeyType) (rsa.PrivateKey, error) {
	var t *rsa.PrivateKey
	var bits int
	switch kt {
	case providers.KEY_TYPE_RSA_2048:
		bits = 2048
	case providers.KEY_TYPE_RSA_4096:
		bits = 4096
	default:
		panic("unhandled default case")
	}
	t, err := rsa.GenerateKey(cryptorand.Reader, bits)
	if err != nil {
		return rsa.PrivateKey{}, err
	}
	return *t, nil
}

func (l *LocalKeychain) MakeNewKey(keyName providers.KeyIdentifier, kt providers.KeyType, ttl int64) (providers.KeyIdentifier, error) {
	if _, ok := keyName.(string); !ok || keyName == nil {
		return nil, errors.New("this provider only takes string key names")
	}
	// allow to rotate the rootKey if needed, but
	if l.store.Contains(keyName.(string)) {
		return nil, errors.New("this key name is already contained in this keyring")
	}
	var key crypto.PrivateKey
	var err error
	switch kt {
	case providers.KEY_TYPE_ED25519:
		var t ed25519.PrivateKey
		_, t, err = ed25519.GenerateKey(cryptorand.Reader)
		key = t
	case providers.KEY_TYPE_ECDSA_256, providers.KEY_TYPE_ECDSA_384, providers.KEY_TYPE_ECDSA_521:
		key, err = makeNewEcdsaKey(kt)
	case providers.KEY_TYPE_RSA_2048, providers.KEY_TYPE_RSA_4096:
		key, err = makeRsaKey(kt)
	default:
		return nil, fmt.Errorf("invalid key format: %s", kt)
	}
	if err != nil {
		return nil, err
	}

	pk := providers.NewPrivateKey(keyName, kt, key, time.Duration(ttl))

	l.store.Set(keyName.(string), pk)

	return keyName, nil
}

func (l *LocalKeychain) SetActiveKey(kid providers.KeyIdentifier) bool {
	_, ok := kid.(string)
	if !ok {
		return false
	}
	key, ok := l.store.Get(kid.(string))
	if !ok {
		return false
	}

	l.activeKey, _ = key.(*providers.PrivateKey)
	return true
}

func (l *LocalKeychain) SetExpKeyHook(f providers.ExpKeyHook) providers.ExpKeyHook {
	old := l.expHook
	l.expHook = f
	return old
}

func (l *LocalKeychain) LookupKey(criteria providers.KeyLookupCriteria) (providers.KeyIdentifier, bool) {
	//TODO implement me
	panic("implement me")
}

func (l *LocalKeychain) RetrieveKey(kid providers.KeyIdentifier) (*providers.PrivateKey, bool) {
	_, ok := kid.(string)
	if !ok {
		return nil, false
	}
	key, ok := l.store.Get(kid.(string))
	if !ok {
		return nil, false
	}

	return key.(*providers.PrivateKey), true
}

func (l *LocalKeychain) DropKey(keyName providers.KeyIdentifier) bool {
	if _, ok := keyName.(string); !ok {
		return false
	}
	return l.store.Drop(keyName.(string), false)
}
func (l *LocalKeychain) MakeAndReplaceKey(keyName providers.KeyIdentifier, kt providers.KeyType, ttl int64) (providers.KeyIdentifier, error) {
	if _, ok := keyName.(string); !ok || keyName == nil {
		return nil, errors.New("this provider only takes string key names")
	}
	l.DropKey(keyName.(string))

	return l.MakeNewKey(keyName, kt, ttl)
}

func (l *LocalKeychain) MakeNewKeyIfNotExists(keyName providers.KeyIdentifier, kt providers.KeyType, ttl int64) (providers.KeyIdentifier, error) {
	if _, ok := keyName.(string); !ok || keyName == nil {
		return nil, errors.New("this provider only takes string key names")
	}
	if _, ok := l.RetrieveKey(keyName); ok {
		return keyName, nil
	}
	return l.MakeNewKey(keyName, kt, ttl)
}

func (l *LocalKeychain) String() string {
	return "LocalKeychain"
}

func (l *LocalKeychain) saveKey(keyName string) {

}
