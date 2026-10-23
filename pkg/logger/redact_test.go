package logger

import "testing"

func TestRedactSecrets(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "json login body",
			in:   `{"email":"a@b.com","password":"hunter2"}`,
			want: `{"email":"a@b.com","password":"***"}`,
		},
		{
			name: "json pretty with confirm_password and token",
			in:   "{\n  \"password\": \"p1\",\n  \"confirm_password\": \"p1\",\n  \"token\": \"abc\"\n}",
			want: "{\n  \"password\": \"***\",\n  \"confirm_password\": \"***\",\n  \"token\": \"***\"\n}",
		},
		{
			name: "json escaped quote in value",
			in:   `{"Password":"a\"b","name":"x"}`,
			want: `{"Password":"***","name":"x"}`,
		},
		{
			name: "form body",
			in:   `email=a%40b.com&password=hunter2&token=xyz`,
			want: `email=a%40b.com&password=***&token=***`,
		},
		{
			name: "no secrets untouched",
			in:   `{"email":"a@b.com","name":"Budi"}`,
			want: `{"email":"a@b.com","name":"Budi"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := redactSecrets(tt.in); got != tt.want {
				t.Errorf("redactSecrets(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
