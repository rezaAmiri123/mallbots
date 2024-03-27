package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterCustomer(t *testing.T) {
	type args struct {
		id        string
		name      string
		smsNumber string
	}

	tests := map[string]struct {
		args args
		want *Customer
		wantErr  error
	}{}

	for name, tt := range tests{
		t.Run(name,func(t *testing.T) {
			got, err := RegisterCustomer(tt.args.id,tt.args.name,tt.args.smsNumber)
			if (err != nil)&&!errors.Is(err, tt.wantErr){
				t.Errorf("RemoveItem() error = %v, wantErr %v", err, tt.wantErr)
			} 

		})
	}
}
