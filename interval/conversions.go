package interval

import (
	"errors"
	"fmt"

	"github.com/teambition/rrule-go"
)

func convertRRULEtoIntervaler(ruleString string) (*IntervalOption, error) {
	r, err := rrule.StrToRRule(ruleString)
	if err != nil {
		return nil, err
	}

	// Implement here
	fmt.Println(r.Options.Dtstart)
	fmt.Println(r.Options.Freq)
	fmt.Println(r.Options.Interval)
	fmt.Println(r.Options.Byhour)
	fmt.Println(r.Options.Until)
	fmt.Println(r.Options.Bymonth)
	fmt.Println(r.Options.Bymonthday)
	fmt.Println(r.Options.Byhour)
	fmt.Println(r.Options.Byminute)

	return nil, errors.New("not implemented")
}

func convertCrontoIntervaler(cronString string) (*IntervalOption, error) {
	return nil, errors.New("not implemented")
}
