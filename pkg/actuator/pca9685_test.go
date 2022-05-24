package actuator

import (
	"go.uber.org/zap"
	"periph.io/x/conn/v3/physic"
	"testing"
)

func init() {
	l, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(l)
}

func Test_convertToDuty(t *testing.T) {
	type fields struct {
	}
	type args struct {
		percent   float32
		leftPWM   PWM
		rightPWM  PWM
		centerPWM PWM
		freq      physic.Frequency
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "center",
			args: args{
				percent:   0.,
				leftPWM:   1000,
				rightPWM:  2000,
				centerPWM: 1500,
				freq:      60 * physic.Hertz,
			},
			want: "9%",
		},
		{
			name: "left",
			args: args{
				percent:   -1.,
				leftPWM:   1000,
				rightPWM:  2000,
				centerPWM: 1500,
				freq:      60 * physic.Hertz,
			},
			want: "6%",
		},
		{
			name: "right",
			args: args{
				percent:   1.,
				leftPWM:   1000,
				rightPWM:  2000,
				centerPWM: 1500,
				freq:      60 * physic.Hertz,
			},
			want: "12%",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got, err := convertToDuty(tt.args.percent, tt.args.freq, tt.args.centerPWM, tt.args.leftPWM, tt.args.rightPWM); got.String() != tt.want {
				if tt.wantErr && err == nil {
					t.Errorf("an error is expected")
					return
				}
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				t.Errorf("convertToDuty() = %v, want %v", got, tt.want)
			}
		})
	}
}
