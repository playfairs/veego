package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/playfairs/veego/internal/api"
	"github.com/playfairs/veego/internal/config"
)

var lightCmd = &cobra.Command{
	Use:   "light",
	Short: "Light control commands",
}

var onCmd = &cobra.Command{
	Use:   "on <device>",
	Short: "Turn on device",
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
		err = client.ControlDevice(device.ID, device.Model, "turn", "on")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error controlling device: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Device turned on")
	},
}

var offCmd = &cobra.Command{
	Use:   "off <device>",
	Short: "Turn off device",
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
		err = client.ControlDevice(device.ID, device.Model, "turn", "off")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error controlling device: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Device turned off")
	},
}

var toggleCmd = &cobra.Command{
	Use:   "toggle <device>",
	Short: "Toggle device power",
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

		var newState string
		if resp.Data.PowerState == "on" {
			newState = "off"
		} else {
			newState = "on"
		}

		err = client.ControlDevice(device.ID, device.Model, "turn", newState)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error controlling device: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Device turned %s\n", newState)
	},
}

var brightnessCmd = &cobra.Command{
	Use:   "brightness <device> <0-100>",
	Short: "Set device brightness",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		deviceName := args[0]
		brightnessStr := args[1]

		brightness, err := strconv.Atoi(brightnessStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid brightness value: %s\n", brightnessStr)
			os.Exit(1)
		}

		if brightness < 0 || brightness > 100 {
			fmt.Fprintf(os.Stderr, "Brightness must be between 0 and 100\n")
			os.Exit(1)
		}

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
		err = client.ControlDevice(device.ID, device.Model, "brightness", brightness)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error controlling device: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Brightness set to %d\n", brightness)
	},
}

var colorCmd = &cobra.Command{
	Use:   "color <device> <color>",
	Short: "Set device color",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		deviceName := args[0]
		colorStr := args[1]

		rgb, err := parseColor(colorStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid color: %v\n", err)
			os.Exit(1)
		}

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
		colorValue := map[string]int{
			"r": rgb[0],
			"g": rgb[1],
			"b": rgb[2],
		}
		err = client.ControlDevice(device.ID, device.Model, "color", colorValue)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error controlling device: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Color set to RGB(%d, %d, %d)\n", rgb[0], rgb[1], rgb[2])
	},
}

func parseColor(color string) ([3]int, error) {
	color = strings.ToLower(strings.TrimSpace(color))

	namedColors := map[string][3]int{
		"red":    {255, 0, 0},
		"green":  {0, 255, 0},
		"blue":   {0, 0, 255},
		"white":  {255, 255, 255},
		"yellow": {255, 255, 0},
		"purple": {128, 0, 128},
		"cyan":   {0, 255, 255},
	}

	if rgb, ok := namedColors[color]; ok {
		return rgb, nil
	}

	if strings.HasPrefix(color, "#") {
		color = color[1:]
	}

	if len(color) == 6 {
		var r, g, b int
		_, err := fmt.Sscanf(color, "%02x%02x%02x", &r, &g, &b)
		if err != nil {
			return [3]int{}, fmt.Errorf("invalid hex color format")
		}
		return [3]int{r, g, b}, nil
	}

	return [3]int{}, fmt.Errorf("unknown color format")
}

func init() {
	RootCmd.AddCommand(lightCmd)
	lightCmd.AddCommand(onCmd)
	lightCmd.AddCommand(offCmd)
	lightCmd.AddCommand(toggleCmd)
	lightCmd.AddCommand(brightnessCmd)
	lightCmd.AddCommand(colorCmd)
}
