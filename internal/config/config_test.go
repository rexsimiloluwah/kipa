package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig_GetEnv(t *testing.T) {
	type args struct {
		key          string
		defaultValue string
	}

	tt := []struct {
		name    string
		args    args
		stubFn  func()
		want    string
		wantErr bool
	}{
		{
			name: "should_return_default_value",
			args: args{
				key:          "name",
				defaultValue: "simi",
			},
			want: "simi",
		},
		{
			name: "should_return_default_value",
			args: args{
				key:          "key",
				defaultValue: "simi",
			},
			stubFn: func() {
				os.Setenv("key", "value")
			},
			want: "value",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn()
			}
			value := getEnv(tc.args.key, tc.args.defaultValue)
			require.Equal(t, value, tc.want)
		})
	}
}

func TestConfig_GetEnvAsInt(t *testing.T) {
	type args struct {
		key          string
		defaultValue int
	}

	tt := []struct {
		name    string
		args    args
		stubFn  func()
		want    int
		wantErr bool
	}{
		{
			name: "should_return_default_value",
			args: args{
				key:          "age",
				defaultValue: 12,
			},
			stubFn: nil,
			want:   12,
		},
		{
			name: "should_return_zero_value_invalid_int",
			args: args{
				key:          "key",
				defaultValue: 12,
			},
			stubFn: func() {
				os.Setenv("key", "oops")
			},
			want: 12,
		},
		{
			name: "should_return_valid_value",
			args: args{
				key:          "key",
				defaultValue: 12,
			},
			stubFn: func() {
				os.Setenv("key", "100")
			},
			want: 100,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn()
			}
			value := getEnvAsInt(tc.args.key, tc.args.defaultValue)
			require.Equal(t, value, tc.want)
			os.Clearenv()
		})
	}
}

func TestConfig_GetEnvAsBool(t *testing.T) {
	type args struct {
		key          string
		defaultValue bool
	}

	tt := []struct {
		name    string
		args    args
		stubFn  func()
		want    bool
		wantErr bool
	}{
		{
			name: "should_return_default_value",
			args: args{
				key:          "key",
				defaultValue: false,
			},
			want: false,
		},
		{
			name: "should_return_default_value_invalid_bool",
			args: args{
				key:          "key",
				defaultValue: true,
			},
			stubFn: func() {
				os.Setenv("key", "oops")
			},
			want: true,
		},
		{
			name: "should_return_valid_value",
			args: args{
				key:          "key",
				defaultValue: false,
			},
			stubFn: func() {
				os.Setenv("key", "true")
			},
			want: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn()
			}
			value := getEnvAsBool(tc.args.key, tc.args.defaultValue)
			require.Equal(t, value, tc.want)
			os.Clearenv()
		})
	}
}

func TestConfig_GetEnvAsFloat(t *testing.T) {
	type args struct {
		key          string
		defaultValue float64
	}

	tt := []struct {
		name    string
		args    args
		stubFn  func()
		want    float64
		wantErr bool
	}{
		{
			name: "should_return_default_value",
			args: args{
				key:          "key",
				defaultValue: 99.5,
			},
			want: 99.5,
		},
		{
			name: "should_return_default_value_invalid_float",
			args: args{
				key:          "key",
				defaultValue: 1.2,
			},
			stubFn: func() {
				os.Setenv("key", "oops")
			},
			want: 1.2,
		},
		{
			name: "should_return_valid_value",
			args: args{
				key:          "key",
				defaultValue: 1.2,
			},
			stubFn: func() {
				os.Setenv("key", "2.5")
			},
			want: 2.5,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn()
			}
			value := getEnvAsFloat(tc.args.key, tc.args.defaultValue)
			require.Equal(t, value, tc.want)
			os.Clearenv()
		})
	}
}

func TestConfig_GetEnvAsSlice(t *testing.T) {
	type args struct {
		key          string
		sep          string
		defaultValue []string
	}

	tt := []struct {
		name    string
		args    args
		stubFn  func()
		want    []string
		wantErr bool
	}{
		{
			name: "should_return_default_value",
			args: args{
				key:          "fruits",
				sep:          ",",
				defaultValue: []string{"orange", "mango"},
			},
			want: []string{"orange", "mango"},
		},
		{
			name: "should_return_valid_value",
			args: args{
				key:          "fruits",
				sep:          ",",
				defaultValue: []string{},
			},
			stubFn: func() {
				os.Setenv("fruits", "orange,mango")
			},
			want: []string{"orange,mango"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.stubFn != nil {
				tc.stubFn()
			}
			value := getEnvAsSlice(tc.args.key, tc.args.defaultValue, ".")
			require.Equal(t, value, tc.want)
			os.Clearenv()
		})
	}
}
