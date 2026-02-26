package room

// Response 房间配置响应
type Response struct {
	Config ConfigData `json:"config"`
	Result int        `json:"result"`
}

// ConfigData 房间配置数据
type ConfigData struct {
	Room        Endpoint `json:"room"`
	Proxy       Endpoint `json:"proxy"`
	ProxyOnly   string   `json:"proxy_only"`
	Punch       Endpoint `json:"punch"`
	NetworkType int      `json:"network_type"`
	RoomName    string   `json:"room_name"`
	BlockType   string   `json:"block_type"`
	AreaType    int      `json:"area_type"`
}

// Endpoint 服务端点
type Endpoint struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}
