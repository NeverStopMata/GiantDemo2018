package common

import (
	"base/glog"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"io"
	"runtime"
	"strconv"
	"strings"

	"github.com/gogo/protobuf/proto"
)

const (
	MaxCompressSize = 4096
	CmdHeaderSize   = 2
	ServerCmdSize   = 1
	CmdRoleSize     = 8
	ServerIdSize    = 4
)

var (
	ShortMsgError = errors.New("[协议] 解码错误, 长度过短")
)

type Message interface {
	Marshal() (data []byte, err error)
	MarshalTo(data []byte) (n int, err error)
	Size() (n int)
	Unmarshal(data []byte) error
}

//zlib压缩，带缓冲，减少gc
type ZlibCompress struct {
	in bytes.Buffer
	w  *zlib.Writer
}

func (comp *ZlibCompress) Init() {
	comp.w = zlib.NewWriter(&comp.in)
}
func (comp *ZlibCompress) ZlibCompressWithBuff(src []byte) error {
	comp.in.Reset()
	comp.w.Reset(&comp.in)
	_, err := comp.w.Write(src)
	if err != nil {
		return err
	}
	comp.w.Close()
	return nil
}

func zlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	_, err := w.Write(src)
	if err != nil {
		return nil
	}
	w.Close()
	return in.Bytes()
}

func zlibUnCompress(src []byte) []byte {
	b := bytes.NewReader(src)
	var out bytes.Buffer
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil
	}
	_, err = io.Copy(&out, r)
	if err != nil {
		return nil
	}
	return out.Bytes()
}

// 生成二进制数据,返回数据和标识
func EncodeCmd(cmd uint16, msg proto.Message) ([]byte, byte, error) {
	data, err := proto.Marshal(msg)
	if err != nil {
		glog.Error("[协议] protobuf生成二进制数据错误 ", err)
		return nil, 0, err
	}
	var (
		mflag byte
		mbuff []byte
	)
	if len(data) >= MaxCompressSize {
		mflag = 1
		mbuff = zlibCompress(data)
	} else {
		mflag = 0
		mbuff = data
	}
	p := make([]byte, len(mbuff)+CmdHeaderSize)
	binary.LittleEndian.PutUint16(p[0:], cmd)
	copy(p[2:], mbuff)
	return p, mflag, nil
}

// 获取指令号
func GetCmd(buf []byte) uint16 {
	if len(buf) < CmdHeaderSize {
		return 0
	}
	return uint16(buf[0]) | uint16(buf[1])<<8
}

// 获取大指令号
func GetPrevCmd(buf []byte) uint8 {
	if len(buf) < 1 {
		return 0
	}
	return uint8(buf[0])
}

// 获取小指令号
func GetNextCmd(buf []byte) uint8 {
	if len(buf) < 2 {
		return 0
	}
	return uint8(buf[1])
}

// 大消息号
func GetBigCmd(buf []byte) uint8 {
	if len(buf) < 2 {
		return 0
	}
	return uint8(buf[1])
}

// 小消息号
func GetLittleCmd(buf []byte) uint8 {
	if len(buf) < 1 {
		return 0
	}
	return uint8(buf[0])
}

// 网关消息号
func GetGatewayCmd(buf []byte) uint8 {
	if len(buf) < 1 {
		return 0
	}
	return uint8(buf[0])
}

func GetPrevID(buf []byte) uint64 {
	if len(buf) < 9 {
		return 0
	}
	return GetUint64(buf[1:])
}

func GetServerID(buf []byte) uint32 {
	if len(buf) < 13 {
		return 0
	}
	return GetUint32(buf[9:])
}

func GetGateaySid(buf []byte) uint32 {
	if len(buf) < 13 {
		return 0
	}
	return uint32(buf[1]) | uint32(buf[2])<<8 | uint32(buf[3])<<8 | uint32(buf[4])<<8
}

func GetSidUid(buf []byte) (uint64, uint32) {
	if len(buf) < 13 {
		return 0, 0
	}
	return GetUint64(buf[1:]), uint32(buf[9]) | uint32(buf[10])<<8 | uint32(buf[11])<<8 | uint32(buf[12])<<8
}

func GetUint32(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func PutUint32(b []byte, v uint32) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
}

func GetUint64(b []byte) uint64 {
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func PutUint64(b []byte, v uint64) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
}

func PutMsgHead(msgNo uint8, v uint64) []byte {
	b := make([]byte, 9)
	b[0] = msgNo
	PutUint64(b[1:], v)
	return b
}

