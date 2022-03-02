package builder

import (
	"fmt"
	"testing"
)

func Test_parseImage(t *testing.T) {
	type args struct {
		img string
	}
	tests := []struct {
		name    string
		args    args
		want    *Image
		wantErr bool
	}{
		{
			name: "case1",
			args: args{img: "nginx"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseImage(tt.args.img)
			//if (err != nil) != tt.wantErr {
			//	t.Errorf("parseImage() error = %v, wantErr %v", err, tt.wantErr)
			//	return
			//}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("parseImage() got = %v, want %v", got, tt.want)
			//}
			if err != nil {
				t.Error(err)
			}
			fmt.Println(got)
		})
	}
}
