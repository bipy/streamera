package server

import (
    "bufio"
    "bytes"
    "encoding/binary"
    "fmt"
    "io"
    "net"
    "streamera/common"
    "time"
)

func runListener(server *Server) {
    for {
        server.TCPListenerMutex.Lock()
        server.FrameMutex.Lock()
        conn, err := server.TCPListener.Accept()
        if err != nil {
            fmt.Println(common.Red("broken connection " + err.Error()))
            continue
        }
        go sendTimeStamp(conn)
        go handleConn(server, conn)
    }
}

func handleConn(server *Server, conn net.Conn) {
    defer func(server *Server, conn net.Conn) {
        err := conn.Close()
        if err != nil {
            fmt.Println(common.Red("tcp conn close failed! " + err.Error()))
        }
    }(server, conn)

    defer server.TCPListenerMutex.Unlock()

    server.FrameMutex.Unlock()

    reader := bufio.NewReader(conn)

    timeStamp := make([]byte, common.TimeStampPacketSize)
    contentLength := make([]byte, common.ContentLengthPacketSize)

    for {
        _, err := io.ReadFull(reader, timeStamp)
        if err != nil {
            fmt.Println(common.Red(conn.RemoteAddr().String() + " Down! " + err.Error()))
            break
        }
        curTime := time.Now().UnixMicro()
        recvTime := int64(binary.LittleEndian.Uint64(timeStamp))

        _, err = io.ReadFull(reader, contentLength)
        if err != nil {
            fmt.Println(common.Red(conn.RemoteAddr().String() + " Down! " + err.Error()))
            break
        }

        length := int64(binary.LittleEndian.Uint64(contentLength))

        buf := make([]byte, length)
        _, err = io.ReadFull(reader, buf)
        if err != nil {
            fmt.Println(common.Red(conn.RemoteAddr().String() + " Down! " + err.Error()))
            break
        }

        server.FrameMutex.Lock()
        server.Frame = new(bytes.Buffer)
        server.Frame.Write(buf)
        server.FrameMutex.Unlock()

        server.Counter.Mutex.Lock()
        server.Counter.LatencyRealTime = (curTime - recvTime) >> 1
        server.Counter.ByteCount += length + common.HeadPacketSize
        server.Counter.Mutex.Unlock()
    }
}

func sendTimeStamp(conn net.Conn) {
    writer := bufio.NewWriter(conn)
    for {
        time.Sleep(time.Millisecond << 7)

        timePkt := make([]byte, common.TimeStampPacketSize)
        binary.LittleEndian.PutUint64(timePkt, uint64(time.Now().UnixMicro()))

        _, err := writer.Write(timePkt)
        if err != nil {
            return
        }

        err = writer.Flush()
        if err != nil {
            return
        }
    }
}
