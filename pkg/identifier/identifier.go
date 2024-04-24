package identifier

import (
	"crypto/sha256"
	"fmt"
	"net"

	"github.com/google/uuid"
)

type CTXKey int

const (
	DB CTXKey = iota + 1
	ReplicaID
)

func MacAddrHex() (string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	var as []byte
	for _, ifa := range ifas {
		a := ifa.HardwareAddr
		if len(a) > 0 {
			as = append(as, a...)
		}
	}
	// as to hex string
	h := sha256.New()
	h.Write(as)
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs), nil
}

func UUID() (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
