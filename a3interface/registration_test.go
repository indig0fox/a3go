package a3interface

import (
	"reflect"
	"testing"
)

func TestNewRegistration(t *testing.T) {
	type args struct {
		command string
	}
	tests := []struct {
		name string
		args args
		want *RVExtensionRegistration
	}{
		{
			name: "test",
			args: args{
				command: "test",
			},
			want: &RVExtensionRegistration{
				Command:         "test",
				DefaultResponse: `["Command test called"]`,
			},
		}, {
			name: "test2",
			args: args{
				command: "test2",
			},
			want: &RVExtensionRegistration{
				Command:         "test2",
				DefaultResponse: `["Command test2 called"]`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRegistration(tt.args.command); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRegistration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRVExtensionRegistration_SetDefaultResponse(t *testing.T) {
	type args struct {
		response string
	}
	tests := []struct {
		name string
		r    *RVExtensionRegistration
		args args
		want *RVExtensionRegistration
	}{
		{
			name: "test",
			r: &RVExtensionRegistration{
				Command:         "test",
				DefaultResponse: `["Command test called"]`,
			},
			args: args{
				response: `["test"]`,
			},
			want: &RVExtensionRegistration{
				Command:         "test",
				DefaultResponse: `["test"]`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.SetDefaultResponse(tt.args.response); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RVExtensionRegistration.SetDefaultResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRVExtensionRegistration_SetRunInBackground(t *testing.T) {
	type args struct {
		runInBackground bool
	}
	tests := []struct {
		name string
		r    *RVExtensionRegistration
		args args
		want *RVExtensionRegistration
	}{
		{
			name: "test",
			r: &RVExtensionRegistration{
				Command:         "test",
				DefaultResponse: `["Command test called"]`,
			},
			args: args{
				runInBackground: true,
			},
			want: &RVExtensionRegistration{
				Command:         "test",
				DefaultResponse: `["Command test called"]`,
				RunInBackground: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.SetRunInBackground(tt.args.runInBackground); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RVExtensionRegistration.SetRunInBackground() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRVExtensionRegistration_SetFunction(t *testing.T) {
	type args struct {
		fnc func(ctx ArmaExtensionContext, data string) (string, error)
	}
	tests := []struct {
		name string
		r    *RVExtensionRegistration
		args args
		want *RVExtensionRegistration
	}{
		{
			name: "add a function",
			r: &RVExtensionRegistration{
				Command:         "test",
				DefaultResponse: `["Command test called"]`,
			},
			args: args{
				fnc: func(ctx ArmaExtensionContext, data string) (string, error) {
					return "", nil
				},
			},
			want: &RVExtensionRegistration{
				Command:         "test",
				DefaultResponse: `["Command test called"]`,
				Function: func(ctx ArmaExtensionContext, data string) (string, error) {
					return "", nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.r.SetFunction(tt.args.fnc)
			if got.Function == nil {
				t.Errorf("RVExtensionRegistration.SetFunction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRVExtensionRegistration_SetArgsFunction(t *testing.T) {
	type args struct {
		fnc func(ctx ArmaExtensionContext, command string, args []string) (string, error)
	}
	tests := []struct {
		name string
		r    *RVExtensionRegistration
		args args
		want *RVExtensionRegistration
	}{
		{
			name: "add a function",
			r: &RVExtensionRegistration{
				Command:         "test",
				DefaultResponse: `["Command test called"]`,
			},
			args: args{
				fnc: func(ctx ArmaExtensionContext, command string, args []string) (string, error) {
					return "", nil
				},
			},
			want: &RVExtensionRegistration{
				Command:         "test",
				DefaultResponse: `["Command test called"]`,
				ArgsFunction: func(ctx ArmaExtensionContext, command string, args []string) (string, error) {
					return "", nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.SetArgsFunction(tt.args.fnc); got.ArgsFunction == nil {
				t.Errorf("RVExtensionRegistration.SetArgsFunction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRVExtensionRegistration_Register(t *testing.T) {
	tests := []struct {
		name    string
		r       *RVExtensionRegistration
		wantErr bool
	}{
		{
			name: "register new",
			r: &RVExtensionRegistration{
				Command:         "test2",
				DefaultResponse: `["Command test2 called"]`,
			},
			wantErr: false,
		},
		{
			name: "register duplicate",
			r: &RVExtensionRegistration{
				Command:         "test",
				DefaultResponse: `["Command test called"]`,
			},
			wantErr: true,
		},
	}

	// create a new registration as 'test' for checking duplicates
	err := NewRegistration("test").Register()
	if err != nil {
		t.Errorf("RVExtensionRegistration.Register() error = %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.r.Register()
			if (err != nil) != tt.wantErr {
				t.Errorf("RVExtensionRegistration.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
