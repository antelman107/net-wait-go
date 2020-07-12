package main

import (
	"encoding/base64"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getFlags(t *testing.T) {
	var tests = []struct {
		name  string
		args  []string
		flags *flags
		err   error
	}{
		{"success",
			[]string{
				"-proto", "udp", "-addrs", "46.174.53.245:27015,185.158.113.136:27015",
				"-packet", "/////1RTb3VyY2UgRW5naW5lIFF1ZXJ5AA==", "-debug", "true"},
			&flags{proto: "udp", addrs: "46.174.53.245:27015,185.158.113.136:27015",
				addrsSlice:   []string{"46.174.53.245:27015", "185.158.113.136:27015"},
				packetBase64: "/////1RTb3VyY2UgRW5naW5lIFF1ZXJ5AA==", packetBytes: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x54, 0x53, 0x6F, 0x75, 0x72, 0x63, 0x65, 0x20, 0x45, 0x6E, 0x67, 0x69, 0x6E, 0x65, 0x20, 0x51, 0x75, 0x65, 0x72, 0x79, 0x00},
				debug:      true,
				deadlineMS: 10000,
				delayMS:    100,
				breakMS:    50,
			},
			nil,
		},

		{"empty addrs",
			[]string{
				"-proto", "udp",
				"-packet", "/////1RTb3VyY2UgRW5naW5lIFF1ZXJ5AA==", "-debug", "true"},
			nil,
			errors.New("addrs are not set"),
		},

		{"invalid base64 packet",
			[]string{
				"-proto", "udp", "-addrs", "46.174.53.245:27015,185.158.113.136:27015",
				"-packet", "/////1RTb3VyY2UgRW5naW5lIFF1ZXJ5AA==111", "-debug", "true"},
			nil,
			base64.CorruptInputError(36),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualFlags, _, err := getFlags("prog", tt.args)
			assert.Equal(t, tt.flags, actualFlags, "flags are not equal")
			assert.Equal(t, tt.err, err, "errors are not equal")
		})
	}
}
