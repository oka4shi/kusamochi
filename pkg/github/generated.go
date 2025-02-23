// Code generated by github.com/Khan/genqlient, DO NOT EDIT.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/oka4shi/kusamochi/pkg/bindings"
)

// __getUserContributionsInput is used internally by genqlient
type __getUserContributionsInput struct {
	UserName string    `json:"userName"`
	To       time.Time `json:"-"`
	From     time.Time `json:"-"`
}

// GetUserName returns __getUserContributionsInput.UserName, and is useful for accessing the field via an interface.
func (v *__getUserContributionsInput) GetUserName() string { return v.UserName }

// GetTo returns __getUserContributionsInput.To, and is useful for accessing the field via an interface.
func (v *__getUserContributionsInput) GetTo() time.Time { return v.To }

// GetFrom returns __getUserContributionsInput.From, and is useful for accessing the field via an interface.
func (v *__getUserContributionsInput) GetFrom() time.Time { return v.From }

func (v *__getUserContributionsInput) UnmarshalJSON(b []byte) error {

	if string(b) == "null" {
		return nil
	}

	var firstPass struct {
		*__getUserContributionsInput
		To   json.RawMessage `json:"to"`
		From json.RawMessage `json:"from"`
		graphql.NoUnmarshalJSON
	}
	firstPass.__getUserContributionsInput = v

	err := json.Unmarshal(b, &firstPass)
	if err != nil {
		return err
	}

	{
		dst := &v.To
		src := firstPass.To
		if len(src) != 0 && string(src) != "null" {
			err = bindings.UnmarshalDateTime(
				src, dst)
			if err != nil {
				return fmt.Errorf(
					"unable to unmarshal __getUserContributionsInput.To: %w", err)
			}
		}
	}

	{
		dst := &v.From
		src := firstPass.From
		if len(src) != 0 && string(src) != "null" {
			err = bindings.UnmarshalDateTime(
				src, dst)
			if err != nil {
				return fmt.Errorf(
					"unable to unmarshal __getUserContributionsInput.From: %w", err)
			}
		}
	}
	return nil
}

type __premarshal__getUserContributionsInput struct {
	UserName string `json:"userName"`

	To json.RawMessage `json:"to"`

	From json.RawMessage `json:"from"`
}

func (v *__getUserContributionsInput) MarshalJSON() ([]byte, error) {
	premarshaled, err := v.__premarshalJSON()
	if err != nil {
		return nil, err
	}
	return json.Marshal(premarshaled)
}

func (v *__getUserContributionsInput) __premarshalJSON() (*__premarshal__getUserContributionsInput, error) {
	var retval __premarshal__getUserContributionsInput

	retval.UserName = v.UserName
	{

		dst := &retval.To
		src := v.To
		var err error
		*dst, err = bindings.MarshalDateTime(
			&src)
		if err != nil {
			return nil, fmt.Errorf(
				"unable to marshal __getUserContributionsInput.To: %w", err)
		}
	}
	{

		dst := &retval.From
		src := v.From
		var err error
		*dst, err = bindings.MarshalDateTime(
			&src)
		if err != nil {
			return nil, fmt.Errorf(
				"unable to marshal __getUserContributionsInput.From: %w", err)
		}
	}
	return &retval, nil
}

// getUserContributionsResponse is returned by getUserContributions on success.
type getUserContributionsResponse struct {
	// Lookup a user by login.
	User getUserContributionsUser `json:"user"`
}

// GetUser returns getUserContributionsResponse.User, and is useful for accessing the field via an interface.
func (v *getUserContributionsResponse) GetUser() getUserContributionsUser { return v.User }

// getUserContributionsUser includes the requested fields of the GraphQL type User.
// The GraphQL type's documentation follows.
//
// A user is an individual's account on GitHub that owns repositories and can make new content.
type getUserContributionsUser struct {
	// The collection of contributions this user has made to different repositories.
	ContributionsCollection getUserContributionsUserContributionsCollection `json:"contributionsCollection"`
}

// GetContributionsCollection returns getUserContributionsUser.ContributionsCollection, and is useful for accessing the field via an interface.
func (v *getUserContributionsUser) GetContributionsCollection() getUserContributionsUserContributionsCollection {
	return v.ContributionsCollection
}

// getUserContributionsUserContributionsCollection includes the requested fields of the GraphQL type ContributionsCollection.
// The GraphQL type's documentation follows.
//
// A contributions collection aggregates contributions such as opened issues and commits created by a user.
type getUserContributionsUserContributionsCollection struct {
	// A calendar of this user's contributions on GitHub.
	ContributionCalendar getUserContributionsUserContributionsCollectionContributionCalendar `json:"contributionCalendar"`
}

// GetContributionCalendar returns getUserContributionsUserContributionsCollection.ContributionCalendar, and is useful for accessing the field via an interface.
func (v *getUserContributionsUserContributionsCollection) GetContributionCalendar() getUserContributionsUserContributionsCollectionContributionCalendar {
	return v.ContributionCalendar
}

