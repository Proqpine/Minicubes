package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
)

type CalendarData struct {
	CalendarID   string `json:"id,omitempty"`
	CalendarName string `json:"summary,omitempty"`
}

type EventData struct {
	Calendar  CalendarData
	EventName string `json:"summary,omitempty"`
	StartTime time.Time
	EndTime   time.Time
}

type Availaibility struct {
	Day       time.Weekday
	StartTime time.Time
	EndTime   time.Time
}

// Retrieve a token, saves the token, then returns the generated client.
func GetClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// func RunCalendarService() {
// 	cldDta := getCalendars(srv)
// 	startDay := time.Now().AddDate(0, 0, -30).Format(time.RFC3339)
// 	endDay := time.Now().AddDate(0, 0, 30).Format(time.RFC3339)
// 	getEvents(startDay, endDay, srv, cldDta)

// }

func GetCalendars(srv *calendar.Service) []CalendarData {
	calendarList := calendar.NewCalendarListService(srv)
	lst, err := calendarList.List().Do()
	if err != nil {
		log.Fatalf("Unable to retrieve Calendars: %v", err)
	}
	var data []CalendarData
	for _, item := range lst.Items {
		tmp := CalendarData{
			CalendarID:   item.Id,
			CalendarName: item.Summary,
		}
		data = append(data, tmp)
	}
	return data
}

func GetEvents(startDay, endDay string, srv *calendar.Service, cldData []CalendarData) []EventData {
	var calendarEvents []EventData

	for _, cld := range cldData {
		events, err := srv.Events.List(cld.CalendarID).ShowDeleted(true).
			SingleEvents(false).TimeMin(startDay).TimeMax(endDay).MaxResults(10).Do()
		if err != nil {
			log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
		}
		if len(events.Items) == 0 {

		} else {
			for _, item := range events.Items {
				date := item.Start.DateTime
				endDate := item.End.DateTime
				if date == "" {
					date = item.Start.Date
				}
				if endDate == "" {
					endDate = item.End.Date
				}

				parsedStartTime, err := parseDateTime(date)
				if err != nil {
					// Handle error
					log.Printf("Error parsing start time: %v", err)
					continue
				}

				parsedEndTime, err := parseDateTime(endDate)
				if err != nil {
					// Handle error
					log.Printf("Error parsing end time: %v", err)
					continue
				}

				event := EventData{
					Calendar:  CalendarData{CalendarID: cld.CalendarID, CalendarName: cld.CalendarName},
					EventName: item.Summary,
					StartTime: parsedStartTime,
					EndTime:   parsedEndTime,
				}

				calendarEvents = append(calendarEvents, event)
			}
		}
	}
	return calendarEvents
}

func parseDateTime(datetime string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, datetime)
	if err == nil {
		return t, nil
	}

	t, err = time.Parse("2006-01-02", datetime)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func getAvailableTimes(bookedEvents []EventData) {

}

// Parse the events and get all the times and dates
// Then for each of them sort them by date/day of the week
// Extract the time and group them
// Then return the times between the availability that are not events
