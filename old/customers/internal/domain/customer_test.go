package domain

// import (
// 	"errors"
// 	"testing"

// 	"github.com/rezaAmiri123/edatV2/ddd"
// 	"github.com/stretchr/testify/assert"
// )

// func TestRegisterCustomer(t *testing.T) {
// 	type args struct {
// 		id        string
// 		name      string
// 		smsNumber string
// 	}

// 	customer := Customer{
// 		Aggregate: ddd.NewAggregate("customer-id", ""),
// 		Name:      "customer-name",
// 		SmsNumber: "customer-sms-number",
// 		Enabled:   true,
// 	}
// 	_ = customer
// 	tests := map[string]struct {
// 		args    args
// 		want    *Customer
// 		wantErr error
// 	}{
// 		"OK": {
// 			args: args{
// 				id:        customer.ID(),
// 				name:      customer.Name,
// 				smsNumber: customer.SmsNumber,
// 			},
// 			want:    &customer,
// 			wantErr: nil,
// 		},
// 		"no id error": {
// 			args: args{
// 				id:        "",
// 				name:      customer.Name,
// 				smsNumber: customer.SmsNumber,
// 			},
// 			want:    nil,
// 			wantErr: ErrCustomerIDCannotBeBlank,
// 		},
// 		"no name error": {
// 			args: args{
// 				id:        customer.ID(),
// 				name:      "",
// 				smsNumber: customer.SmsNumber,
// 			},
// 			want:    nil,
// 			wantErr: ErrNameCannotBeBlank,
// 		},
// 		"no sms number error": {
// 			args: args{
// 				id:        customer.ID(),
// 				name:      customer.Name,
// 				smsNumber: "",
// 			},
// 			want:    nil,
// 			wantErr: ErrSmsNumberCannotBeBlank,
// 		},

// 	}

// 	for name, tt := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			got, err := RegisterCustomer(tt.args.id, tt.args.name, tt.args.smsNumber)

// 			if (err != nil) && !errors.Is(err, tt.wantErr) {
// 				t.Errorf("RegisterCustomer() error = %v, wantErr %v", err, tt.wantErr)
// 			}

// 			if got == nil {
// 				if tt.want != nil {
// 					t.Errorf("RegisterCustomer() got = %v, want %v", got, tt.want)
// 				} 
// 				return
// 			}

// 			assert.Equal(t, tt.want.ID(), tt.args.id)
// 			assert.Equal(t, tt.want.Name, tt.args.name)
// 			assert.Equal(t, tt.want.SmsNumber, tt.args.smsNumber)
// 			assert.True(t, tt.want.Enabled)

// 			events := got.Events()
// 			assert.Greater(t, len(events), 0)
// 		})
// 	}
// }
