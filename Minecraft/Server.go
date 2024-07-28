package Minecraft

import (
	"bufio"
	"encoding/binary"
	json2 "encoding/json"
	"net"
	"strconv"
	"strings"
	"time"
)

type ServerInfo struct {
	Version struct {
		Name     string `json:"name"`
		Software string
		Protocol int32 `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int32 `json:"max"`
		Online int32 `json:"online"`
		Sample []struct {
			Name string `json:"name"`
			Id   string `json:"id"`
		} `json:"sample"`
	} `json:"players"`
	Description struct {
		Text string `json:"text"`
	} `json:"description"`
	Favicon            string `json:"favicon"`
	EnforcesSecureChat bool   `json:"enforcesSecureChat"`
	PreviewsChat       bool   `json:"previewsChat"`
}

type Server struct {
	ip      string
	port    uint16
	version string
	sock    net.Conn
}

func NewMinecraftServer(ip string, port uint16, version string) (server Server, err error) {
	addr := net.JoinHostPort(ip, strconv.Itoa(int(port)))
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return
	}
	server.sock = conn
	server.ip = ip
	server.port = port
	server.version = version
	return
}

func (srv *Server) GetStatusRequest() (info ServerInfo, err error) {
	handshake, err := CreateHandshake(srv.ip, srv.port, uint32(GetProtocol(srv.version)))
	if err != nil {
		return
	}
	_, err = srv.sock.Write(handshake.GetBuffer())
	if err != nil {
		return
	}
	serverListRequest, err := CreateStatusRequest()
	if err != nil {
		return
	}
	_, err = srv.sock.Write(serverListRequest.GetBuffer())
	if err != nil {
		return
	}

	reader := bufio.NewReader(srv.sock)
	_, err = binary.ReadUvarint(reader)
	if err != nil {
		return
	}

	_, err = reader.ReadByte()
	if err != nil {
		return
	}

	jsonLength, err := binary.ReadUvarint(reader)
	if err != nil {
		return
	}

	var read uint64 = 0
	json := make([]byte, jsonLength)
	for read < jsonLength {
		n, _ := reader.Read(json[read:jsonLength])
		read = read + uint64(n)
	}

	err = json2.Unmarshal(json, &info)
	if err != nil {
		return
	}

	versionString := info.Version.Name
	versionSplit := strings.Split(versionString, " ")
	if len(versionSplit) > 1 {
		info.Version.Name = versionSplit[1]
		info.Version.Software = versionSplit[0]
	} else {
		info.Version.Name = versionString
		info.Version.Software = "Vanilla"
	}

	return
}
