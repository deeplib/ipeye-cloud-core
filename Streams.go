package icc

import (
	"errors"
	"net"
	"sync"
	"time"
)

//Streams global
//var Streams = NewStreams()

//StreamStorage struct
type StreamStorage struct {
	socket  net.Conn
	control chan bool
	tx      chan []byte
}

//StreamsT struct
type StreamsT struct {
	mutex sync.RWMutex
	m     map[string]StreamStorage
}

//NewStreams func
func NewStreams() *StreamsT {
	return &StreamsT{m: make(map[string]StreamStorage)}
}

//New Tunel chanel
func (element *StreamsT) New(key string, conn net.Conn, out chan *Packet) {
	//log.Println("tunel", "new tunel", key)
	element.mutex.Lock()
	defer element.mutex.Unlock()
	tx := make(chan []byte, 300)
	control := make(chan bool, 10)
	element.m[key] = StreamStorage{socket: conn, control: control, tx: tx}
	go element.Writer(key, conn, tx, control)
	go element.Reader(key, conn, out, control)
}

//Writer func
func (element *StreamsT) Writer(key string, conn net.Conn, tx chan []byte, control chan bool) {
	defer func() {
		//log.Println("tunel", "writer exit", key)
		conn.Close()
		element.Close(key)
	}()
	for {
		err := conn.SetDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			return
		}
		select {
		case msg := <-tx:
			_, err := conn.Write(msg)
			if err != nil {
				return
			}
		case <-control:
			return
		}
	}
}

//Reader func
func (element *StreamsT) Reader(key string, conn net.Conn, out chan *Packet, control chan bool) {
	defer func() {
		//log.Println("tunel", "reader exit", key)
		control <- true
		element.Close(key)
	}()
	ntmp := make([]byte, 65535)
	for {
		err := conn.SetDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			return
		}
		n, err := conn.Read(ntmp)
		if err != nil {
			return
		}
		stmp := make([]byte, n)
		copy(stmp, ntmp[:n])
		if len(out) > MaxLenPacketChanel {
			return
		}
		out <- &Packet{PackageType: DATA, Payload: stmp, TunelUUID: key}
	}
}

//Close close tunel
func (element *StreamsT) Close(key string) error {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	if _, ok := element.m[key]; ok {
		//log.Println("tunel", "close tunel", key)
		element.m[key].socket.Close()
		element.m[key].control <- true
		delete(element.m, key)
		return nil
	}
	return errors.New("tunel not found")
}

//Запись данных в тунель
func (element *StreamsT) Write(key string, val []byte) error {
	element.mutex.RLock()
	defer element.mutex.RUnlock()
	if tmp, ok := element.m[key]; ok {
		if len(tmp.tx) < MaxLenPacketChanel {
			tmp.tx <- val
			return nil
		}
		return errors.New("tunel " + key + "write error chanel full")
	}
	return errors.New("tunel " + key + "write error tunel not found")
}

/*
//Create for you bro )
func (element *StreamsT) Create(key string, conn net.Conn, out chan *Packet) {
	tx := make(chan []byte, 200)
	control := make(chan bool, 10)
	element.mutex.Lock()
	element.m[key] = StreamStorage{socket: conn, control: control, tx: tx}
	element.mutex.Unlock()
	go func() {
		defer func() {
			//log.Println("TUNEL DELETED", key)
			element.Delete(key)
		}()
		for {
			select {
			case msg := <-tx:
				_, err := conn.Write(msg)
				if err != nil {
					return
				}
			case <-control:
				return
			}
		}
	}()
	go func() {
		defer func() {
			control <- true
		}()
		ntmp := make([]byte, 65535)
		for {
			n, err := conn.Read(ntmp)
			if err != nil {
				return
			}
			stmp := make([]byte, n)
			copy(stmp, ntmp[:n])
			if len(out) > MaxLenPacketChanel {
				return
			}
			out <- &Packet{PackageType: DATA, Payload: stmp, TunelUUID: key}
		}
	}()
}

//Delete func
func (element *StreamsT) Delete(key string) error {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	if tmp, ok := element.m[key]; ok {
		tmp.socket.Close()
		delete(element.m, key)
		return nil
	}
	return errors.New("no found tunel")
}

//GetAllR func
func (element *StreamsT) GetAllR() map[string]StreamStorage {
	element.mutex.RLock()
	defer element.mutex.RUnlock()
	return element.m
}

//GetAll func
func (element *StreamsT) GetAll() []byte {
	element.mutex.RLock()
	defer element.mutex.RUnlock()
	b, _ := json.MarshalIndent(element.m, "", "  ")
	return b
}

//Write wunc
func (element *StreamsT) Write(key string, val []byte) error {
	element.mutex.RLock()
	defer element.mutex.RUnlock()
	if tmp, ok := element.m[key]; ok {
		if len(tmp.tx) < MaxLenPacketChanel {
			tmp.tx <- val
			return nil
		}
	}
	return errors.New("tunel " + key + "write error")
}
*/
