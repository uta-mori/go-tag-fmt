package main

import (
	"strings"
	"testing"
)

func TestAlign(t *testing.T) {
	type args struct {
		tags []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "single tag",
			args: args{tags: []string{
				`json:"id"`,
				`json:"data,omitempty"`,
			}},
			want: []string{
				`json:"id"`,
				`json:"data,omitempty"`,
			},
		},
		{
			name: "multi tag",
			args: args{tags: []string{
				`json:"id" xml:"id"`,
				`json:"data,omitempty" xml:"data"`,
			}},
			want: []string{
				`json:"id"             xml:"id"`,
				`json:"data,omitempty" xml:"data"`,
			},
		},
		{
			name: "gorm with json",
			args: args{tags: []string{
				`gorm:"primary,int"`,
				`json:"id" gorm:"int"`,
				`json:"data,omitempty"`,
			}},
			want: []string{
				`gorm:"primary,int"`,
				`gorm:"int"         json:"id"`,
				`                   json:"data,omitempty"`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T){
			tags := Align(tt.args.tags)
			got := strings.Join(tags,"\n")
			want := strings.Join(tt.want, "\n")
			if got != want {
				t.Errorf("err%v %v",got,want)
			}
		})
	}
}
