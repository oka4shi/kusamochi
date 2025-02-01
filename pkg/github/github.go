package github

import (
	"context"
	"time"

	"github.com/Khan/genqlient/graphql"
)

type DateRange struct {
	To   time.Time
	From time.Time
}

func getDateRange(now time.Time, origin int, duration int) DateRange {
	return DateRange{
		To:   now.AddDate(0, 0, -origin),
		From: now.AddDate(0, 0, -origin-(duration-1)),
	}
}

type WeeklyContributions []getUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeekContributionDaysContributionCalendarDay

func GetLastWeekContributions(c *graphql.Client, now time.Time, user string) (WeeklyContributions, error) {
	d := (int(now.Weekday()) + 1) % 7
	r := getDateRange(now, d, 7)

	var resp *getUserContributionsResponse
	resp, err := getUserContributions(context.Background(), *c, user, r.To, r.From)
	if err != nil {
		return WeeklyContributions{}, err
	}

	return WeeklyContributions(resp.User.ContributionsCollection.ContributionCalendar.Weeks[0].ContributionDays), nil
}
