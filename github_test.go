package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestGetLastWeekContributions(t *testing.T) {
	cases := map[string]struct {
		in   string
		want []string
	}{
		"Sanday": {
			in:   "2024-02-18",
			want: []string{"2024-02-11", "2024-02-12", "2024-02-13", "2024-02-14", "2024-02-15", "2024-02-16", "2024-02-17"},
		},
		"Monday": {
			in:   "2024-02-19",
			want: []string{"2024-02-11", "2024-02-12", "2024-02-13", "2024-02-14", "2024-02-15", "2024-02-16", "2024-02-17"},
		},
		"Saturday": {
			in:   "2024-02-24",
			want: []string{"2024-02-18", "2024-02-19", "2024-02-20", "2024-02-21", "2024-02-22", "2024-02-23", "2024-02-24"},
		},
	}

	client, err := createClient()
	if err != nil {
		t.Fatal(err)
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := getLastWeekContributions(client, toDateTime(t, tt.in), "oka4shi")
			if err != nil {
				t.Fatal(err)
			}

			var want []time.Time
			for _, v := range tt.want {
				want = append(want, toDate(t, v))
			}

			var days []time.Time
			for _, v := range got {
				days = append(days, v.Date)
			}

			if len(days) != 7 {
				t.Fatal("The data doesn't have 7 days")
			}

			if !reflect.DeepEqual(want, days) {
				fmt.Println("want: ", want)
				fmt.Println("data: ", days)
				t.Fatal("Unexpected data")
			}
		})
	}

}

func toDate(t *testing.T, date string) time.Time {
	t.Helper()
	d, err := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s 00:00:00", date))
	if err != nil {
		t.Fatalf("toDate: %v", err)
	}

	return d
}

func toDateTime(t *testing.T, date string) time.Time {
	t.Helper()
	d, err := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s 12:34:56", date))
	if err != nil {
		t.Fatalf("toDate: %v", err)
	}

	return d

}
