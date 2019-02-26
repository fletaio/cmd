package main

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"git.fleta.io/fleta/cmd/closer"
	"git.fleta.io/fleta/core/block"
	"git.fleta.io/fleta/core/data"
	"git.fleta.io/fleta/core/formulator"
	"git.fleta.io/fleta/core/kernel"
	"git.fleta.io/fleta/core/key"
	"git.fleta.io/fleta/core/reward"
	"git.fleta.io/fleta/framework/config"
	"git.fleta.io/fleta/framework/peer"
	"git.fleta.io/fleta/framework/router"
	"git.fleta.io/fleta/framework/router/evilnode"
	"git.fleta.io/fleta/framework/rpc"
	"github.com/dgraph-io/badger"

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
	ForceRecover   bool
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

	cm := closer.NewManager()
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigc
		cm.CloseAll()
	}()
	defer cm.CloseAll()

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
	cm.Add("kernel.Store", ks)

	rd := &reward.TestNetRewarder{}
	kn, err := kernel.NewKernel(&kernel.Config{
		ChainCoord:     GenCoord,
		ObserverKeyMap: ObserverKeyBoolMap,
	}, ks, rd, GenesisContextData)
	if err != nil {
		panic(err)
	}
	cm.RemoveAll()
	cm.Add("kernel.Kernel", kn)

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
	cm.RemoveAll()
	cm.Add("cmd.Formulator", fr)

	go fr.Run()

	rm := rpc.NewManager()
	cm.RemoveAll()
	cm.Add("rpc.Manager", rm)
	cm.Add("cmd.Formulator", fr)
	kn.AddEventHandler(rm)

	// Chain
	rm.Add("Version", func(kn *kernel.Kernel, ID interface{}, arg *rpc.Argument) (interface{}, error) {
		return kn.Provider().Version(), nil
	})
	rm.Add("Height", func(kn *kernel.Kernel, ID interface{}, arg *rpc.Argument) (interface{}, error) {
		return kn.Provider().Height(), nil
	})
	rm.Add("LastHash", func(kn *kernel.Kernel, ID interface{}, arg *rpc.Argument) (interface{}, error) {
		return kn.Provider().LastHash(), nil
	})
	rm.Add("Hash", func(kn *kernel.Kernel, ID interface{}, arg *rpc.Argument) (interface{}, error) {
		if arg.Len() < 1 {
			return nil, rpc.ErrInvalidArgument
		}
		height, err := arg.Uint32(0)
		if err != nil {
			return nil, err
		}
		h, err := kn.Provider().Hash(height)
		if err != nil {
			return nil, err
		}
		return h, nil
	})
	rm.Add("Header", func(kn *kernel.Kernel, ID interface{}, arg *rpc.Argument) (interface{}, error) {
		if arg.Len() < 1 {
			return nil, rpc.ErrInvalidArgument
		}
		height, err := arg.Uint32(0)
		if err != nil {
			return nil, err
		}
		h, err := kn.Provider().Header(height)
		if err != nil {
			return nil, err
		}
		return h, nil
	})
	rm.Add("Block", func(kn *kernel.Kernel, ID interface{}, arg *rpc.Argument) (interface{}, error) {
		if arg.Len() < 1 {
			return nil, rpc.ErrInvalidArgument
		}
		height, err := arg.Uint32(0)
		if err != nil {
			return nil, err
		}
		cd, err := kn.Provider().Data(height)
		if err != nil {
			return nil, err
		}
		b := &block.Block{
			Header: cd.Header.(*block.Header),
			Body:   cd.Body.(*block.Body),
		}
		return b, nil
	})

	go func() {
		if err := rm.Run(kn, ":58000"); err != nil {
			if http.ErrServerClosed != err {
				panic(err)
			}
		}
	}()

	cm.Wait()
}
