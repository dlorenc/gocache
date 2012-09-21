package data

type DataLock struct {
    writeLock chan int
    readLock chan int
    _numReaders int
}

func NewDataLock(numReaders int) (*DataLock){
    self := *new(DataLock)
    self.writeLock = make(chan int, 1)
    self.readLock = make(chan int, numReaders)
    self._numReaders = numReaders
    return &self
}

func (p *DataLock) GetRead() {
    self := *p
    self.readLock <- 1
}

func (p *DataLock) ReleaseRead() {
    self := *p
    <- self.readLock
}

func (p *DataLock) GetWrite() {
    self := *p
    self.writeLock <- 1 
    for i:=0; i < self._numReaders; i++ {
        self.GetRead()
    }
    <- self.writeLock
}

func (p *DataLock) ReleaseWrite() {
    self := *p
    for i:=0; i < self._numReaders; i++ {
        self.ReleaseRead()
    }
}