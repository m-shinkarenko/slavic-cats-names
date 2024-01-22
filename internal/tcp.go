package internal

import (
	"encoding/binary"
	"errors"
	"net"
)

func TCPRead(c net.Conn, maxContentLength uint32) ([]byte, error) {
	contentLengthBytes := make([]byte, 4)

	_, err := c.Read(contentLengthBytes)
	if err != nil {
		return nil, errors.Join(err, errors.New("error reading content length"))
	}

	contentLength := binary.BigEndian.Uint32(contentLengthBytes)
	if contentLength > maxContentLength {
		return nil, errors.New("the request is too big")
	}

	content := make([]byte, contentLength)
	_, err = c.Read(content)
	if err != nil {
		return nil, errors.Join(err, errors.New("error reading content"))
	}

	return content, nil
}

func TCPWrite(c net.Conn, content []byte) error {
	res := binary.BigEndian.AppendUint32(make([]byte, 0, len(content)+4), uint32(len(content)))

	_, err := c.Write(append(res, content...))
	if err != nil {
		return errors.Join(err, errors.New("error writing response"))
	}
	return nil
}
