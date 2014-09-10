package protocol

import (
	"testing"
)

func Test_CheckUDPFlag(t *testing.T) {
	if CheckUDPFlag(0x01) != true {
		t.Errorf("CheckUDPFlag fail")
	}

	if CheckUDPFlag(0x00) != false {
		t.Errorf("CheckUDPFlag fail")
	}
}

func Test_AddressTypeFlag(t *testing.T) {
	if CheckDomainFlag(0xFD) != true {
		t.Errorf("CheckDomainFlag error")
	}

	if CheckDomainFlag(0xFF) != false {
		t.Errorf("CheckDomainFlag error")
	}

	if CheckIPv4Flag(0xF9) != true {
		t.Errorf("CheckIPv4Flag error")
	}

	if CheckIPv6Flag(0xFB) != true {
		t.Errorf("CheckIPv6Flag error")
	}
}
