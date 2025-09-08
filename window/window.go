package window

import (
	"errors"
	"iter"
	"time"

	"github.com/copito/coptime/interval"
)

type Windower struct {
	opt WindowOption
}

func New(option WindowOption) Windower {
	return Windower{
		opt: option,
	}
}

func (w *Windower) Iterate(direction interval.Direction, maxAttempt *int32) (iter.Seq[WindowResult], error) {
	intervalOption := interval.IntervalOption{
		AnchorDate:    w.opt.AnchorDate,
		StartDate:     w.opt.StartDate,
		EndDate:       w.opt.EndDate,
		FrequencyUnit: w.opt.FrequencyUnit,
		IntervalValue: w.opt.IntervalValue,
	}

	iv := interval.New(intervalOption)
	iterator, err := iv.Iterate(direction, maxAttempt)
	if err != nil {
		return nil, errors.Join(errors.New("failed to create iterator for window"), err)
	}

	return func(yield func(WindowResult) bool) {
		next, _ := iter.Pull(iterator)
		// defer stop()

		previousTime, ok := next()
		if !ok {
			return
		}

		for {
			nextTime, ok := next()
			if !ok {
				return
			}

			w := WindowResult{
				StartWindow: previousTime,
				EndWindow:   nextTime,
				SubWindow:   []SubWindowResult{},
			}

			ok = yield(w)
			if !ok {
				return
			}

			previousTime = nextTime
		}
	}, nil
}

func (w *Windower) All(direction interval.Direction, maxAttempt *int32) ([]WindowResult, error) {
	// panic("unimplemented")
	return nil, nil
}

func (w *Windower) Between(direction interval.Direction, startTime time.Time, endTime time.Time, maxAttempt *int32) ([]WindowResult, error) {
	// panic("unimplemented")
	return nil, nil
}
