package main

import "testing"

func TestValidateDevelopmentListenAddress(t *testing.T) {
	for _, test := range []struct {
		address string
		wantErr bool
	}{
		{address: "127.0.0.1:8080"},
		{address: "[::1]:8080"},
		{address: "0.0.0.0:8080", wantErr: true},
		{address: ":8080", wantErr: true},
		{address: "192.168.1.2:8080", wantErr: true},
	} {
		t.Run(test.address, func(t *testing.T) {
			err := validateDevelopmentListenAddress(test.address)
			if (err != nil) != test.wantErr {
				t.Fatalf("validateDevelopmentListenAddress() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
