// Copyright (c) 2018 Yuriy Lisovskiy
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/YuriyLisovskiy/blockchain-go/src/utils"
)

var configLocation string

func init() {
	// get path of running app
	absPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	configLocation = absPath + "/config.json"
}

// Represents node configuration
type Config struct {
	Ip          string `json:"ip"`
	Port        int    `json:"port"`
	ChainPath   string `json:"chain_path"`
	WalletsPath string `json:"wallets_path"`
}

// Default returns default node configuration.
func Default() (Config, error) {
	var cfg Config

	// get path of running app
	absPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return cfg, err
	}

	// create data directory for storing
	err = os.MkdirAll(absPath+"/data", os.ModePerm)
	if err != nil {
		return cfg, err
	}

	// get node ip address
	ip, err := getIp()
	if err != nil {
		return cfg, err
	}

	// setup data
	cfg.Ip = ip
	cfg.Port = 8000
	cfg.ChainPath = absPath + "/data/" + fmt.Sprintf(utils.DBFile, cfg.Port)
	cfg.WalletsPath = absPath + "/data/" + fmt.Sprintf(utils.WalletFile, cfg.Port)

	return cfg, nil
}

// SetIp sets node's ip.
func (cfg Config) SetIp(ip string) Config {
	cfg.Ip = ip
	return cfg
}

// SetPort sets node's port.
func (cfg Config) SetPort(port int) Config {
	cfg.Port = port
	if cfg.ChainPath != "" {
		cfg.ChainPath = cfg.ChainPath[:strings.LastIndex(cfg.ChainPath, "/")+1] + fmt.Sprintf(utils.DBFile, cfg.Port)
	}
	if cfg.WalletsPath != "" {
		cfg.WalletsPath = cfg.WalletsPath[:strings.LastIndex(cfg.WalletsPath, "/")+1] + fmt.Sprintf(utils.WalletFile, cfg.Port)
	}
	return cfg
}

// SetChainPath sets a path to block chain database.
func (cfg Config) SetChainPath(path string) Config {
	cfg.ChainPath = path
	return cfg
}

// SetWalletsPath sets a path to wallets location.
func (cfg Config) SetWalletsPath(path string) Config {
	cfg.WalletsPath = path
	return cfg
}

// Exists checks if configuration file exists on disk.
func Exists() bool {
	_, err := os.Stat(configLocation)
	return !os.IsNotExist(err)
}

// Save marshals config and saves it to a file.
func (cfg Config) Save() error {
	// marshal config content to json
	bCfg, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configLocation, bCfg, 0666)
}

// Load loads config from file.
func LoadConfig() (Config, error) {
	var cfg Config
	b, err := ioutil.ReadFile(configLocation)
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
