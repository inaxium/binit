package binit

import (
	"fmt"
	spinconfig "github.com/fermyon/spin/sdk/go/config"
	"github.com/inaxium/bconst/arango/atxt"
	"github.com/inaxium/bconst/emsg"
	"github.com/inaxium/bconst/pub/ptxt"
	"github.com/inaxium/bmod"
	auth "github.com/inaxium/btool/auth/arango"
	"github.com/inaxium/btool/txn"
	"runtime"
	"time"
)

func ReadArangoCfg(source string) (cfg bmod.ArangoCfg, err error) {
	if source == ptxt.Spin {
		if cfg.AuthUrl, err = getSpinConfig(atxt.ArangoUrlAuth); err != nil {
			return cfg, err
		}
		if cfg.Url, err = getSpinConfig(atxt.ArangoUrlDocument); err != nil {
			return cfg, err
		}
	} else {
		return cfg, fmt.Errorf("%s : %s", emsg.InvalidSource, source)
	}

	return cfg, nil
}

func ReadVersion() (version string, err error) {
	version, err = getSpinConfig(ptxt.VersionL)
	return version, err
}

func ReadMeta(goVersion string) bmod.Meta {

	return bmod.Meta{
		TransactionId:   txn.GenId(),
		TransactionTime: time.Now(),
		GoOS:            runtime.GOOS,
		GoVersion:       goVersion,
	}
}

func Bootstrap(goVersion string) (meta bmod.Meta,
	arangoCfg bmod.ArangoCfg, token bmod.Token, err error) {

	creds := bmod.Credentials{}
	meta = ReadMeta(goVersion)
	meta.Version, err = ReadVersion()
	arangoCfg, err = ReadArangoCfg(ptxt.Spin)
	creds, err = auth.ReadCredentials(ptxt.Spin)
	token, err = auth.Login(creds, arangoCfg)

	return meta, arangoCfg, token, nil

}

func getSpinConfig(key string) (string, error) {
	return spinconfig.Get(key)
}
