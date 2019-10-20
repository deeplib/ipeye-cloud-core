package icc

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"sync"
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

//Create for you bro )
func (element *StreamsT) Create(key string, conn net.Conn, out chan *Packet) {
	tx := make(chan []byte, 200)
	control := make(chan bool, 10)
	element.mutex.Lock()
	element.m[key] = StreamStorage{socket: conn, control: control, tx: tx}
	element.mutex.Unlock()
	go func() {
		defer func() {
			log.Println("TUNEL DELETED", key)
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
func (element *StreamsT) Delete(key string) {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	if tmp, ok := element.m[key]; ok {
		tmp.socket.Close()
		delete(element.m, key)
	}
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
		}
	}
	return errors.New("Stream " + key + " Not Found")
}
