package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

func main() {
	var (
		expr     *cronexpr.Expression
		err      error
		now      time.Time
		nextTime time.Time
	)

	//every minutes
	//if expr, err = cronexpr.Parse("* * * * *"); err != nil {
	//	fmt.Println(err)
	//	return
	//}

	//every 5ms
	if expr = cronexpr.MustParse("*/5 * * * * * *"); err != nil {
		fmt.Println(err)
		return
	}

	//Time of next execution
	now = time.Now()
	nextTime = expr.Next(now)

	time.AfterFunc(nextTime.Sub(now), func() {
		fmt.Println("Be scheduled", nextTime)
	})

	time.Sleep(1 * time.Second)
}
