package client

import (
    "bufio"
    "encoding/binary"
    "fmt"
    "io"
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
        _, err := writer.Write(pack(pkt, client.TimeDiff))
        if err != nil {
            fmt.Println(common.Red("Packet Send Failed! " + err.Error()))
            go retry(client)
            return
        }
    }
}

func pack(frame []byte, td int64) []byte {
    // ------ Packet ------
    // timestamp (8 bytes)
    // content-length (8 bytes)
    // content (content-length bytes)
    // ------  End   ------

    timePkt := make([]byte, common.TimeStampPacketSize)
    binary.LittleEndian.PutUint64(timePkt, uint64(time.Now().UnixMicro() - td))

    contentLengthPkt := make([]byte, common.ContentLengthPacketSize)
    binary.LittleEndian.PutUint64(contentLengthPkt, uint64(len(frame)))

    var pkt []byte
    pkt = append(pkt, timePkt...)
    pkt = append(pkt, contentLengthPkt...)
    pkt = append(pkt, frame...)

    return pkt
}

func handleReceive(client *Client) {
    reader := bufio.NewReader(client.TCPConn)
    timeStamp := make([]byte, common.TimeStampPacketSize)

    for {
        _, err := io.ReadFull(reader, timeStamp)
        if err != nil {
            fmt.Println(common.Red(client.TCPConn.RemoteAddr().String() + " Down! " + err.Error()))
            break
        }
        curTime := time.Now().UnixMicro()
        recvTime := int64(binary.LittleEndian.Uint64(timeStamp))

        client.TimeDiff = curTime - recvTime
    }
}
