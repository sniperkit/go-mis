package email

import "testing"

func TestSendEmailVerificationFailed(t *testing.T) {
	type args struct {
		email   string
		name    string
		reasons []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Email success sent",
			args: args{
				email: "didi.yudha@amartha.com",
				name:  "Didi Yudha Perwira",
				reasons: []string{
					"Nomor Rekening atau nama bank tidak valid",
					"Nama rekening pemilik bank tidak sesuai dengan nama registrasi",
					"Scan foto KTP tidak valid",
					"Resolusi gambar/scan foto KTP tidak jelas",
					"Nama registrasi berbeda dengan nama pada KTP",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendEmailVerificationFailed(tt.args.email, tt.args.name, tt.args.reasons)
		})
	}
}
