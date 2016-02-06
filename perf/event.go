package perf

import (
	"errors"
	"fmt"
	"time"
)

var ZeroTime = time.Unix(0, 0)

type Event struct {
	name             string
	first_start_time time.Time
	start_time       time.Time
	end_time         time.Time
	event_duration   time.Duration
	total_duration   time.Duration
	event_count      int64
}

func MakeEvent(name string) *Event {
	pe := &Event{
		name:             name,
		first_start_time: ZeroTime,
		start_time:       ZeroTime,
		end_time:         ZeroTime,
		event_duration:   time.Duration(0),
		total_duration:   time.Duration(0),
		event_count:      0,
	}

	return pe
}

func (pe *Event) Start() error {
	if !pe.start_time.Equal(ZeroTime) {
		return errors.New("Tried to start event that already started!")
	}
	pe.start_time = time.Now()
	if pe.first_start_time.Equal(ZeroTime) {
		pe.first_start_time = pe.start_time
	}
	return nil
}

func (pe *Event) Stop() error {
	if pe.start_time.Equal(ZeroTime) {
		return errors.New("Tried to stop event that didn't start!")
	}
	pe.end_time = time.Now()
	pe.event_duration = pe.end_time.Sub(pe.start_time)
	pe.total_duration += pe.event_duration
	pe.event_count++
	pe.start_time = time.Unix(0, 0)
	return nil
}

func (pe *Event) Print(debug bool) {
	var avg_resp float64 = -1.0
	var bench_time time.Duration = time.Duration(0)
	var trans_rate float64 = -1.0

	if pe.event_count != 0 {
		avg_resp = float64(pe.total_duration) / float64(pe.event_count) /
			float64(time.Millisecond)
	}
	if !pe.first_start_time.Equal(ZeroTime) && !pe.end_time.Equal(ZeroTime) {
		bench_time = pe.end_time.Sub(pe.first_start_time)
	}
	if pe.event_count != 0 && bench_time != time.Duration(0) {
		trans_rate = float64(pe.event_count) / bench_time.Seconds()
	}

	fmt.Printf("%10s,%10s,%12s,%15s,%12s,%15s\n",
		"name", "e_cnt_t", "e_time_t(s)", "e_avg_rt(ms/e)",
		"b_time_t(s)", "eff_rate(e/s)")
	fmt.Printf("%10s,%10d,%12.6f,%15.6f,%12.6f,%15.6f\n",
		pe.name, pe.event_count, pe.total_duration.Seconds(),
		avg_resp, bench_time.Seconds(), trans_rate)

	if debug {
		fmt.Printf("   s: %v, e: %v, d: %v\n", pe.start_time, pe.end_time, pe.event_duration)
	}
}

func (pe *Event) Clear() {
	pe.start_time = time.Unix(0, 0)
	pe.end_time = time.Unix(0, 0)
	pe.event_duration = time.Duration(0)
	pe.total_duration = time.Duration(0)
	pe.event_count = 0
}
