package main

import "testing"

var tests = []struct {
	name    string
	wantErr bool
}{
    {
        name: "Test with no error",
        wantErr: false,
    },
}

func TestRun(t *testing.T) {
	for _, e := range tests {
		t.Run(e.name, func(t *testing.T) {
			if _, err := run(); (err != nil) != e.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, e.wantErr)
			}
		})
	}
}
