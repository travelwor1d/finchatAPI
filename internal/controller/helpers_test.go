package controller

import "testing"

func Test_getUserTypes(t *testing.T) {
	type args struct {
		t string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"#1", args{"user, goat"}, "'USER','GOAT'", false},
		{"#2", args{"goat"}, "'GOAT'", false},
		{"#3", args{"user"}, "'USER'", false},
		{"#4", args{""}, "'GOAT','USER'", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getUserTypes(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("getUserTypes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getUserTypes() = %v, want %v", got, tt.want)
			}
		})
	}
}
