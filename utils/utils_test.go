package utils

import (
	"fmt"
	"testing"
)

func TestNewCurl(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name     string
		args     args
		wantCurl *Curl
	}{
		{"test", args{"curl 'https://www.baidu.com/' -H 'User-Agent: Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:79.0) Gecko/20100101 Firefox/79.0' -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8' -H 'Accept-Language: zh-CN,en-US;q=0.7,en;q=0.3' --compressed -H 'Connection: keep-alive' -H 'Upgrade-Insecure-Requests: 1' -H 'Cache-Control: max-age=0' -H 'TE: Trailers'"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCurl := NewCurl(tt.args.data)
			fmt.Println(gotCurl)

		})
	}
}

func TestParseTheFileC(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name      string
		args      args
		wantCurls []*Curl
		wantErr   bool
	}{
		{"test", args{path: "curl.txt"},nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCurls, _ := ParseTheFileC(tt.args.path)
			for _,v := range gotCurls {
				fmt.Println(v)
			}
		})
	}
}