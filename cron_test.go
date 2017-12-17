package gocron

import (
	"testing"
	"fmt"
	"time"
)

func TestNewCronTicker(t *testing.T) {
	_, err := NewCronTicker("0 * *")
	if nil != err {
		t.Fatal(err)
	}

	_, err = NewCronTicker("0/5 * *")
	if nil != err {
		t.Fatal(err)
	}

	_, err = NewCronTicker("* */10 *")
	if nil != err {
		t.Fatal(err)
	}
}

func TestCronTicker_Tick(t *testing.T) {
	tick, _ := NewCronTicker("*/1 * *")
	for {
		tick.Tick()
		now := time.Now()
		fmt.Printf("%d:%d:%d\n", now.Hour(), now.Minute(), now.Second())
	}
}