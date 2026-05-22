package cli

import (
	"fmt"
	"os"

	"github.com/playfairs/veego/internal/api"
	"github.com/playfairs/veego/internal/config"
)

var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Device commands",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all devices",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		client := api.NewClient(cfg.APIKey)
		resp, err := client.GetDevices()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting devices: %v\n", err)
			os.Exit(1)
		}

		for _, device := range resp.Data.Devices {
			fmt.Printf("%s (%s) - %s\n", device.DeviceName, device.Device, device.Model)
		}
	},
}

var statusCmd = &cobra.Command{
	Use:   "status <device>",
	Short: "Show device status",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deviceName := args[0]

		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		device, ok := cfg.Devices[deviceName]
		if !ok {
			fmt.Fprintf(os.Stderr, "Device not found: %s\n", deviceName)
			os.Exit(1)
		}

		client := api.NewClient(cfg.APIKey)
		resp, err := client.GetDeviceState(device.ID, device.Model)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting device state: %v\n", err)
			os.Exit(1)
		}

		state := resp.Data
		fmt.Printf("Power: %s\n", state.PowerState)
		fmt.Printf("Brightness: %d\n", state.Brightness)
		if state.Color.R != 0 || state.Color.G != 0 || state.Color.B != 0 {
			fmt.Printf("Color: RGB(%d, %d, %d)\n", state.Color.R, state.Color.G, state.Color.B)
		}
	},
}

func init() {
	RootCmd.AddCommand(deviceCmd)
	deviceCmd.AddCommand(listCmd)
	deviceCmd.AddCommand(statusCmd)
}
