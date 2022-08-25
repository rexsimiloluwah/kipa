package config

import (
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
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
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
		want    int
		wantErr bool
	}{
		{
			name: "should_return_default_value",
			args: args{
				key:          "age",
				defaultValue: 12,
			},
			want: 12,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			value := getEnvAsInt(tc.args.key, tc.args.defaultValue)
			require.Equal(t, value, tc.want)
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
		want    bool
		wantErr bool
	}{
		{
			name: "should_return_default_value",
			args: args{
				key:          "age",
				defaultValue: false,
			},
			want: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			value := getEnvAsBool(tc.args.key, tc.args.defaultValue)
			require.Equal(t, value, tc.want)
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
		want    float64
		wantErr bool
	}{
		{
			name: "should_return_default_value",
			args: args{
				key:          "percentage",
				defaultValue: 99.5,
			},
			want: 99.5,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			value := getEnvAsFloat(tc.args.key, tc.args.defaultValue)
			require.Equal(t, value, tc.want)
		})
	}
}

func TestConfig_GetEnvAsSlice(t *testing.T) {
	type args struct {
		key          string
		defaultValue []string
	}

	tt := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "should_return_default_value",
			args: args{
				key:          "fruits",
				defaultValue: []string{"orange", "mango"},
			},
			want: []string{"orange", "mango"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			value := getEnvAsSlice(tc.args.key, tc.args.defaultValue, ".")
			require.Equal(t, value, tc.want)
		})
	}
}
