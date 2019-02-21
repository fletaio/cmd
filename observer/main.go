package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"git.fleta.io/fleta/core/reward"
	"git.fleta.io/fleta/framework/config"
	"github.com/dgraph-io/badger"

	"git.fleta.io/fleta/core/data"
	"git.fleta.io/fleta/core/kernel"
	"git.fleta.io/fleta/core/key"
	"git.fleta.io/fleta/core/observer"

	"git.fleta.io/fleta/common"
)

// Config is a configuration for the cmd
type Config struct {
	ObserverKeyMap map[string]string
	KeyHex         string
	ObseverPort    int
	FormulatorPort int
	StoreRoot      string
}

func main() {
	var cfg Config
	if err := config.LoadFile("./config.toml", &cfg); err != nil {
		panic(err)
	}
	if len(cfg.StoreRoot) == 0 {
		cfg.StoreRoot = "./observer"
	}

	var obkey key.Key
	if bs, err := hex.DecodeString(cfg.KeyHex); err != nil {
		panic(err)
	} else if Key, err := key.NewMemoryKeyFromBytes(bs); err != nil {
		panic(err)
	} else {
		obkey = Key
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
	if s, err := kernel.NewStore(cfg.StoreRoot+"/kernel", BlockchainVersion, act, tran, false); err != nil {
		if err != badger.ErrTruncateNeeded {
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
		ObserverKeyMap: ObserverKeyBoolMap,
	}, ks, rd, GenesisContextData)
	if err != nil {
		panic(err)
	}
	closable = kn

	obcfg := &observer.Config{
		ChainCoord:     GenCoord,
		Key:            obkey,
		ObserverKeyMap: ObserverKeyMap,
	}
	ob, err := observer.NewObserver(obcfg, kn)
	if err != nil {
		panic(err)
	}
	closable = ob

	ob.Run(":"+strconv.Itoa(cfg.ObseverPort), ":"+strconv.Itoa(cfg.FormulatorPort))
}

type Closable interface {
	Close()
}
