package a3interface

import (
	"reflect"
	"testing"
)

func TestRemoveEscapeQuotes(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "stringInArray",
			args: args{
				input: `"[""this is a string in an array""]`,
			},
			want: `["this is a string in an array"]`,
		}, {
			name: "nested stringInArray",
			args: args{
				input: `"[""[""""this is a string in an array""""]""]`,
			},
			want: `["[\"this is a string in an array\"]"]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveEscapeQuotes(tt.args.input); got != tt.want {
				t.Errorf("RemoveEscapeQuotes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseSQF(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "arrayOfArrays1",
			args: args{
				input: `"[""data1"", ""data2"", [""data3"", ""data4""]]"`,
			},
			want: []interface{}{
				"data1",
				"data2",
				[]interface{}{
					"data3",
					"data4",
				},
			},
			wantErr: false,
		}, {
			name: "arrayOfArrays2",
			args: args{
				input: `"[""data1"", ""data2"", [""data3"", ""data4"", [""data5"", ""data6""]]]"`,
			},
			want: []interface{}{
				"data1",
				"data2",
				[]interface{}{
					"data3",
					"data4",
					[]interface{}{
						"data5",
						"data6",
					},
				},
			},
			wantErr: false,
		}, {
			name: "arrayOfArrays3",
			args: args{
				input: `"[""data1"", ""data2"", [""data3"", 34, [""data5"",22]]]"`,
			},
			want: []interface{}{
				"data1",
				"data2",
				[]interface{}{
					"data3",
					float64(34),
					[]interface{}{
						"data5",
						float64(22),
					},
				},
			},
			wantErr: false,
		}, {
			name: "badArray1",
			args: args{
				input: `"[""data1"", ""data2"", [, ""data4"", [""data5"", ""data6""]]]`,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSQF(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSQF() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSQF() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseSQFHashMap(t *testing.T) {
	type args struct {
		input interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "basic hashmap",
			args: args{
				input: []interface{}{
					[]interface{}{
						"key1",
						"value1",
					},
					[]interface{}{
						"key2",
						"value2",
					},
				},
			},
			want: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			wantErr: false,
		}, {
			name: "nested hashmap",
			args: args{
				input: []interface{}{
					[]interface{}{
						"key1",
						"value1",
					},
					[]interface{}{
						"key2",
						[]interface{}{
							[]interface{}{
								"key3",
								"value3",
							},
						},
					},
				},
			},
			want: map[string]interface{}{
				"key1": "value1",
				"key2": map[string]interface{}{
					"key3": "value3",
				},
			},
			wantErr: false,
		}, {
			name: "bad hashmap",
			args: args{
				input: []interface{}{
					[]interface{}{
						"key1",
						"value1",
					},
					[]interface{}{
						"key2",
						"value2",
						"value3",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSQFHashMap(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSQFHashMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSQFHashMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