// 解析protobuf数据
// Todo(Jinq): 改为返回bool.
func DecodeCmd(buf []byte, flag byte, pb proto.Message) proto.Message {
	if len(buf) < CmdHeaderSize {
		glog.Error("[协议] 解析protobuf数据错误 ", buf)
		return nil
	}
	var mbuff []byte
	if flag == 1 {
		mbuff = zlibUnCompress(buf[CmdHeaderSize:])
	} else {
		mbuff = buf[CmdHeaderSize:]
	}
	err := proto.Unmarshal(mbuff, pb)
	if err != nil {
		glog.Error("[协议] 生成protobuf数据错误 ", err)
		return nil
	}
	return pb
}

// 生成二进制数据,返回数据和是否压缩标识
func EncodeGoCmdWithCache(zlibComp *ZlibCompress, cmd uint16, msg Message) (data []byte, flag byte, err error) {
	msglen := msg.Size()
	if msglen >= MaxCompressSize {
		flag = 1
		data, err = msg.Marshal()
		if err != nil {
			glog.Error("[协议] 生成gogo二进制数据编码错误 ", err)
			return
		}
		if errz := zlibComp.ZlibCompressWithBuff(data); errz != nil {
			glog.Error("[协议] ZlibCompressWithBuff fail:", errz)
			return
		}
		mbuff := zlibComp.in.Bytes()
		data = make([]byte, CmdHeaderSize+len(mbuff))
		data[0] = byte(cmd)
		data[1] = byte(cmd >> 8)
		copy(data[CmdHeaderSize:], mbuff)
		return
	}
	data = make([]byte, CmdHeaderSize+msglen)
	_, err = msg.MarshalTo(data[CmdHeaderSize:])
	if err != nil {
		glog.Error("[协议] 编码错误 ", err)
		return nil, 0, err
	}
	data[0] = byte(cmd)
	data[1] = byte(cmd >> 8)
	return
}

// 生成二进制数据,返回数据和是否压缩标识
func EncodeGoCmd(cmd uint16, msg Message) (data []byte, flag byte, err error) {
	return EncodeGoCmdWithCompressSize(cmd, msg, MaxCompressSize)
}

// 生成二进制数据,返回数据和是否压缩标识(指定开启压缩的消息长度)
func EncodeGoCmdWithCompressSize(cmd uint16, msg Message, compressSize int) (data []byte, flag byte, err error) {
	msglen := msg.Size()

	if msglen >= compressSize {
		data, err = msg.Marshal()
		if err != nil {
			glog.Error("[协议] 生成gogo二进制数据编码错误 ", err)
			return
		}
		mbuff := zlibCompress(data)
		mbufflen := len(mbuff)

		// 压缩后数据长度有减少，该压缩才有意义
		if mbufflen < msglen {
			flag = 1
			data = make([]byte, CmdHeaderSize+mbufflen)
			data[0] = byte(cmd)
			data[1] = byte(cmd >> 8)
			copy(data[CmdHeaderSize:], mbuff)
			return
		}
	}
	data = make([]byte, CmdHeaderSize+msglen)
	_, err = msg.MarshalTo(data[CmdHeaderSize:])
	if err != nil {
		glog.Error("[协议] 编码错误 ", err)
		return nil, 0, err
	}
	data[0] = byte(cmd)
	data[1] = byte(cmd >> 8)
	return
}

// protobuf解码数据
func DecodeGoCmd(buf []byte, flag byte, pb Message) (err error) {
	if len(buf) < CmdHeaderSize {
		err = ShortMsgError
		glog.Error(err.Error())
		return
	}
	var mbuff []byte
	if flag == 1 {
		mbuff = zlibUnCompress(buf[CmdHeaderSize:])
	} else {
		mbuff = buf[CmdHeaderSize:]
	}
	err = pb.Unmarshal(mbuff)
	if err != nil {
		glog.Error("[协议] gogo解码错误 ", err)
	}
	return
}

// 生成二进制数据,返回数据和是否压缩标识
func EncodeGatewayGoCmd(cmd uint32, msg Message) (data []byte, flag byte, err error) {
	msglen := msg.Size()
	if msglen >= MaxCompressSize {
		flag = 1
		data, err = msg.Marshal()
		if err != nil {
			glog.Error("[协议] 生成gogo二进制数据编码错误 ", err)
			return
		}
		mbuff := zlibCompress(data)
		data = make([]byte, ServerCmdSize+len(mbuff))
		data[0] = byte(cmd)
		copy(data[ServerCmdSize:], mbuff)
		return
	}
	data = make([]byte, ServerCmdSize+msglen)
	_, err = msg.MarshalTo(data[ServerCmdSize:])
	if err != nil {
		glog.Error("[协议] 编码错误 ", err)
		return nil, 0, err
	}
	data[0] = byte(cmd)
	return
}

