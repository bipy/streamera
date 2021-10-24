package client

import (
	"fmt"
	"gocv.io/x/gocv"
	"net"
	"streamera/common"
)

type Client struct {
	RemoteIP   net.IP
	RemotePort int
	TCPConn    *net.TCPConn
	Camera     *gocv.VideoCapture
	SendQueue  chan []byte
	TimeDiff   int64
}

func NewClient(ip net.IP, port int, deviceID int) (*Client, error) {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   ip,
		Port: port,
	})
	if err != nil {
		fmt.Println(common.Red("generate udp client failed!"))
		return nil, err
	}

	sendQueue := make(chan []byte, 180)

	cam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		return nil, err
	}

	client := &Client{
		RemoteIP:   ip,
		RemotePort: port,
		TCPConn:    conn,
		Camera:     cam,
		SendQueue:  sendQueue,
		TimeDiff:   0,
	}

	return client, nil
}

func RunClient(client *Client) {
	defer func(conn *net.TCPConn) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(client.TCPConn)

	go handleReceive(client)
	go handleSend(client)

	mat := gocv.NewMat()
	for {
		client.Camera.Read(&mat)
		client.SendQueue <- encodeImage(&mat)
		gocv.WaitKey(1)
	}
}
