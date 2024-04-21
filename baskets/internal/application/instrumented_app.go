package application

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

type instrumentedApp struct {
	App
	basketsStarted prometheus.Counter
}

var _ App = (*instrumentedApp)(nil)

func NewInstrumentedApp(app App, basketsStarted prometheus.Counter) App {
	return instrumentedApp{
		App:            app,
		basketsStarted: basketsStarted,
	}
}

func (a instrumentedApp) StartBasket(ctx context.Context, start StartBasket) error {
	err := a.App.StartBasket(ctx, start)
	if err != nil {
		return err
	}
	a.basketsStarted.Inc()
	return nil
}
