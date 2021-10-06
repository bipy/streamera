package client

import (
    "bufio"
    "encoding/binary"
    "fmt"
    "net"
    "streamera/common"
    "time"
)

func retry(client *Client) {
    conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
        IP:   client.RemoteIP,
        Port: client.RemotePort,
    })
    if err != nil {
        fmt.Println(common.Red(err.Error()))
        time.Sleep(time.Second * 3)
        go retry(client)
        return
    }
    client.TCPConn = conn
    go handleSend(client)
}

func handleSend(client *Client) {
    writer := bufio.NewWriter(client.TCPConn)
    for pkt := range client.SendQueue {
        _, err := writer.Write(pack(pkt))
        if err != nil {
            fmt.Println(common.Red("Packet Send Failed! " + err.Error()))
            go retry(client)
            return
        }
    }
}

func pack(frame []byte) []byte {
    // ------ Packet ------
    // timestamp (8 bytes)
    // content (2,764,800 bytes)
    // ------  End   ------

    timePkt := make([]byte, common.TimeStampPacketSize)
    binary.LittleEndian.PutUint64(timePkt, uint64(time.Now().UnixMicro()))

    var pkt []byte
    pkt = append(pkt, timePkt...)
    pkt = append(pkt, frame...)

    return pkt
}