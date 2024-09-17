package project

import (
	"testing"

	"github.com/xhd2015/less-gen/go/load"
)

func TestGen(t *testing.T) {
	type args struct {
		args []string
		opts *load.LoadOptions
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.args.args, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Gen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_ = got
			// if got != tt.want {
			// 	t.Errorf("Gen() = %v, want %v", got, tt.want)
			// }
		})
	}
}
