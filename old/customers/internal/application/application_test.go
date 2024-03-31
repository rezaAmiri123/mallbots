package application

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/rezaAmiri123/edatV2/ddd"
// 	"github.com/rezaAmiri123/mallbots/customers/internal/domain"
// 	"github.com/stretchr/testify/mock"
// )

// func TestApplication_RegisterCustomer(t *testing.T) {
// 	type mocks struct {
// 		customers       *domain.MockCustomerRepository
// 		domainPublisher *ddd.MockEventPublisher[ddd.AggregateEvent]
// 	}

// 	type args struct {
// 		ctx      context.Context
// 		register RegisterCustomer
// 	}
// 	register := RegisterCustomer{
// 		ID:        "customer-id",
// 		Name:      "customer-name",
// 		SmsNumber: "sms number",
// 	}
// 	tests := map[string]struct {
// 		args    args
// 		on      func(m mocks)
// 		wantErr bool
// 	}{
// 		"Success": {
// 			args: args{
// 				ctx:      context.Background(),
// 				register: register,
// 			},
// 			on: func(m mocks) {
// 				m.customers.On("Save", mock.Anything, mock.AnythingOfType("*domain.Customer")).Return(nil)
// 				m.domainPublisher.On("Publish", mock.Anything, mock.Anything).Return(nil)
// 			},
// 			wantErr: false,
// 		},
// 		"publishFailed": {
// 			args: args{
// 				ctx:      context.Background(),
// 				register: register,
// 			},
// 			on: func(m mocks) {
// 				m.customers.On("Save", mock.Anything, mock.AnythingOfType("*domain.Customer")).Return(nil)
// 				m.domainPublisher.On("Publish", mock.Anything, mock.Anything).Return(errors.New("publish failed"))
// 			},
// 			wantErr: true,
// 		},
// 		"SaveFailed": {
// 			args: args{
// 				ctx:      context.Background(),
// 				register: register,
// 			},
// 			on: func(m mocks) {
// 				m.customers.On("Save", mock.Anything, mock.AnythingOfType("*domain.Customer")).Return(errors.New("publish failed"))
// 			},
// 			wantErr: true,
// 		},
// 		"RegisterFailed": {
// 			args: args{
// 				ctx:      context.Background(),
// 				register: RegisterCustomer{}, // no argumant
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			m := mocks{
// 				customers:       domain.NewMockCustomerRepository(t),
// 				domainPublisher: ddd.NewMockEventPublisher[ddd.AggregateEvent](t),
// 			}
// 			a := Application{
// 				customers:       m.customers,
// 				domainPublisher: m.domainPublisher,
// 			}

// 			if tc.on != nil {
// 				tc.on(m)
// 			}

// 			err := a.RegisterCustomer(tc.args.ctx, tc.args.register)
// 			if (err != nil) != tc.wantErr {
// 				t.Errorf("RegisterCustomer() error = %v, wantErr %v", err, tc.wantErr)
// 			}
// 		})
// 	}
// }
