// Copyright 2020 Coinbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package configuration

import (
	"errors"
	"os"
	"testing"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/stretchr/testify/assert"
	"github.com/ubiq/go-ubiq/v7/params"
	"github.com/ubiq/rosetta-ubiq/ubiq"
)

func TestLoadConfiguration(t *testing.T) {
	tests := map[string]struct {
		Mode          string
		Network       string
		Port          string
		Geth          string
		SkipGethAdmin string

		cfg *Configuration
		err error
	}{
		"no envs set": {
			err: errors.New("MODE must be populated"),
		},
		"only mode set": {
			Mode: string(Online),
			err:  errors.New("NETWORK must be populated"),
		},
		"only mode and network set": {
			Mode:    string(Online),
			Network: Mainnet,
			err:     errors.New("PORT must be populated"),
		},
		"all set (mainnet)": {
			Mode:          string(Online),
			Network:       Mainnet,
			Port:          "1000",
			SkipGethAdmin: "FALSE",
			cfg: &Configuration{
				Mode: Online,
				Network: &types.NetworkIdentifier{
					Network:    ubiq.MainnetNetwork,
					Blockchain: ubiq.Blockchain,
				},
				Params:                 params.MainnetChainConfig,
				GenesisBlockIdentifier: ubiq.MainnetGenesisBlockIdentifier,
				Port:                   1000,
				GethURL:                DefaultGethURL,
				GethArguments:          ubiq.MainnetGethArguments,
				SkipGethAdmin:          false,
			},
		},
		"all set (mainnet) + geth": {
			Mode:          string(Online),
			Network:       Mainnet,
			Port:          "1000",
			Geth:          "http://blah",
			SkipGethAdmin: "TRUE",
			cfg: &Configuration{
				Mode: Online,
				Network: &types.NetworkIdentifier{
					Network:    ubiq.MainnetNetwork,
					Blockchain: ubiq.Blockchain,
				},
				Params:                 params.MainnetChainConfig,
				GenesisBlockIdentifier: ubiq.MainnetGenesisBlockIdentifier,
				Port:                   1000,
				GethURL:                "http://blah",
				RemoteGeth:             true,
				GethArguments:          ubiq.MainnetGethArguments,
				SkipGethAdmin:          true,
			},
		},
		"all set (testnet)": {
			Mode:          string(Online),
			Network:       Testnet,
			Port:          "1000",
			SkipGethAdmin: "TRUE",
			cfg: &Configuration{
				Mode: Online,
				Network: &types.NetworkIdentifier{
					Network:    ubiq.RopstenNetwork,
					Blockchain: ubiq.Blockchain,
				},
				Params:                 params.RopstenChainConfig,
				GenesisBlockIdentifier: ubiq.RopstenGenesisBlockIdentifier,
				Port:                   1000,
				GethURL:                DefaultGethURL,
				GethArguments:          ubiq.RopstenGethArguments,
				SkipGethAdmin:          true,
			},
		},
		"invalid mode": {
			Mode:    "bad mode",
			Network: Ropsten,
			Port:    "1000",
			err:     errors.New("bad mode is not a valid mode"),
		},
		"invalid network": {
			Mode:    string(Offline),
			Network: "bad network",
			Port:    "1000",
			err:     errors.New("bad network is not a valid network"),
		},
		"invalid port": {
			Mode:    string(Offline),
			Network: Ropsten,
			Port:    "bad port",
			err:     errors.New("unable to parse port bad port"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			os.Setenv(ModeEnv, test.Mode)
			os.Setenv(NetworkEnv, test.Network)
			os.Setenv(PortEnv, test.Port)
			os.Setenv(GethEnv, test.Geth)
			os.Setenv(SkipGethAdminEnv, test.SkipGethAdmin)

			cfg, err := LoadConfiguration()
			if test.err != nil {
				assert.Nil(t, cfg)
				assert.Contains(t, err.Error(), test.err.Error())
			} else {
				assert.Equal(t, test.cfg, cfg)
				assert.NoError(t, err)
			}
		})
	}
}
