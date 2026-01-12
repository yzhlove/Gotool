package main

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"strconv"
)

type Reader interface {
	Read(reader io.Reader) error
}

type Writer interface {
	Write(writer io.Writer) error
}

const (
	SocksVer            = 0x05
	UserVer             = 0x01
	NoAuthMethod        = 0x00
	UserMethod          = 0x02
	AuthSucceed         = 0x00
	AuthFailed          = 0x01
	NoAcceptableMethods = 0xFF
	AddrIpv4            = 0x01
	AddrIpv6            = 0x04
	AddrDomain          = 0x03
)

/*
X’00’ succeeded
X’01’ general SOCKS server failure
X’02’ connection not allowed by ruleset
X’03’ Network unreachable
X’04’ Host unreachable
X’05’ Connection refused
X’06’ TTL expired
X’07’ Command not supported
X’08’ Address type not supported
X’09’ to X’FF’ unassigned
*/

const (
	RepSucceeded       = 0x00
	RepFailure         = 0x01
	RepAllowed         = 0x02
	RepNetUnreachable  = 0x03
	RepHostUnreachable = 0x04
	RepConnRefused     = 0x05
	RepTTLExpired      = 0x06
	RepCmdUnsupported  = 0x07
	RepAddrUnsupported = 0x08
)

var (
	errBadVersion  = errors.New("Bad socks version! ")
	errBadMethod   = errors.New("Bad method! ")
	errBadUser     = errors.New("Bad user! ")
	errBadPassword = errors.New("Bad password! ")
	errBadNetwork  = errors.New("Bad network! ")
)

type MethodsReq struct {
	Methods []byte
}

func (opt *MethodsReq) Read(reader io.Reader) error {
	buf := GetBytes()
	defer PutBytes(buf)

	if _, err := io.ReadFull(reader, buf[:2]); err != nil {
		return err
	}

	ver, nMethods := int(buf[0]), int(buf[1])
	if ver != SocksVer {
		return errBadVersion
	}

	if nMethods == 0 {
		return errBadMethod
	}

	if _, err := io.ReadFull(reader, buf[:nMethods]); err != nil {
		return err
	}

	opt.Methods = make([]byte, nMethods)
	copy(opt.Methods, buf[:nMethods])
	return nil
}

type MethodResp struct {
	Method byte
}

func (opt *MethodResp) Write(writer io.Writer) error {
	_, err := writer.Write([]byte{SocksVer, opt.Method})
	return err
}

type MethodsResp struct {
	Methods []byte
}

func (opt *MethodsResp) Write(writer io.Writer) error {
	data := make([]byte, 2+len(opt.Methods))
	data[0] = SocksVer
	data[1] = byte(len(opt.Methods))
	copy(data[2:], opt.Methods)
	_, err := writer.Write(data)
	return err
}

type UserAuthReq struct {
	Version  byte
	Username string
	Password string
}

func (opt *UserAuthReq) Read(reader io.Reader) error {

	buf := GetBytes()
	defer PutBytes(buf)

	if _, err := io.ReadFull(reader, buf[:2]); err != nil {
		return err
	}

	userVer, userLen := int(buf[0]), int(buf[1])
	if userVer != UserVer {
		return errBadVersion
	}
	if userLen == 0 {
		return errBadUser
	}

	opt.Version = byte(userVer)
	if _, err := io.ReadFull(reader, buf[:userLen]); err != nil {
		return err
	}
	opt.Username = string(buf[:userLen])

	if _, err := io.ReadFull(reader, buf[:1]); err != nil {
		return err
	}

	passwordLen := int(buf[0])
	if passwordLen == 0 {
		return errBadPassword
	}

	if _, err := io.ReadFull(reader, buf[:passwordLen]); err != nil {
		return err
	}
	opt.Password = string(buf[:passwordLen])
	return nil
}

func (opt *UserAuthReq) Write(writer io.Writer) (err error) {

	data := make([]byte, 0, 2+len(opt.Username)+1+len(opt.Password))
	data = append(data, opt.Version)
	data = append(data, byte(len(opt.Username)))
	data = append(data, opt.Username...)
	data = append(data, byte(len(opt.Password)))
	data = append(data, opt.Password...)
	_, err = writer.Write(data)
	return
}

