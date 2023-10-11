package a3interface

import "testing"

func Test_escapeForSQF(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: `Basic quotation`,
			args: args{
				str: `He said "Hi there!"`,
			},
			want: `He said ""Hi there!""`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := escapeForSQF(tt.args.str); got != tt.want {
				t.Errorf("escapeForSQF() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToArmaHashMap(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: `Basic string`,
			args: args{
				data: `He said "Hi there!"`,
			},
			want: []interface{}{`"He said ""Hi there!"""`},
		},
		{
			name: `Basic int`,
			args: args{
				data: 1,
			},
			want: []interface{}{`1`},
		},
		{
			name: "map[string]interface{}",
			args: args{
				data: map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
			},
			want: []interface{}{
				// results are not ordered
				`[["key1", "value1"], ["key2", "value2"]]`,
				`[["key2", "value2"], ["key1", "value1"]]`,
			},
		},
		{
			name: "map[string]interface{} nested",
			args: args{
				data: map[string]interface{}{
					"key1": "value1",
					"key2": map[string]interface{}{
						"key3": "value3",
					},
				},
			},
			want: []interface{}{
				// results are not ordered
				`[["key1", "value1"], ["key2", [["key3", "value3"]]]]`,
				`[["key2", [["key3", "value3"]]], ["key1", "value1"]]`,
			},
		},
		{
			name: "[]map[string]interface{}",
			args: args{
				data: []map[string]interface{}{
					{
						"key1": "value1",
						"key2": []interface{}{
							24,
							"test",
							"equal",
						},
					},
					{
						"key3": "value3",
						"key4": 4,
					},
				},
			},
			want: []interface{}{
				// results are not ordered
				`[[["key1", "value1"], ["key2", [24, "test", "equal"]]], [["key3", "value3"], ["key4", 4]]]`,
				`[[["key3", "value3"], ["key4", 4]], [["key1", "value1"], ["key2", [24, "test", "equal"]]]]`,
				`[[["key1", "value1"], ["key2", [24, "test", "equal"]]], [["key4", 4], ["key3", "value3"]]]`,
				`[[["key4", 4], ["key3", "value3"]], [["key1", "value1"], ["key2", [24, "test", "equal"]]]]`,
				`[[["key2", [24, "test", "equal"]], ["key1", "value1"]], [["key3", "value3"], ["key4", 4]]]`,
				`[[["key3", "value3"], ["key4", 4]], [["key2", [24, "test", "equal"]], ["key1", "value1"]]]`,
				`[[["key2", [24, "test", "equal"]], ["key1", "value1"]], [["key4", 4], ["key3", "value3"]]]`,
				`[[["key4", 4], ["key3", "value3"]], [["key2", [24, "test", "equal"]], ["key1", "value1"]]]`,
			},
		},
		{
			name: "map[string]string",
			args: args{
				data: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			want: []interface{}{
				// results are not ordered
				`[["key1", "value1"], ["key2", "value2"]]`,
				`[["key2", "value2"], ["key1", "value1"]]`,
			},
		},
		{
			name: "[]interface{}",
			args: args{
				data: []interface{}{
					"test",
					24,
					[]interface{}{
						"test",
						24,
					},
				},
			},
			want: []interface{}{
				// results are not ordered
				`["test", 24, ["test", 24]]`,
				`["test", ["test", 24], 24]`,
				`[24, "test", ["test", 24]]`,
				`[24, ["test", 24], "test"]`,
				`[["test", 24], "test", 24]`,
				`[["test", 24], 24, "test"]`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToArmaHashMap(tt.args.data)
			found := false
			for _, want := range tt.want {
				if got == want {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("ToArmaHashMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
