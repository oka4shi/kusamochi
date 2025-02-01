package github

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestGetDataRange(t *testing.T) {
	type getDateTestIn struct {
		now      string
		origin   int
		duration int
	}

	cases := map[string]struct {
		in   getDateTestIn
		want DateRange
	}{
		"origin: -, duration: +": {
			in:   getDateTestIn{now: "2025-02-02", origin: -1, duration: 7},
			want: DateRange{From: toDate(t, "2025-02-01"), To: toDate(t, "2025-02-07")},
		},
		"origin: -, duration: -": {
			in:   getDateTestIn{now: "2025-02-02", origin: -1, duration: -7},
			want: DateRange{From: toDate(t, "2025-01-26"), To: toDate(t, "2025-02-01")},
		},
		"origin: +, duration: +": {
			in:   getDateTestIn{now: "2025-02-02", origin: 1, duration: 7},
			want: DateRange{From: toDate(t, "2025-02-03"), To: toDate(t, "2025-02-09")},
		},
		"origin: +, duration: -": {
			in:   getDateTestIn{now: "2025-02-02", origin: 1, duration: -7},
			want: DateRange{From: toDate(t, "2025-01-28"), To: toDate(t, "2025-02-03")},
		},
		"origin: 0, duration: +": {
			in:   getDateTestIn{now: "2025-02-02", origin: 0, duration: 7},
			want: DateRange{From: toDate(t, "2025-02-02"), To: toDate(t, "2025-02-08")},
		},
		"origin: 0, duration: -": {
			in:   getDateTestIn{now: "2025-02-02", origin: 0, duration: -7},
			want: DateRange{From: toDate(t, "2025-01-27"), To: toDate(t, "2025-02-02")},
		},
		"origin: 0, duration: 0": {
			in:   getDateTestIn{now: "2025-02-02", origin: 0, duration: -0},
			want: DateRange{From: toDate(t, "2025-02-02"), To: toDate(t, "2025-02-02")},
		},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			got := getDateRange(toDate(t, tt.in.now), tt.in.origin, tt.in.duration)

			if !reflect.DeepEqual(tt.want, got) {
				fmt.Println("want: ", tt.want)
				fmt.Println("data: ", got)
				t.Fatal("Unexpected data")
			}
		})
	}
}

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

	ghToken := os.Getenv("KUSAMOCHI_GITHUB_TOKEN")
	if ghToken == "" {
		err := fmt.Errorf("must set KUSAMOCHI_GITHUB_TOKEN")
		t.Fatal(err)
	}
	client := CreateGraphQLClient(ghToken)

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := GetLastWeekContributions(client, toDateTime(t, tt.in), "oka4shi")
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