type UserAuthResp struct {
	Version byte
	Status  byte
}

func (opt *UserAuthResp) Write(writer io.Writer) error {
	_, err := writer.Write([]byte{opt.Version, opt.Status})
	return err
}

type Address struct {
	Type byte
	Host string
	Port uint16
}

func (a *Address) String() string {
	return net.JoinHostPort(a.Host, strconv.Itoa(int(a.Port)))
}

func (a *Address) Encode(b []byte) (pos int, err error) {
	b[0] = a.Type
	pos = 1
	switch a.Type {
	case AddrIpv4:
		ip4 := net.ParseIP(a.Host).To4()
		if ip4 == nil {
			ip4 = net.IPv4zero.To4()
		}
		pos += copy(b[pos:], ip4)
	case AddrIpv6:
		ip16 := net.ParseIP(a.Host).To16()
		if ip16 == nil {
			ip16 = net.IPv6zero.To16()
		}
		pos += copy(b[pos:], ip16)
	case AddrDomain:
		b[pos] = byte(len(a.Host))
		pos++
		pos += copy(b[pos:], a.Host)
	default:
		b[0] = AddrIpv4
		pos += copy(b[pos:pos+4], net.IPv4zero.To4())
	}
	binary.BigEndian.PutUint16(b[pos:pos+2], a.Port)
	pos += 2
	return
}

func (a *Address) Length() (n int) {
	switch a.Type {
	case AddrIpv4:
		n = 4 + net.IPv4len + 2
	case AddrIpv6:
		n = 4 + net.IPv6len + 2
	case AddrDomain:
		n = 4 + 1 + len(a.Host) + 2
	default:
		n = 4 + net.IPv4len + 2
	}
	return
}

func (a *Address) Decode(b []byte) (err error) {
	a.Type = b[0]
	pos := 1
	switch a.Type {
	case AddrIpv4:
		a.Host = net.IPv4(b[pos], b[pos+1], b[pos+2], b[pos+3]).String()
		pos += net.IPv4len
	case AddrIpv6:
		a.Host = net.IP(b[pos : pos+net.IPv6len]).String()
		pos += net.IPv6len
	case AddrDomain:
		domainLen := int(b[pos])
		pos++
		a.Host = string(b[pos : pos+domainLen])
		pos += domainLen
	default:
		return errBadNetwork
	}
	a.Port = binary.BigEndian.Uint16(b[pos : pos+2])
	return nil
}

type ConnectReq struct {
	Cmd  byte
	Addr *Address
}

func (opt *ConnectReq) Read(reader io.Reader) error {

	buf := GetBytes()
	defer PutBytes(buf)

	if _, err := io.ReadFull(reader, buf[:5]); err != nil {
		return err
	}

	if int(buf[0]) != SocksVer {
		return errBadVersion
	}
	opt.Cmd = buf[1]
	var pos = 5
	switch int(buf[3]) {
	case AddrIpv4:
		pos += net.IPv4len - 1
	case AddrIpv6:
		pos += net.IPv6len - 1
	case AddrDomain:
		pos += int(buf[4])
	default:
		return errBadNetwork
	}

	// 还需要解析 port，所以必须再加2个 byte
	pos += 2

	if _, err := io.ReadFull(reader, buf[5:pos]); err != nil {
		return err
	}

	opt.Addr = new(Address)
	// 这里必须包含 AType , 所以从 3 开始
	return opt.Addr.Decode(buf[3:pos])
}

type AddressResp struct {
	Req  byte
	Addr *Address
}

func NewAddressResp(req byte, addr *Address) *AddressResp {
	return &AddressResp{Req: req, Addr: addr}
}

func (opt *AddressResp) Write(writer io.Writer) (err error) {

	buf := GetBytes()
	defer PutBytes(buf)

	buf[0] = SocksVer
	buf[1] = opt.Req
	buf[2] = 0 // rsv
	buf[3] = AddrIpv4
	length := 10
	buf[4], buf[5], buf[6], buf[7], buf[8], buf[9] = 0, 0, 0, 0, 0, 0 // ipv4len + port

	if opt.Addr != nil {
		n, err := opt.Addr.Encode(buf[3:])
		if err != nil {
			return err
		}
		length = n + 3
	}

	_, err = writer.Write(buf[:length])
	return
}