// protobuf解码数据
func DecodeGatewayGoCmd(buf []byte, flag byte, pb Message) (err error) {
	if len(buf) < ServerCmdSize {
		err = ShortMsgError
		glog.Error(err.Error())
		return
	}
	var mbuff []byte
	if flag == 1 {
		mbuff = zlibUnCompress(buf[ServerCmdSize:])
	} else {
		mbuff = buf[ServerCmdSize:]
	}
	err = pb.Unmarshal(mbuff)
	if err != nil {
		glog.Error("[协议] gogo解码错误 ", err)
	}
	return
}

// 解析UDP数据
func DecodeUdpCmd(buf []byte, pb proto.Message) proto.Message {
	err := proto.Unmarshal(buf, pb)
	if err != nil {
		glog.Error("[协议] 生成protobuf数据错误 ", err)
		return nil
	}
	return pb
}

// 生成UDP二进制数据,返回数据
func EncodeUdpGoCmd(msg Message) (data []byte, err error) {
	msglen := msg.Size()
	data = make([]byte, 3+msglen)
	_, err = msg.MarshalTo(data[3:])
	if err != nil {
		glog.Error("[协议] 编码错误 ", err)
		return nil, err
	}
	return
}

// 获取玩家ID
func GetRoleID(buf []byte) uint64 {
	if len(buf) < (CmdHeaderSize + CmdRoleSize) {
		return 0
	}
	return uint64(buf[2]) | uint64(buf[3])<<8 | uint64(buf[4])<<16 | uint64(buf[5])<<24 | uint64(buf[6])<<32 | uint64(buf[7])<<40 | uint64(buf[8])<<48 | uint64(buf[9])<<56
}

func SetRoleID(buf []byte, userId uint64) bool {
	if len(buf) < (CmdHeaderSize + CmdRoleSize) {
		return false
	}
	copy(buf[CmdHeaderSize:], makeRoleByte(userId))
	return true
}

func makeRoleByte(userId uint64) []byte {
	var rolebyte []byte
	rolebyte = append(rolebyte, byte(userId), byte(userId>>8), byte(userId>>16), byte(userId>>24), byte(userId>>32), byte(userId>>40), byte(userId>>48), byte(userId>>56))
	return rolebyte
}

// 生成二进制数据,返回数据和标识
func LineByte2ChatByte(userId uint64, data []byte) ([]byte, error) {
	p := make([]byte, len(data)+CmdRoleSize)
	copy(p[0:], data[0:CmdHeaderSize])
	copy(p[CmdHeaderSize:], makeRoleByte(userId))
	copy(p[CmdHeaderSize+CmdRoleSize:], data[CmdHeaderSize:])
	return p, nil
}

//将地图发来二进制转化为网关二进制 区别在于里面是否带玩家ID
func ChatByte2LineByte(data []byte) ([]byte, error) {
	p := make([]byte, len(data)-CmdRoleSize)
	copy(p[0:], data[0:CmdHeaderSize])
	copy(p[CmdHeaderSize:], data[CmdHeaderSize+CmdRoleSize:])
	return p, nil
}

// 生成二进制数据,返回数据和标识
func Proto2Server(userId uint64, serverCmd uint8, cmd uint16, msg proto.Message) (data []byte, flag byte, err error) {
	data, err = proto.Marshal(msg)
	if err != nil {
		glog.Error("[协议] 编码错误 ", err)
		return nil, 0, err
	}
	var (
		mflag byte
		mbuff []byte
	)
	if len(data) >= MaxCompressSize {
		mflag = 1
		mbuff = zlibCompress(data)
	} else {
		mflag = 0
		mbuff = data
	}
	headLen := ServerCmdSize + CmdHeaderSize + CmdRoleSize

	p := make([]byte, headLen+len(mbuff))

	p[0] = byte(serverCmd)
	PutUint64(p[1:], userId)
	p[9] = byte(cmd)
	p[10] = byte(cmd >> 8)
	copy(p[headLen:], mbuff)
	return p, mflag, nil
}

