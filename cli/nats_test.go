package cli

import "testing"

func TestBuildConnectString(t *testing.T) {
	type args struct {
		urlStr   string
		user     string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "host only",
			args: args{
				urlStr:   "nats.server.com",
				user:     "",
				password: "",
			},
			want:    "nats://nats.server.com",
			wantErr: false,
		},
		{
			name: "host:port",
			args: args{
				urlStr:   "nats.server.com:1234",
				user:     "",
				password: "",
			},
			want:    "nats://nats.server.com:1234",
			wantErr: false,
		},
		{
			name: "scheme host:port",
			args: args{
				urlStr:   "wss://nats.server.com:1234",
				user:     "",
				password: "",
			},
			want:    "wss://nats.server.com:1234",
			wantErr: false,
		},
		{
			name: "scheme user pass host port",
			args: args{
				urlStr:   "wss://nats.server.com:443",
				user:     "user",
				password: "pass",
			},
			want:    "wss://user:pass@nats.server.com:443",
			wantErr: false,
		},
		{
			name: "scheme user host port",
			args: args{
				urlStr:   "wss://user@nats.server.com:443",
				user:     "",
				password: "",
			},
			want:    "wss://user@nats.server.com:443",
			wantErr: false,
		},
		{
			name: "pass only",
			args: args{
				urlStr:   "wss://user@nats.server.com:443",
				user:     "",
				password: "pass",
			},
			want:    "wss://pass@nats.server.com:443",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildConnectString(tt.args.urlStr, tt.args.user, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildConnectString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BuildConnectString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
