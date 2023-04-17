package routine

import (
	"github.com/rfw141/i/context"
	"github.com/rfw141/i/log"
)

func Go(ctx context.IContext, fn func(context.IContext) error) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("Go panic: %v", err)
			}
		}()
		err := fn(ctx)
		if err != nil {
			log.Errorf("Go fail: %v", err)
		}
	}()

}
