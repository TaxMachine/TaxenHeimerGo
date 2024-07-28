package Minecraft

type Version struct {
	version  string
	protocol int
}

var ProtocolVersions = []Version{
	{"1.20.5", 766},
	{"1.20.4", 765},
	{"1.20.3", 765},
	{"1.20.2", 764},
	{"1.20.1", 763},
	{"1.20", 763},
	{"1.19.4", 762},
	{"1.19.3", 761},
	{"1.19.2", 760},
	{"1.19.1", 760},
	{"1.19", 759},
	{"1.18.2", 758},
	{"1.18.1", 757},
	{"1.18", 757},
	{"1.17.1", 756},
	{"1.17", 755},
	{"1.16.5", 754},
	{"1.16.4", 754},
	{"1.16.3", 753},
	{"1.16.2", 751},
	{"1.16.1", 736},
	{"1.16", 735},
	{"1.15.2", 578},
	{"1.15.1", 575},
	{"1.15", 573},
	{"1.14.4", 498},
	{"1.14.3", 490},
	{"1.14.2", 485},
	{"1.14.1", 480},
	{"1.14", 477},
	{"1.13.2", 404},
	{"1.13.1", 401},
	{"1.13", 393},
	{"1.12.2", 340},
	{"1.12.1", 338},
	{"1.12", 335},
	{"1.11.2", 316},
	{"1.11.1", 316},
	{"1.11", 315},
	{"1.10.2", 210},
	{"1.10.1", 210},
	{"1.10", 210},
	{"1.9.4", 110},
	{"1.9.3", 109},
	{"1.9.2", 176},
	{"1.9.1", 108},
	{"1.9", 169},
	{"1.8.9", 47},
	{"1.8.8", 47},
	{"1.8.7", 47},
	{"1.8.6", 47},
	{"1.8.5", 47},
	{"1.8.4", 47},
	{"1.8.3", 47},
	{"1.8.2", 47},
	{"1.8.1", 47},
	{"1.8", 47},
	{"1.7.10", 5},
	{"1.7.9", 5},
	{"1.7.8", 5},
	{"1.7.7", 5},
	{"1.7.6", 5},
	{"1.7.5", 4},
	{"1.7.4", 4},
	{"1.7.3", 4},
	{"1.7.2", 4},
	{"1.7.1", 3},
	{"1.7", 3},
}

func GetProtocol(version string) int {
	for _, p := range ProtocolVersions {
		if p.version == version {
			return p.protocol
		}
	}
	return 0
}

func CreateHandshake(ip string, port uint16, protocolVersion uint32) (packet Packet, err error) {
	handshake := NewMinecraftPacket()
	err = handshake.WriteVarInt(0x00)
	if err != nil {
		return
	}
	err = handshake.WriteVarInt(protocolVersion)
	if err != nil {
		return
	}
	err = handshake.WriteString(ip)
	if err != nil {
		return
	}
	handshake.WriteShort(port)
	err = handshake.WriteVarInt(1)
	if err != nil {
		return
	}

	err = packet.WriteVarInt(uint32(handshake.Size()))
	if err != nil {
		return
	}
	packet.WriteBuffer(handshake.GetBuffer())
	return packet, nil
}

func CreatePing() (packet Packet, err error) {
	ping := NewMinecraftPacket()
	err = ping.WriteVarInt(0x01)
	if err != nil {
		return
	}
	err = ping.WriteVarInt(343838093)
	if err != nil {
		return
	}
	err = packet.WriteVarInt(uint32(ping.Size()))
	if err != nil {
		return
	}
	packet.WriteBuffer(ping.GetBuffer())
	return packet, nil
}

func CreateStatusRequest() (packet Packet, err error) {
	var status Packet
	err = status.WriteVarInt(0x00)
	if err != nil {
		return
	}

	err = packet.WriteVarInt(uint32(status.Size()))
	if err != nil {
		return
	}
	packet.WriteBuffer(status.GetBuffer())
	return packet, nil
}
