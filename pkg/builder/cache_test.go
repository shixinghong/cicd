package builder

import (
	"fmt"
	"reflect"
	"testing"
)

func TestInitCache(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "case1",
			args: args{size: 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitCache(tt.args.size)
			img1, _ := parseImage("nginx")
			img2, _ := parseImage("alpine")
			img3, _ := parseImage("mysql")
			ImageCache.Add(img1.Digest, img1)
			ImageCache.Add(img2.Digest, img2)
			ImageCache.Add(img3.Digest, img3)
			for _, k := range ImageCache.Keys() {
				v, _ := ImageCache.Get(k)
				fmt.Println(reflect.TypeOf(v))
				fmt.Println(v.(*Image).Name)

			}
		})
	}
}
