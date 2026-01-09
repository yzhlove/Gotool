package main

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
)

type SocksReadOpt interface {
	Read(reader io.Reader) error
}

type SocksWriteOpt interface {
	Replay(writer io.Writer) error
}

const (
	SocksVer            = 0x05
	UserVer             = 0x01
	NoAuthMethod        = 0x00
	UsernameMethod      = 0x02
	NoAcceptableMethods = 0xFF
	AddrIpv4            = 0x01
	AddrIpv6            = 0x04
	AddrDomain          = 0x03
)

var (
	errBadVersion  = errors.New("socks version not supported! ")
	errBadMethod   = errors.New("bad method! ")
	errBadUser     = errors.New("bad user! ")
	errBadPassword = errors.New("bad password! ")
	errBadNetwork  = errors.New("bad network! ")
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

func (opt *MethodResp) Replay(writer io.Writer) error {
	_, err := writer.Write([]byte{SocksVer, opt.Method})
	return err
}

type MethodsResp struct {
	Methods []byte
}

func (opt *MethodsResp) Replay(writer io.Writer) error {
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

type UserAuthResp struct {
	Version byte
	Status  byte
}

func (opt *UserAuthResp) Replay(writer io.Writer) error {
	_, err := writer.Write([]byte{opt.Version, opt.Status})
	return err
}

type AddrReq struct {
	Cmd  byte
	Host string
	Port uint16
}

func (opt *AddrReq) Read(reader io.Reader) error {

	buf := GetBytes()
	defer PutBytes(buf)

	if _, err := io.ReadFull(reader, buf[:4]); err != nil {
		return err
	}

	if int(buf[0]) != SocksVer {
		return errBadVersion
	}
	opt.Cmd = buf[1]

	switch buf[3] {
	case AddrIpv4:
		if _, err := io.ReadFull(reader, buf[:net.IPv4len]); err != nil {
			return err
		}
		opt.Host = net.IPv4(buf[0], buf[1], buf[2], buf[3]).String()
	case AddrIpv6:
		if _, err := io.ReadFull(reader, buf[:net.IPv6len]); err != nil {
			return err
		}
		opt.Host = net.IP(buf[:net.IPv6len]).String()
	case AddrDomain:
		if _, err := io.ReadFull(reader, buf[:1]); err != nil {
			return err
		}
		domainLen := int(buf[0])
		if _, err := io.ReadFull(reader, buf[:domainLen]); err != nil {
			return err
		}
		opt.Host = string(buf[:domainLen])
	default:
		return errBadNetwork
	}

	if _, err := io.ReadFull(reader, buf[:2]); err != nil {
		return err
	}
	opt.Port = binary.BigEndian.Uint16(buf[:2])
	return nil
}
