package api

type Device struct {
	Device     string `json:"device"`
	Model      string `json:"model"`
	DeviceName string `json:"deviceName"`
}

type DevicesResponse struct {
	Message string   `json:"message"`
	Data    Devices `json:"data"`
}

type Devices struct {
	Devices []Device `json:"devices"`
}

type DeviceState struct {
	Online       bool    `json:"online"`
	PowerState   string  `json:"powerState"`
	Brightness  int     `json:"brightness"`
	Color        Color   `json:"color"`
}

type Color struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
}

type DeviceStateResponse struct {
	Message string       `json:"message"`
	Data    DeviceState  `json:"data"`
}

type ControlRequest struct {
	Device string `json:"device"`
	Model  string `json:"model"`
	Cmd    Cmd    `json:"cmd"`
}

type Cmd struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type ControlResponse struct {
	Message string `json:"message"`
}

func (c *Client) GetDevices() (*DevicesResponse, error) {
	var result DevicesResponse
	err := c.do("GET", "/devices", nil, &result)
	return &result, err
}

func (c *Client) GetDeviceState(deviceID, model string) (*DeviceStateResponse, error) {
	path := fmt.Sprintf("/devices/state?device=%s&model=%s", deviceID, model)
	var result DeviceStateResponse
	err := c.do("GET", path, nil, &result)
	return &result, err
}

func (c *Client) ControlDevice(deviceID, model string, cmdName string, cmdValue interface{}) error {
	req := ControlRequest{
		Device: deviceID,
		Model:  model,
		Cmd: Cmd{
			Name:  cmdName,
			Value: cmdValue,
		},
	}
	var result ControlResponse
	return c.do("PUT", "/devices/control", req, &result)
}
