package clashapi

import "testing"

func TestValidateExternalController(t *testing.T) {
	tests := []struct {
		name    string
		address string
		secret  string
		wantErr bool
	}{
		{name: "disabled", address: ""},
		{name: "loopback without secret", address: "127.0.0.1:9090"},
		{name: "ipv6 loopback without secret", address: "[::1]:9090"},
		{name: "wildcard without secret", address: "0.0.0.0:9090", wantErr: true},
		{name: "lan with short secret", address: "192.168.1.2:9090", secret: "short", wantErr: true},
		{name: "lan with secret", address: "192.168.1.2:9090", secret: "0123456789abcdef"},
		{name: "invalid address", address: "not-an-address", wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateExternalController(test.address, test.secret)
			if (err != nil) != test.wantErr {
				t.Fatalf("validateExternalController() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
