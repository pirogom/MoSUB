package main

import "time"

//
// random util from tbaMud
//
type dice struct {
	randM    int64
	randQ    int64
	randA    int64
	randR    int64
	randSeed int64
}

func newDice() *dice {
	nd := dice{}
	nd.init()
	return &nd
}

func (d *dice) init() {
	d.randM = 2147483647
	d.randQ = 127773
	d.randA = 16807
	d.randR = 2836
	d.randSeed = time.Now().Unix()
}

func (d *dice) random() int64 {
	hi := d.randSeed / d.randQ
	lo := d.randSeed % d.randQ

	test := d.randA*lo - d.randR*hi

	if test > 0 {
		d.randSeed = test
	} else {
		d.randSeed = test + d.randM
	}

	return d.randSeed
}

func (d *dice) randNumber(from int64, to int64) int64 {
	if from > to {
		tmp := from
		from = to
		to = tmp
	}
	return ((d.random() % (to - from + 1)) + from)
}

func (d *dice) dice(num int64, size int64) int64 {
	var sum int64
	sum = 0

	if size <= 0 || num <= 0 {
		return 0
	}

	for {
		if num <= 0 {
			break
		}
		sum += d.randNumber(1, size)
		num--
	}

	return sum
}
