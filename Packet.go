package icc

import (
	"time"
)

//update
//Packet описание типов комманд сервера - клиента
const (
	PING = iota
	PONG
	REGISTER
	COMMAND_QUERY
	COMMAND_RESPONCE
	TCP_QUERY
	TCP_RESPONCE
	HTTP_QUERY
	HTTP_RESPONCE
	STREAM_CREATE
	STREAM_STOP
	DATA
	CRYPT
	EXEC_QUERY
	EXEC_RESPONCE
)

//Protocol описание типов протоколов нужно для серверной части
const (
	RTSP = iota
	HTTP
	HTTPS
)

//Type тип устройства трансляции
const (
	SERVER = iota
	CAMERA
	NVR
	PC
	CUSTOM
	OTHER
)

//MaxTCPPacketSize максимальный размер пакета TCP
var MaxTCPPacketSize = 65535

//MaxLenPacketChanel Максимальная длинна очереди собщений любой из буферизируемых очередений
//Требудетпровести тесты памяти что бы понять на что способна камера.
var MaxLenPacketChanel = 200
var PreMaxLenPacketChanel = 20
var MaxLenPacketChanelServer = 9000

type Options struct {
	HostPort     string
	ConnDeadline time.Duration
	RWDeadline   time.Duration
}

//Packet стандартная структура данных передаваямая между сервером и клиентом
type Packet struct {
	PackageType    int     //Тип пакета данных
	Payload        []byte  //Полезные данные тело пакета
	TunelUUID      string  //Уникальный иденитификатор тунеля
	TunelOptions   string  //Опции тунеля
	TunelError     string  //Опции тунеля
	TunelOptionsV2 Options //опции тунеля
}

//Register строка регистрации между сервером и клиентом
//Возможно потребуется передача описания uri для протокола rtsp
type Register struct {
	DeviceUUID              string                         `json:"DeviceUUID"` //Уникальный индентификатор устройства
	ServerPort              string                         `json:"-"`          //Сервер и порт подключения клиента
	ServiceIP               string                         `json:"ServiceIP"`
	ServiceAlarm            bool                           `json:"ServiceAlarm,omitempty"`            //Включить коллектор тревог
	ServiceRTSPProxy        bool                           `json:"SeviceRTSPProxy,omitempty"`         //Включить встроенный proxy
	ServiceRTSPPort         string                         `json:"SeviceRTSPPort,omitempty"`          //Исходящий порт proxy
	ServiceRTSPDelay        int                            `json:"ServiceRTSPDelay,omitempty"`        //Задержка proxy frame num
	ManufacturerDescription *ManufacturerDescriptionStruct `json:"ManufacturerDescription,omitempty"` //Описание производителя устройства
	DeviceDescription       *DeviceDescriptionStruct       `json:"DeviceDescription,omitempty"`       //Опасание устройства
	ClientDescription       *ClientDescriptionStruct       `json:"ClientDescription,omitempty"`       //Описание облачного клиента
	ResellerDescription     *ResellerDescriptionStruct     `json:"ResellerDescription,omitempty"`     //Описание компании интегратора
}

//ManufacturerDescriptionStruct Описание производителя
type ManufacturerDescriptionStruct struct {
	Title string `json:"Title,omitempty"`
}

//DeviceDescriptionStruct описание устройства
type DeviceDescriptionStruct struct {
	Title        string      `json:"Title,omitempty"`              //Название устройства
	Type         int         `json:"Type,omitempty"`               //Тип устройства
	Model        string      `json:"Model,omitempty"`              //Модель устройства
	SerialNumber string      `json:"SerialNumber,omitempty"`       //Серийный номер устройства
	IPAdress     string      `json:"IPAdress,omitempty"`           //IP адрес устройства
	IPAdresses   [][2]string `json:"IPAdresses,omitempty"`         //IP адреса устройства
	MACAdress    string      `json:"MACAdress,omitempty"`          //Апаратный адрес устройства
	Firmware     string      `json:"Firmware,omitempty,omitempty"` //Версия ПО устройства
	Version      string      `json:"Version,omitempty,omitempty"`  //Версия устройства
	Login        string      `json:"Login,omitempty,omitempty"`    //Логин от устройства
	Password     string      `json:"Password,omitempty"`           //Пароль от устройства
	PortRTSP     string      `json:"PortRTSP,omitempty"`           //порт rtsp устройства
	PortHTTP     string      `json:"PortHTTP,omitempty"`           //порт http сервера
	RTSPUri      []string    `json:"RTSPUri,omitempty"`            //список потдерживамых потоков
}

//ClientDescriptionStruct строка описания клиента доставки потока
type ClientDescriptionStruct struct {
	Language string `json:"Language,omitempty"` //Язык клиента
	Version  string `json:"Version,omitempty"`  //Версия клиента
}

//ResellerDescriptionStruct строка описания клиента доставки потока
type ResellerDescriptionStruct struct {
	Title string `json:"Title,omitempty"` //Название продовца
}
