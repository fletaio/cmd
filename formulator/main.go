package main

import (
	"encoding/hex"

	"git.fleta.io/fleta/core/data"
	"git.fleta.io/fleta/core/formulator"
	"git.fleta.io/fleta/core/kernel"
	"git.fleta.io/fleta/core/key"
	"git.fleta.io/fleta/core/reward"
	"git.fleta.io/fleta/framework/config"
	"git.fleta.io/fleta/framework/peer"
	"git.fleta.io/fleta/framework/router"
	"git.fleta.io/fleta/framework/router/evilnode"

	"git.fleta.io/fleta/common"
)

// Config is a configuration for the cmd
type Config struct {
	SeedNodes      []string
	ObserverKeyMap map[string]string
	KeyHex         string
	Formulator     string
	Port           int
	StoreRoot      string
}

func main() {
	var cfg Config
	if err := config.LoadFile("./config.toml", &cfg); err != nil {
		panic(err)
	}
	if len(cfg.StoreRoot) == 0 {
		cfg.StoreRoot = "./formulator"
	}

	var frkey key.Key
	if bs, err := hex.DecodeString(cfg.KeyHex); err != nil {
		panic(err)
	} else if Key, err := key.NewMemoryKeyFromBytes(bs); err != nil {
		panic(err)
	} else {
		frkey = Key
	}

	ObserverKeyMap := map[common.PublicHash]string{}
	ObserverKeyBoolMap := map[common.PublicHash]bool{}
	for k, netAddr := range cfg.ObserverKeyMap {
		pubhash, err := common.ParsePublicHash(k)
		if err != nil {
			panic(err)
		}
		ObserverKeyMap[pubhash] = netAddr
		ObserverKeyBoolMap[pubhash] = true
	}

	GenCoord := common.NewCoordinate(0, 0)
	act := data.NewAccounter(GenCoord)
	tran := data.NewTransactor(GenCoord)
	if err := initChainComponent(act, tran); err != nil {
		panic(err)
	}
	GenesisContextData, err := initGenesisContextData(act, tran)
	if err != nil {
		panic(err)
	}

	ks, err := kernel.NewStore(cfg.StoreRoot+"/kernel", 1, act, tran)
	if err != nil {
		panic(err)
	}

	rd := &reward.TestNetRewarder{}
	kn, err := kernel.NewKernel(&kernel.Config{
		ChainCoord:     GenCoord,
		ObserverKeyMap: ObserverKeyBoolMap,
	}, ks, rd, GenesisContextData)
	if err != nil {
		panic(err)
	}

	frcfg := &formulator.Config{
		ChainCoord:     GenCoord,
		Key:            frkey,
		SeedNodes:      cfg.SeedNodes,
		ObserverKeyMap: ObserverKeyMap,
		Formulator:     common.MustParseAddress(cfg.Formulator),
		Router: router.Config{
			Network: "tcp",
			Port:    cfg.Port,
			EvilNodeConfig: evilnode.Config{
				StorePath: cfg.StoreRoot + "/router",
			},
		},
		Peer: peer.Config{
			StorePath: cfg.StoreRoot + "/peers",
		},
	}
	fr, err := formulator.NewFormulator(frcfg, kn)
	if err != nil {
		panic(err)
	}
	fr.Run()
}
