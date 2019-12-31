package appsvr

import (
	"github.com/jinzhu/gorm"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_initialMysqlTables(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := initialMysqlTables(tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("initialMysqlTables() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
