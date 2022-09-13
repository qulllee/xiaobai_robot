package ver

import "testing"

func TestGetURL(t *testing.T) {
	type args struct {
		endpoint uri
		Sandbox  bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1-获取非沙箱环境的SettingGuideURI的url",
			args: args{
				endpoint: SettingGuideURI,
				Sandbox:  false,
			},
			want: "https://api.sgroup.qq.com/channels/{channel_id}/settingguide",
		},
		{
			name: "test2-获取非沙箱环境的未定义字符的url",
			args: args{
				endpoint: "123",
				Sandbox:  false,
			},
			want: "https://api.sgroup.qq.com123",
		},
		{
			name: "test3-获取沙箱环境的SettingGuideURI的url",
			args: args{
				endpoint: SettingGuideURI,
				Sandbox:  true,
			},
			want: "https://sandbox.api.sgroup.qq.com/channels/{channel_id}/settingguide",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetURL(tt.args.endpoint, tt.args.Sandbox); got != tt.want {
				t.Errorf("GetURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
