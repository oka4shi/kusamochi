package main

import (
	"context"
	"time"

	"github.com/Khan/genqlient/graphql"
)

func getDateRange(now time.Time, origin int, duration int) dateRange {
	return dateRange{
		To:   now.AddDate(0, 0, -origin),
		From: now.AddDate(0, 0, -origin-(duration-1)),
	}
}

type weeklyContributions getUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeek

func getLastWeekContributions(c graphql.Client, now time.Time, user string) (weeklyContributions, error) {
	d := (int(now.Weekday()) + 1) % 7
	r := getDateRange(now, d, 7)

	var resp *getUserContributionsResponse
	resp, err := getUserContributions(context.Background(), c, user, r.To, r.From)
	if err != nil {
		return weeklyContributions{}, err
	}

	return weeklyContributions(resp.User.ContributionsCollection.ContributionCalendar.Weeks[0]), nil
}