// 生成二进制数据,返回数据和标识
func Proto2GoServer(userId uint64, serverCmd uint8, cmd uint16, msg Message) (data []byte, flag byte, err error) {
	headLen := ServerCmdSize + CmdHeaderSize + CmdRoleSize
	msglen := msg.Size()
	if msglen >= MaxCompressSize {
		flag = 1
		data, err = msg.Marshal()
		if err != nil {
			glog.Error("[协议] 生成gogo二进制数据编码错误 ", err)
			return
		}
		mbuff := zlibCompress(data)
		data = make([]byte, headLen+len(mbuff))
		data[0] = byte(serverCmd)
		PutUint64(data[1:], userId)
		data[9] = byte(cmd)
		data[10] = byte(cmd >> 8)
		copy(data[headLen:], mbuff)
		return data, 1, nil
	}
	data = make([]byte, headLen+msglen)
	_, err = msg.MarshalTo(data[headLen:])
	if err != nil {
		glog.Error("[协议] 编码错误 ", err)
		return
	}
	data[0] = byte(serverCmd)
	PutUint64(data[1:], userId)
	data[9] = byte(cmd)
	data[10] = byte(cmd >> 8)
	return
}

// 生成二进制数据,返回数据和标识
func Proto2ServerWithUserId(serverId uint32, userId uint64, serverCmd uint8, cmd uint16, msg proto.Message) (data []byte, flag byte, err error) {
	data, err = proto.Marshal(msg)
	if err != nil {
		glog.Error("[协议] 编码错误 ", err)
		return nil, 0, err
	}
	var (
		mflag byte
		mbuff []byte
	)
	if len(data) >= MaxCompressSize {
		mflag = 1
		mbuff = zlibCompress(data)
	} else {
		mflag = 0
		mbuff = data
	}
	headLen := ServerCmdSize + CmdHeaderSize + CmdRoleSize + ServerIdSize

	p := make([]byte, headLen+len(mbuff))

	p[0] = byte(serverCmd)
	PutUint64(p[1:], userId)
	PutUint32(p[9:], serverId)
	p[13] = byte(cmd)
	p[14] = byte(cmd >> 8)
	copy(p[headLen:], mbuff)
	return p, mflag, nil
}

// 生成二进制数据,返回数据和标识
func Proto2ChatByte(userId uint64, cmd uint16, msg proto.Message) ([]byte, byte, error) {
	data, err := proto.Marshal(msg)
	if err != nil {
		glog.Error("[协议] 编码错误 ", err)
		return nil, 0, err
	}
	var (
		mflag byte
		mbuff []byte
	)
	if len(data) >= MaxCompressSize {
		mflag = 1
		mbuff = zlibCompress(data)
	} else {
		mflag = 0
		mbuff = data
	}

	p := make([]byte, CmdHeaderSize+CmdRoleSize+len(mbuff))
	binary.LittleEndian.PutUint16(p[0:], cmd)
	copy(p[CmdHeaderSize:], makeRoleByte(userId))
	copy(p[CmdHeaderSize+CmdRoleSize:], mbuff)
	return p, mflag, nil
}

// 生成二进制数据,返回数据和标识
func Empty2ChatByte(userId uint64, cmd uint16) ([]byte, byte, error) {
	p := make([]byte, CmdHeaderSize+CmdRoleSize)
	binary.LittleEndian.PutUint16(p[0:], cmd)
	copy(p[CmdHeaderSize:], makeRoleByte(userId))
	return p, 0, nil
}

// 解析protobuf数据
func ChatByte2Proto(buf []byte, flag byte, pb proto.Message) proto.Message {
	if len(buf) < CmdHeaderSize+CmdRoleSize {
		glog.Error("[协议] 数据错误 ", buf)
		return nil
	}
	var mbuff []byte
	if flag == 1 {
		mbuff = zlibUnCompress(buf[CmdHeaderSize+CmdRoleSize:])
	} else {
		mbuff = buf[CmdHeaderSize+CmdRoleSize:]
	}

	err := proto.Unmarshal(mbuff, pb)
	if err != nil {
		_, file, line, ok := runtime.Caller(1)
		glog.Error("[协议] 解码错误 ", err, ",", mbuff, " ", ok, " ", file, " ", line)
		return nil
	}
	return pb
}

func ParseAttr(attr string) []*GiftItem {
	lastPos := strings.Index(attr, ",")
	if lastPos == -1 {
		return nil
	}
	gifts := make([]*GiftItem, 0, 1)
	sps := strings.Split(attr, "|")
	for _, item := range sps {
		if len(item) > 0 {
			giftStr := strings.Split(item, ",")
			if len(giftStr) == 2 {
				giftId, _ := strconv.Atoi(giftStr[0])
				giftNum, _ := strconv.Atoi(giftStr[1])
				gifts = append(gifts, &GiftItem{
					GiftId:  uint32(giftId),
					GiftNum: uint32(giftNum),
				})
			}
		}
	}
	return gifts
}
