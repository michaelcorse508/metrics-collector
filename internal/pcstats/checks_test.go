package pcstats

import "testing"

func TestCheckNameIsNotEmpty(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid name",
			args: args{
				name: "test",
			},
			wantErr: false,
		},
		{
			name: "invalid name",
			args: args{
				name: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckIDIsNotEmpty(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("CheckIDIsNotEmpty() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckTypeIsValid(t *testing.T) {
	type args struct {
		mType MetricType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid type gauge",
			args: args{
				mType: Gauge,
			},
			wantErr: false,
		},
		{
			name: "valid type counter",
			args: args{
				mType: Counter,
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			args: args{
				mType: MetricType("random"),
			},
			wantErr: true,
		},
		{
			name: "empty type",
			args: args{
				mType: MetricType(""),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckTypeIsValid(tt.args.mType); (err != nil) != tt.wantErr {
				t.Errorf("CheckTypeIsValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckGaugeIsNotNil(t *testing.T) {
	type args struct {
		value *float64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "not nil",
			args: args{
				value: new(float64),
			},
			wantErr: false,
		},
		{
			name: "nil",
			args: args{
				value: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckGaugeIsNotNil(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("CheckGaugeIsNotNil() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckCounterIsNotNil(t *testing.T) {
	type args struct {
		delta *int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "not nil",
			args: args{
				delta: new(int64),
			},
			wantErr: false,
		},
		{
			name: "nil",
			args: args{
				delta: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckCounterIsNotNil(tt.args.delta); (err != nil) != tt.wantErr {
				t.Errorf("CheckCounterIsNotNil() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
