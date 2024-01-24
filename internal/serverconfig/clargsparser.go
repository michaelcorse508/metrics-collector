package serverconfig

import (
	"flag"
	"fmt"
	"github.com/bazookajoe1/metrics-collector/internal/auxiliary"
	"os"
	"strings"
)

type CLArgParams struct {
	Host          string
	Port          string
	StoreInterval int
	FilePath      string
	Restore       *bool
	DatabaseDSN   string
	SecretKey     string
}

func (a *CLArgParams) String() string {
	return a.Host + ":" + a.Port
}

func (a *CLArgParams) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return fmt.Errorf("need address in the form Host:Port")
	}

	err := auxiliary.ValidateIP(hp[0])
	if err != nil {
		return err
	}

	err = auxiliary.ValidatePort(hp[1])
	if err != nil {
		return err
	}
	a.Host = hp[0]
	a.Port = hp[1]
	return nil

}

func (a *CLArgParams) GetAddress() string {
	return a.Host
}

func (a *CLArgParams) GetPort() string {
	return a.Port
}

func (a *CLArgParams) GetStoreInterval() int {
	return a.StoreInterval
}

func (a *CLArgParams) GetFilePath() string {
	return a.FilePath
}

func (a *CLArgParams) GetRestoreFlag() *bool {
	return a.Restore
}

func (a *CLArgParams) GetDatabaseDSN() string {
	return a.DatabaseDSN
}

func (a *CLArgParams) GetSecretKey() string {
	return a.SecretKey
}

func ArgParse() *CLArgParams {
	na := &CLArgParams{
		Host:          "localhost",
		Port:          "8080",
		StoreInterval: 300,
		FilePath:      os.TempDir() + "\\metrics-db.json",
		Restore:       new(bool),
	}

	flag.Var(na, "a", "Server listen point in format: `Host:Port`")
	flag.IntVar(&na.StoreInterval, "i", 300, "Save metrics to file interval")
	flag.StringVar(&na.FilePath, "f", os.TempDir()+"\\metrics-db.json", "File name where to store metrics")
	flag.BoolVar(na.Restore, "r", false, "Load storage from file or not")
	flag.StringVar(&na.DatabaseDSN, "d", "", "Database connection string")
	flag.StringVar(&na.SecretKey, "k", "", "Sign key")
	flag.Parse()

	return na
}
