package app

import (
	"reflect"
	"testing"
)

func Test_getChartInfo(t *testing.T) {
	// version string = the first semver-valid string after the last hypen in the chart string.
	type args struct {
		r *Release
	}
	tests := []struct {
		name string
		args args
		want *ChartInfo
	}{
		{
			name: "getChartInfo - local chart should return given release info",
			args: args{
				r: &Release{
					Name:      "release1",
					Namespace: "namespace",
					Version:   "1.0.0",
					Chart:     "./../../tests/chart-test",
					Enabled:   True,
				},
			},
			want: &ChartInfo{Name: "chart-test", Version: "1.0.0"},
		},
		{
			name: "getChartInfo - local chart semver should return latest matching release",
			args: args{
				r: &Release{
					Name:      "release1",
					Namespace: "namespace",
					Version:   "1.0.*",
					Chart:     "./../../tests/chart-test",
					Enabled:   True,
				},
			},
			want: &ChartInfo{Name: "chart-test", Version: "1.0.0"},
		},
		{
			name: "getChartInfo - unknown chart should error",
			args: args{
				r: &Release{
					Name:      "release1",
					Namespace: "namespace",
					Version:   "1.0.0",
					Chart:     "random-chart-name-1f8147",
					Enabled:   True,
				},
			},
			want: nil,
		},
		{
			name: "getChartInfo - wrong local version should error",
			args: args{
				r: &Release{
					Name:      "release1",
					Namespace: "namespace",
					Version:   "0.9.0",
					Chart:     "./../../tests/chart-test",
					Enabled:   True,
				},
			},
			want: nil,
		},
		{
			name: "getChartInfo - latest version should return chart info",
			args: args{
				r: &Release{
					Name:      "release1",
					Namespace: "namespace",
					Version:   "latest",
					Chart:     "./../../tests/chart-test",
					Enabled:   True,
				},
			},
			want: &ChartInfo{Name: "chart-test", Version: "1.0.0"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getChartInfo(tt.args.r.Chart, tt.args.r.Version)
			if err != nil && tt.want != nil {
				t.Errorf("getChartInfo() = Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getChartInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
