package main

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func mustRFC3339(t *testing.T, value string) time.Time {
	tm, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatal(err)
	}
	return tm
}

func assertPeriods(t *testing.T, periods []period, expected string) {
	var buf bytes.Buffer
	for _, period := range periods {
		fmt.Fprintf(&buf, "%s [%s, %s)\n",
			period.label,
			period.from.Format(time.RFC3339),
			period.to.Format(time.RFC3339))
	}

	if got := buf.String(); got != expected {
		t.Errorf("unexpected periods:\ngot:\n%swant:\n%s", got, expected)
	}
}

func TestCumulativePeriods(t *testing.T) {
	// Include DST transition in the period.
	start := mustRFC3339(t, "2021-03-27T09:00:00+02:00")
	stop := mustRFC3339(t, "2021-03-29T20:15:00+03:00")
	loc, err := time.LoadLocation("Europe/Tallinn")
	if err != nil {
		t.Fatal(err)
	}

	periods := cumulativePeriods(start, stop, loc)

	assertPeriods(t, periods, ""+
		"3 [2021-03-27T09:00:00+02:00, 2021-03-28T00:00:00+02:00)\n"+
		"2 [2021-03-27T09:00:00+02:00, 2021-03-29T00:00:00+03:00)\n"+
		"1 [2021-03-27T09:00:00+02:00, 2021-03-29T20:15:00+03:00)\n")
}

func TestDayBeforeInternetVotingCumulativePeriodsNotAfter(t *testing.T) {
	// Include DST transition in the period.
	start := mustRFC3339(t, "2024-06-03T09:00:00+02:00")
	notAfter := mustRFC3339(t, "2024-06-08T20:15:00+02:00")
	stop := mustRFC3339(t, "2024-06-02T13:00:00+03:00")
	loc, err := time.LoadLocation("Europe/Tallinn")
	if err != nil {
		t.Fatal(err)
	}

	periods := cumulativePeriodsNotAfter(start, stop, notAfter, loc)

	assertPeriods(t, periods, ""+
		"6 [2024-06-03T09:00:00+02:00, 2024-06-04T00:00:00+03:00)\n")
}

func TestFirstDayOfInternetVotingCumulativePeriodsNotAfter(t *testing.T) {
	// Include DST transition in the period.
	start := mustRFC3339(t, "2024-06-03T09:00:00+02:00")
	notAfter := mustRFC3339(t, "2024-06-08T20:15:00+02:00")
	stop := mustRFC3339(t, "2024-06-03T23:59:59+03:00")

	loc, err := time.LoadLocation("Europe/Tallinn")
	if err != nil {
		t.Fatal(err)
	}

	periods := cumulativePeriodsNotAfter(start, stop, notAfter, loc)

	assertPeriods(t, periods, ""+
		"6 [2024-06-03T09:00:00+02:00, 2024-06-04T00:00:00+03:00)\n")
}

func TestSecondDayOfInternetVotingCumulativePeriodsNotAfter(t *testing.T) {
	// Include DST transition in the period.
	start := mustRFC3339(t, "2024-06-03T09:00:00+02:00")
	notAfter := mustRFC3339(t, "2024-06-08T20:15:00+02:00")
	stop := mustRFC3339(t, "2024-06-04T23:59:59+03:00")

	loc, err := time.LoadLocation("Europe/Tallinn")
	if err != nil {
		t.Fatal(err)
	}

	periods := cumulativePeriodsNotAfter(start, stop, notAfter, loc)

	assertPeriods(t, periods, ""+
		"6 [2024-06-03T09:00:00+02:00, 2024-06-04T00:00:00+03:00)\n"+
		"5 [2024-06-03T09:00:00+02:00, 2024-06-05T00:00:00+03:00)\n")
}

func TestThirdDayOfInternetVotingCumulativePeriodsNotAfter(t *testing.T) {
	// Include DST transition in the period.
	start := mustRFC3339(t, "2024-06-03T09:00:00+02:00")
	notAfter := mustRFC3339(t, "2024-06-08T20:15:00+02:00")
	stop := mustRFC3339(t, "2024-06-05T23:59:59+03:00")

	loc, err := time.LoadLocation("Europe/Tallinn")
	if err != nil {
		t.Fatal(err)
	}

	periods := cumulativePeriodsNotAfter(start, stop, notAfter, loc)

	assertPeriods(t, periods, ""+
		"6 [2024-06-03T09:00:00+02:00, 2024-06-04T00:00:00+03:00)\n"+
		"5 [2024-06-03T09:00:00+02:00, 2024-06-05T00:00:00+03:00)\n"+
		"4 [2024-06-03T09:00:00+02:00, 2024-06-06T00:00:00+03:00)\n")
}

func TestFourthDayOfInternetVotingCumulativePeriodsNotAfter(t *testing.T) {
	// Include DST transition in the period.
	start := mustRFC3339(t, "2024-06-03T09:00:00+02:00")
	notAfter := mustRFC3339(t, "2024-06-08T20:15:00+02:00")
	stop := mustRFC3339(t, "2024-06-06T00:00:00+03:00")

	loc, err := time.LoadLocation("Europe/Tallinn")
	if err != nil {
		t.Fatal(err)
	}

	periods := cumulativePeriodsNotAfter(start, stop, notAfter, loc)

	assertPeriods(t, periods, ""+
		"6 [2024-06-03T09:00:00+02:00, 2024-06-04T00:00:00+03:00)\n"+
		"5 [2024-06-03T09:00:00+02:00, 2024-06-05T00:00:00+03:00)\n"+
		"4 [2024-06-03T09:00:00+02:00, 2024-06-06T00:00:00+03:00)\n"+
		"3 [2024-06-03T09:00:00+02:00, 2024-06-07T00:00:00+03:00)\n")
}

