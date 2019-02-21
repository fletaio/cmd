package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"git.fleta.io/fleta/core/node"
	"git.fleta.io/fleta/core/reward"
	"git.fleta.io/fleta/framework/config"
	"git.fleta.io/fleta/framework/peer"
	"git.fleta.io/fleta/framework/router"
	"git.fleta.io/fleta/framework/router/evilnode"
	"github.com/dgraph-io/badger"

	"git.fleta.io/fleta/core/data"
	"git.fleta.io/fleta/core/kernel"

	"git.fleta.io/fleta/common"
)

// Config is a configuration for the cmd
type Config struct {
	SeedNodes    []string
	ObserverKeys []string
	Port         int
	StoreRoot    string
	ForceRecover bool
}

func main() {
	var cfg Config
	if err := config.LoadFile("./config.toml", &cfg); err != nil {
		panic(err)
	}
	if len(cfg.StoreRoot) == 0 {
		cfg.StoreRoot = "./data"
	}

	ObserverKeyMap := map[common.PublicHash]bool{}
	for _, k := range cfg.ObserverKeys {
		pubhash, err := common.ParsePublicHash(k)
		if err != nil {
			panic(err)
		}
		ObserverKeyMap[pubhash] = true
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

	var closable Closable
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigc
		if closable != nil {
			closable.Close()
		}
	}()

	var ks *kernel.Store
	if s, err := kernel.NewStore(cfg.StoreRoot+"/kernel", BlockchainVersion, act, tran, cfg.ForceRecover); err != nil {
		if cfg.ForceRecover || err != badger.ErrTruncateNeeded {
			panic(err)
		} else {
			fmt.Println(err)
			fmt.Println("Do you want to recover database(it can be failed)? [y/n]")
			var answer string
			fmt.Scanf("%s", &answer)
			if strings.ToLower(answer) == "y" {
				if s, err := kernel.NewStore(cfg.StoreRoot+"/kernel", BlockchainVersion, act, tran, true); err != nil {
					panic(err)
				} else {
					ks = s
				}
			} else {
				os.Exit(1)
			}
		}
	} else {
		ks = s
	}
	closable = ks

	rd := &reward.TestNetRewarder{}
	kn, err := kernel.NewKernel(&kernel.Config{
		ChainCoord:     GenCoord,
		ObserverKeyMap: ObserverKeyMap,
	}, ks, rd, GenesisContextData)
	if err != nil {
		panic(err)
	}
	closable = kn

	ndcfg := &node.Config{
		ChainCoord: GenCoord,
		SeedNodes:  cfg.SeedNodes,
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
	nd, err := node.NewNode(ndcfg, kn)
	if err != nil {
		panic(err)
	}
	closable = nd
	nd.Run()
}

type Closable interface {
	Close()
}
