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

    go handleSend(client)

    img := gocv.NewMat()
    for {
        client.Camera.Read(&img)
        client.SendQueue <- img.ToBytes()
        gocv.WaitKey(1)
    }
}
