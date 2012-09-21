package data

import (
    "time"
)

type DataPacket struct {
    timeout int64
    data string
}

func (p *DataPacket) ShouldEvict() (bool) {
    self := *p
    now := time.Now().Unix()
    if self.timeout < now && self.data != "" {
        return true
    }
    return false
}

type DataDict struct {
    data map[string] *DataPacket
    lock DataLock
}

func NewDataDict() (*DataDict) {
    self := *new(DataDict)
    self.data = make(map[string] *DataPacket)
    self.lock = *NewDataLock(5)
    return &self
}

func (p *DataDict) Set(key, value string, timeout int) {
    self := *p
    if timeout >=0 {
        go self.Evicter(timeout, key)
    }

    now := time.Now().Unix()
    ttl := now + int64(timeout)

    var dp DataPacket
    dp.timeout = ttl
    dp.data = value

    self.lock.GetWrite()
    self.data[key] = &dp
    self.lock.ReleaseWrite()
}

func (p *DataDict) Get(key string) (string){
    self := *p

    self.lock.GetRead()
    dp := self.data[key]
    self.lock.ReleaseRead()

    if dp.ShouldEvict() {
        go self.Evicter(0, key)
        return ""
    }
    return dp.data
}

func (p *DataDict) Evicter(timeout int, key string){
    self := *p
    <- time.After(time.Duration(timeout) * time.Second)
    println("Autoevicting ", key)
    self.Set(key, "", -1)
}