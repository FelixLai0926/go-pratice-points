package config

import (
	"reflect"
	"testing"
)

func Test_setFieldValue(t *testing.T) {
	type TestStruct struct {
		IntField    int
		FloatField  float64
		BoolField   bool
		StringField string
	}
	type args struct {
		field reflect.Value
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test setFieldValue (string)",
			args: args{
				field: reflect.ValueOf(&TestStruct{}).Elem().FieldByName("StringField"),
				value: "test",
			},
			wantErr: false,
		},
		{
			name: "Test setFieldValue (int)",
			args: args{
				field: reflect.ValueOf(&TestStruct{}).Elem().FieldByName("IntField"),
				value: "10",
			},
			wantErr: false,
		},
		{
			name: "Test setFieldValue (float)",
			args: args{
				field: reflect.ValueOf(&TestStruct{}).Elem().FieldByName("FloatField"),
				value: "10.5",
			},
			wantErr: false,
		},
		{
			name: "Test setFieldValue (bool)",
			args: args{
				field: reflect.ValueOf(&TestStruct{}).Elem().FieldByName("BoolField"),
				value: "true",
			},
			wantErr: false,
		},
		{
			name: "Test setFieldValue (invalid int)",
			args: args{
				field: reflect.ValueOf(&TestStruct{}).Elem().FieldByName("IntField"),
				value: "test",
			},
			wantErr: true,
		},
		{
			name: "Test setFieldValue (invalid float)",
			args: args{
				field: reflect.ValueOf(&TestStruct{}).Elem().FieldByName("FloatField"),
				value: "test",
			},
			wantErr: true,
		},
		{
			name: "Test setFieldValue (invalid bool)",
			args: args{
				field: reflect.ValueOf(&TestStruct{}).Elem().FieldByName("BoolField"),
				value: "test",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setFieldValue(tt.args.field, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("setFieldValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
