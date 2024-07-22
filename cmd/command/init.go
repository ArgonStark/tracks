package command

import (
	"fmt"
	"github.com/airchains-network/tracks/config"
	logs "github.com/airchains-network/tracks/log"

	//logs "github.com/airchains-network/tracks/log"
	"github.com/airchains-network/tracks/p2p"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

type Configs struct {
	moniker     string
	stationType string
	daType      string
	daRPC       string
	daKey       string
	stationRPC  string
	stationAPI  string
}

func InitConfigs(cmd *cobra.Command) (*Configs, error) {
	var configs Configs
	var err error

	configs.moniker, err = cmd.Flags().GetString("moniker")
	if err != nil {
		return nil, fmt.Errorf("failed to get flag 'moniker': %w", err)
	}

	configs.stationType, err = cmd.Flags().GetString("stationType")
	if err != nil {
		return nil, fmt.Errorf("failed to get flag 'stationType': %w", err)
	}

	configs.daType, err = cmd.Flags().GetString("daType")
	if err != nil {
		return nil, fmt.Errorf("failed to get flag 'daType': %w", err)
	}

	validTypes := map[string]bool{
		"avail":    true,
		"celestia": true,
		"eigen":    true,
		"mock":     true,
	}
	if _, isValid := validTypes[configs.daType]; !isValid {
		logs.Log.Error("invalid daType. Must be one of: avail, celestia, eigen, mock")
		return nil, fmt.Errorf("invalid daType: %s", configs.daType)
	}

	configs.daRPC, err = cmd.Flags().GetString("daRpc")
	if err != nil {
		return nil, fmt.Errorf("failed to get flag 'daRpc': %w", err)
	}

	configs.daKey, err = cmd.Flags().GetString("daKey")
	if err != nil {
		return nil, fmt.Errorf("failed to get flag 'daKey': %w", err)
	}

	configs.stationRPC, err = cmd.Flags().GetString("stationRpc")
	if err != nil {
		return nil, fmt.Errorf("failed to get flag 'stationRPC': %w", err)
	}

	configs.stationAPI, err = cmd.Flags().GetString("stationAPI")
	if err != nil {
		return nil, fmt.Errorf("failed to get flag 'stationAPI': %w", err)
	}

	return &configs, nil
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize the sequencer nodes",
	Run: func(cmd *cobra.Command, args []string) {
		configs, err := InitConfigs(cmd)
		if err != nil {
			logs.Log.Error(err.Error())
			return
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			logs.Log.Error("Failed to get user home directory:" + err.Error())
			return
		}

		tracksDir := filepath.Join(homeDir, config.DefaultTracksDir)

		conf := config.DefaultConfig()
		peerGen := p2p.NewPeerGenerator("/ip4/0.0.0.0/tcp/2300", false)
		peerID, err := peerGen.GeneratePeerID()

		conf.BaseConfig.RootDir = tracksDir
		conf.DA.DaType = configs.daType
		conf.DA.DaRPC = configs.daRPC
		conf.DA.DaKey = configs.daKey
		conf.Station.StationType = configs.stationType
		conf.Station.StationRPC = configs.stationRPC
		conf.Station.StationAPI = configs.stationAPI
		conf.P2P.NodeId = peerID
		conf.SetRoot(conf.BaseConfig.RootDir)

		success := config.CreateConfigFile(conf.BaseConfig.RootDir, conf)
		if !success {
			logs.Log.Error("Unable to generate a config file. Please check the error and try again.")
			return
		}

		logs.Log.Info("Track initialization successful")
	},
}
