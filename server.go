package main

import (
    "net"
    "os"
    "strconv"
    "./data"
)


func main() {
    println("Starting the server")

    listener, err := net.Listen("tcp", "0.0.0.0:6666")
    if err != nil {
        println("error listening:", err.Error())
        os.Exit(1)
    }

    dataDict := data.NewDataDict()

    for {
        conn, err := listener.Accept()
        if err != nil {
            println("Error accept:", err.Error())
            return
        }
        go HandleClient(conn, dataDict)
    }
}


func HandleClient(conn net.Conn, dataDict *data.DataDict) {
    cmdBuf := make([]byte, 8)
    for {
        _, err := conn.Read(cmdBuf)
        if err != nil {
            println("Error reading:", err.Error())
            continue
        }

        tuple := string(cmdBuf)
        l, cmd := tuple[0:4], tuple[4:8]
        length, _ := strconv.Atoi(l)

        dataBuf := make([]byte, length)
        _, err = conn.Read(dataBuf)
        if err != nil {
            println("Error reading:", err.Error())
            break
        }
        data := string(dataBuf)

        var result string
        switch {
            case cmd == "GET_":
                result = GetCommand(data, dataDict)
            case cmd == "SET_":
                result = SetCommand(data, dataDict)
        }
        if err != nil {
            conn.Write([]byte(result))
        } else {
            conn.Write([]byte(result))
        }
    }
}

func GetCommand(data string, dataDict *data.DataDict) (string) {
    key := data[0:32]
    value := (*dataDict).Get(key)
    return "key:"+key+"value:"+value
}

func SetCommand(data string, dataDict *data.DataDict) (string) {
    key := data[0:32]
    timeout, _ := strconv.Atoi(data[32:64])

    value := data[64:]
    (*dataDict).Set(key, value, timeout)
    return "key:"+key+"value:"+value
}