func TestFifthDayOfInternetVotingCumulativePeriodsNotAfter(t *testing.T) {
	// Include DST transition in the period.
	start := mustRFC3339(t, "2024-06-03T09:00:00+02:00")
	notAfter := mustRFC3339(t, "2024-06-08T20:15:00+02:00")
	stop := mustRFC3339(t, "2024-06-07T00:00:00+03:00")

	loc, err := time.LoadLocation("Europe/Tallinn")
	if err != nil {
		t.Fatal(err)
	}

	periods := cumulativePeriodsNotAfter(start, stop, notAfter, loc)

	assertPeriods(t, periods, ""+
		"6 [2024-06-03T09:00:00+02:00, 2024-06-04T00:00:00+03:00)\n"+
		"5 [2024-06-03T09:00:00+02:00, 2024-06-05T00:00:00+03:00)\n"+
		"4 [2024-06-03T09:00:00+02:00, 2024-06-06T00:00:00+03:00)\n"+
		"3 [2024-06-03T09:00:00+02:00, 2024-06-07T00:00:00+03:00)\n"+
		"2 [2024-06-03T09:00:00+02:00, 2024-06-08T00:00:00+03:00)\n")
}

func TestSixthAkaLastDayOfInternetVotingCumulativePeriodsNotAfter(t *testing.T) {
	// Include DST transition in the period.
	start := mustRFC3339(t, "2024-06-03T09:00:00+02:00")
	notAfter := mustRFC3339(t, "2024-06-08T20:15:00+02:00")
	stop := mustRFC3339(t, "2024-06-08T00:00:00+03:00")

	loc, err := time.LoadLocation("Europe/Tallinn")
	if err != nil {
		t.Fatal(err)
	}

	periods := cumulativePeriodsNotAfter(start, stop, notAfter, loc)

	assertPeriods(t, periods, ""+
		"6 [2024-06-03T09:00:00+02:00, 2024-06-04T00:00:00+03:00)\n"+
		"5 [2024-06-03T09:00:00+02:00, 2024-06-05T00:00:00+03:00)\n"+
		"4 [2024-06-03T09:00:00+02:00, 2024-06-06T00:00:00+03:00)\n"+
		"3 [2024-06-03T09:00:00+02:00, 2024-06-07T00:00:00+03:00)\n"+
		"2 [2024-06-03T09:00:00+02:00, 2024-06-08T00:00:00+03:00)\n"+
		"1 [2024-06-03T09:00:00+02:00, 2024-06-08T20:15:00+02:00)\n")
}

func TestSixthAkaLastDayOfInternetVotingLastSecondCumulativePeriodsNotAfter(t *testing.T) {
	// Include DST transition in the period.
	start := mustRFC3339(t, "2024-06-03T09:00:00+02:00")
	notAfter := mustRFC3339(t, "2024-06-08T20:15:00+02:00")
	stop := mustRFC3339(t, "2024-06-08T20:14:59+03:00")

	loc, err := time.LoadLocation("Europe/Tallinn")
	if err != nil {
		t.Fatal(err)
	}

	periods := cumulativePeriodsNotAfter(start, stop, notAfter, loc)

	assertPeriods(t, periods, ""+
		"6 [2024-06-03T09:00:00+02:00, 2024-06-04T00:00:00+03:00)\n"+
		"5 [2024-06-03T09:00:00+02:00, 2024-06-05T00:00:00+03:00)\n"+
		"4 [2024-06-03T09:00:00+02:00, 2024-06-06T00:00:00+03:00)\n"+
		"3 [2024-06-03T09:00:00+02:00, 2024-06-07T00:00:00+03:00)\n"+
		"2 [2024-06-03T09:00:00+02:00, 2024-06-08T00:00:00+03:00)\n"+
		"1 [2024-06-03T09:00:00+02:00, 2024-06-08T20:15:00+02:00)\n")
}

func TestSixthAkaLastDayOfInternetVotingSecondAfterEndCumulativePeriodsNotAfter(t *testing.T) {
	// Include DST transition in the period.
	start := mustRFC3339(t, "2024-06-03T09:00:00+02:00")
	notAfter := mustRFC3339(t, "2024-06-08T20:15:00+02:00")
	stop := mustRFC3339(t, "2024-06-08T20:15:01+03:00")

	loc, err := time.LoadLocation("Europe/Tallinn")
	if err != nil {
		t.Fatal(err)
	}

	periods := cumulativePeriodsNotAfter(start, stop, notAfter, loc)

	assertPeriods(t, periods, ""+
		"6 [2024-06-03T09:00:00+02:00, 2024-06-04T00:00:00+03:00)\n"+
		"5 [2024-06-03T09:00:00+02:00, 2024-06-05T00:00:00+03:00)\n"+
		"4 [2024-06-03T09:00:00+02:00, 2024-06-06T00:00:00+03:00)\n"+
		"3 [2024-06-03T09:00:00+02:00, 2024-06-07T00:00:00+03:00)\n"+
		"2 [2024-06-03T09:00:00+02:00, 2024-06-08T00:00:00+03:00)\n"+
		"1 [2024-06-03T09:00:00+02:00, 2024-06-08T20:15:00+02:00)\n")
}