// getUserContributionsUserContributionsCollectionContributionCalendar includes the requested fields of the GraphQL type ContributionCalendar.
// The GraphQL type's documentation follows.
//
// A calendar of contributions made on GitHub by a user.
type getUserContributionsUserContributionsCollectionContributionCalendar struct {
	// A list of the weeks of contributions in this calendar.
	Weeks []getUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeek `json:"weeks"`
}

// GetWeeks returns getUserContributionsUserContributionsCollectionContributionCalendar.Weeks, and is useful for accessing the field via an interface.
func (v *getUserContributionsUserContributionsCollectionContributionCalendar) GetWeeks() []getUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeek {
	return v.Weeks
}

// getUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeek includes the requested fields of the GraphQL type ContributionCalendarWeek.
// The GraphQL type's documentation follows.
//
// A week of contributions in a user's contribution graph.
type getUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeek struct {
	// The days of contributions in this week.
	ContributionDays []getUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeekContributionDaysContributionCalendarDay `json:"contributionDays"`
}

// GetContributionDays returns getUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeek.ContributionDays, and is useful for accessing the field via an interface.
func (v *getUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeek) GetContributionDays() []getUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeekContributionDaysContributionCalendarDay {
	return v.ContributionDays
}

// getUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeekContributionDaysContributionCalendarDay includes the requested fields of the GraphQL type ContributionCalendarDay.
// The GraphQL type's documentation follows.
//
// Represents a single day of contributions on GitHub by a user.
type getUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeekContributionDaysContributionCalendarDay struct {
	// The day this square represents.
	Date time.Time `json:"-"`
	// How many contributions were made by the user on this day.
	ContributionCount int `json:"contributionCount"`
}

// GetDate returns getUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeekContributionDaysContributionCalendarDay.Date, and is useful for accessing the field via an interface.
func (v *getUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeekContributionDaysContributionCalendarDay) GetDate() time.Time {
	return v.Date
}

// GetContributionCount returns getUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeekContributionDaysContributionCalendarDay.ContributionCount, and is useful for accessing the field via an interface.
func (v *getUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeekContributionDaysContributionCalendarDay) GetContributionCount() int {
	return v.ContributionCount
}

func (v *getUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeekContributionDaysContributionCalendarDay) UnmarshalJSON(b []byte) error {

	if string(b) == "null" {
		return nil
	}

	var firstPass struct {
		*getUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeekContributionDaysContributionCalendarDay
		Date json.RawMessage `json:"date"`
		graphql.NoUnmarshalJSON
	}
	firstPass.getUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeekContributionDaysContributionCalendarDay = v

	err := json.Unmarshal(b, &firstPass)
	if err != nil {
		return err
	}

	{
		dst := &v.Date
		src := firstPass.Date
		if len(src) != 0 && string(src) != "null" {
			err = bindings.UnmarshalDateTime(
				src, dst)
			if err != nil {
				return fmt.Errorf(
					"unable to unmarshal getUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeekContributionDaysContributionCalendarDay.Date: %w", err)
			}
		}
	}
	return nil
}

type __premarshalgetUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeekContributionDaysContributionCalendarDay struct {
	Date json.RawMessage `json:"date"`

	ContributionCount int `json:"contributionCount"`
}

func (v *getUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeekContributionDaysContributionCalendarDay) MarshalJSON() ([]byte, error) {
	premarshaled, err := v.__premarshalJSON()
	if err != nil {
		return nil, err
	}
	return json.Marshal(premarshaled)
}

func (v *getUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeekContributionDaysContributionCalendarDay) __premarshalJSON() (*__premarshalgetUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeekContributionDaysContributionCalendarDay, error) {
	var retval __premarshalgetUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeekContributionDaysContributionCalendarDay

	{

		dst := &retval.Date
		src := v.Date
		var err error
		*dst, err = bindings.MarshalDateTime(
			&src)
		if err != nil {
			return nil, fmt.Errorf(
				"unable to marshal getUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeekContributionDaysContributionCalendarDay.Date: %w", err)
		}
	}
	retval.ContributionCount = v.ContributionCount
	return &retval, nil
}

// The query or mutation executed by getUserContributions.
const getUserContributions_Operation = `
query getUserContributions ($userName: String!, $to: DateTime!, $from: DateTime!) {
	user(login: $userName) {
		contributionsCollection(to: $to, from: $from) {
			contributionCalendar {
				weeks {
					contributionDays {
						date
						contributionCount
					}
				}
			}
		}
	}
}
`

func getUserContributions(
	ctx context.Context,
	client graphql.Client,
	userName string,
	to time.Time,
	from time.Time,
) (*getUserContributionsResponse, error) {
	req := &graphql.Request{
		OpName: "getUserContributions",
		Query:  getUserContributions_Operation,
		Variables: &__getUserContributionsInput{
			UserName: userName,
			To:       to,
			From:     from,
		},
	}
	var err error

	var data getUserContributionsResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}
