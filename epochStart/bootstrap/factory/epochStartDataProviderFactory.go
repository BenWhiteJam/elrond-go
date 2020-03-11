package factory

import (
	"time"

	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-go/epochStart/bootstrap"
	"github.com/ElrondNetwork/elrond-go/hashing"
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/ElrondNetwork/elrond-go/p2p"
	"github.com/ElrondNetwork/elrond-go/sharding"
)

type epochStartDataProviderFactory struct {
	pubKey              crypto.PublicKey
	messenger           p2p.Messenger
	marshalizer         marshal.Marshalizer
	hasher              hashing.Hasher
	nodesConfigProvider bootstrap.NodesConfigProviderHandler
	shouldSync          bool
}

// EpochStartDataProviderFactoryArgs holds the arguments needed for creating aa factory for the epoch start data
// provider component
type EpochStartDataProviderFactoryArgs struct {
	PubKey                crypto.PublicKey
	Messenger             p2p.Messenger
	Marshalizer           marshal.Marshalizer
	Hasher                hashing.Hasher
	NodesConfigProvider   bootstrap.NodesConfigProviderHandler
	StartTime             time.Time
	OriginalNodesConfig   *sharding.NodesSetup
	GeneralConfig         *config.Config
	IsEpochFoundInStorage bool
}

// NewEpochStartDataProviderFactory returns a new instance of epochStartDataProviderFactory
func NewEpochStartDataProviderFactory(args EpochStartDataProviderFactoryArgs) (*epochStartDataProviderFactory, error) {
	if check.IfNil(args.PubKey) {
		return nil, bootstrap.ErrNilPublicKey
	}
	if check.IfNil(args.Messenger) {
		return nil, bootstrap.ErrNilMessenger
	}
	if check.IfNil(args.Marshalizer) {
		return nil, bootstrap.ErrNilMarshalizer
	}
	if check.IfNil(args.Hasher) {
		return nil, bootstrap.ErrNilHasher
	}
	if check.IfNil(args.NodesConfigProvider) {
		return nil, bootstrap.ErrNilNodesConfigProvider
	}

	shouldSync := bootstrap.ShouldSyncWithTheNetwork(
		args.StartTime,
		args.IsEpochFoundInStorage,
		args.OriginalNodesConfig,
		args.GeneralConfig,
	)
	shouldSync = true // harcoded so we can test we can sync

	return &epochStartDataProviderFactory{
		pubKey:              args.PubKey,
		messenger:           args.Messenger,
		marshalizer:         args.Marshalizer,
		hasher:              args.Hasher,
		nodesConfigProvider: args.NodesConfigProvider,
		shouldSync:          shouldSync,
	}, nil
}

// Create will init and return an instance of an epoch start data provider
func (esdpf *epochStartDataProviderFactory) Create() (bootstrap.EpochStartDataProviderHandler, error) {
	if !esdpf.shouldSync {
		return &disabledEpochStartDataProvider{}, nil
	}

	epochStartMetaBlockInterceptor, err := bootstrap.NewSimpleEpochStartMetaBlockInterceptor(esdpf.marshalizer, esdpf.hasher)
	if err != nil {
		return nil, err
	}
	metaBlockInterceptor, err := bootstrap.NewSimpleMetaBlockInterceptor(esdpf.marshalizer, esdpf.hasher)
	if err != nil {
		return nil, err
	}
	shardHdrInterceptor, err := bootstrap.NewSimpleShardHeaderInterceptor(esdpf.marshalizer, esdpf.hasher)
	if err != nil {
		return nil, err
	}
	miniBlockInterceptor, err := bootstrap.NewSimpleMiniBlockInterceptor(esdpf.marshalizer, esdpf.hasher)
	if err != nil {
		return nil, err
	}

	argsEpochStart := bootstrap.ArgsEpochStartDataProvider{
		PublicKey:                      esdpf.pubKey,
		Messenger:                      esdpf.messenger,
		Marshalizer:                    esdpf.marshalizer,
		Hasher:                         esdpf.hasher,
		NodesConfigProvider:            esdpf.nodesConfigProvider,
		EpochStartMetaBlockInterceptor: epochStartMetaBlockInterceptor,
		MetaBlockInterceptor:           metaBlockInterceptor,
		ShardHeaderInterceptor:         shardHdrInterceptor,
		MiniBlockInterceptor:           miniBlockInterceptor,
	}
	epochStartDataProvider, err := bootstrap.NewEpochStartDataProvider(argsEpochStart)

	return epochStartDataProvider, nil
}