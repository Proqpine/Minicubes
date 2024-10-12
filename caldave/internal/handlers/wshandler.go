package handlers

import (
	"caldave/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/websocket"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type MessageType string

const (
	RequestAvailability  MessageType = "REQUEST_AVAILABILITY"
	AvailabilityResponse MessageType = "AVAILABILITY_RESPONSE"
	UpdateAvailaibilty   MessageType = "UPDATE_AVAILABILITY"
	EventUpdated         MessageType = "EVENTS_UPDATED"
)

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type AvailabilityRequest struct {
	Date string `json:"date"` // Format: "2024-10-11"
}

type TimeSlot struct {
	Start string `json:"start"` // Format: "HH:MM"
	End   string `json:"end"`   // Format: "HH:MM"
}

type BusinessHours struct {
	StartTime string // Business start time, e.g., "08:00"
	EndTime   string // Business end time, e.g., "17:00"
}

type ScheduleConfig struct {
	WeekdayHours  map[time.Weekday]BusinessHours // Custom business hours for specific weekdays
	DefaultHours  BusinessHours                  // Default business hours
	BufferMinutes int                            // Buffer time around events, in minutes
}

type AvailabilityResponseData struct {
	Date           string     `json:"date"`
	AvailableTimes []TimeSlot `json:"availableTimes"`
}

type UpdateEventsRequest struct {
	StartDate string `json:"startDate"` // Format: "2024-10-11"
	EndDate   string `json:"endDate"`   // Format: "2024-11-11"
}

type Client struct {
	ID         string
	Connection *websocket.Conn
	Hub        *Hub
	Send       chan Message
}

type Hub struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
	mutex      sync.RWMutex
	wsHandler  *WebSocketHandler
}

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	hub             *Hub
	calendarService *calendar.Service
	events          []utils.EventData
}

func NewHub(handler *WebSocketHandler) *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message, 256),
		wsHandler:  handler,
	}
}

func NewWebSocketHandler() *WebSocketHandler {
	handler := &WebSocketHandler{}
	hub := NewHub(handler)
	handler.hub = hub

	go hub.Run()

	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := utils.GetClient(config)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	handler.calendarService = srv

	handler.updateEvents()

	go handler.refreshEvents()

	return handler
}

// Run starts the Hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mutex.Lock()
			h.Clients[client.ID] = client
			h.mutex.Unlock()
			log.Printf("Client %s connected", client.ID)

		case client := <-h.Unregister:
			h.mutex.Lock()
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.Send)
			}
			h.mutex.Unlock()
			log.Printf("Client %s disconnected", client.ID)

		case message := <-h.Broadcast:
			h.mutex.RLock()
			for _, client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client.ID)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.Connection.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				return
			}

			err := websocket.JSON.Send(c.Connection, message)
			if err != nil {
				log.Printf("Error sending message to client %s: %v", c.ID, err)
				return
			}
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Connection.Close()
	}()

	for {
		var message Message
		err := websocket.JSON.Receive(c.Connection, &message)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		switch message.Type {
		case string(RequestAvailability):
			c.handleAvailabilityRequest(message)
		case string(UpdateAvailaibilty):
			c.handleUpdateEventsRequest(message)
		default:
			c.Hub.Broadcast <- message
		}

	}
}

// Add method to handle availability requests
func (c *Client) handleAvailabilityRequest(message Message) {
	reqBytes, _ := json.Marshal(message.Payload)
	var request AvailabilityRequest
	if err := json.Unmarshal(reqBytes, &request); err != nil {
		log.Printf("Error parsing availability request: %v", err)
		return
	}

	datePart := strings.Split(request.Date, "T")[0]

	requestedDate, err := time.Parse("2006-01-02", datePart)
	if err != nil {
		log.Printf("Error parsing date: %v", err)
		return
	}

	handler := c.Hub.wsHandler

	schedule := ScheduleConfig{
		WeekdayHours: map[time.Weekday]BusinessHours{
			time.Monday:    {StartTime: "08:00", EndTime: "17:00"},
			time.Tuesday:   {StartTime: "08:30", EndTime: "17:00"},
			time.Wednesday: {StartTime: "09:00", EndTime: "18:00"},
			time.Thursday:  {StartTime: "08:00", EndTime: "17:00"},
			time.Friday:    {StartTime: "08:00", EndTime: "17:30"},
			time.Saturday:  {StartTime: "10:00", EndTime: "14:00"}, // Optional business hours on weekends
			time.Sunday:    {StartTime: "00:00", EndTime: "00:00"}, // Closed on Sundays
		},
		DefaultHours: BusinessHours{
			StartTime: "08:00",
			EndTime:   "17:30", // Default fallback if no specific day is defined
		},
		BufferMinutes: 10, // 10-minute buffer before and after events
	}

	availableTimes := getAvailableTimesForDate(requestedDate, handler.events, schedule)

	response := Message{
		Type: string(AvailabilityResponse),
		Payload: AvailabilityResponseData{
			Date:           request.Date,
			AvailableTimes: availableTimes,
		},
	}

	c.Send <- response
}

func (c *Client) handleUpdateEventsRequest(message Message) {
	reqBytes, _ := json.Marshal(message.Payload)
	var request UpdateEventsRequest
	if err := json.Unmarshal(reqBytes, &request); err != nil {
		log.Printf("Error parsing update events request: %v", err)
		return
	}

	handler := c.Hub.wsHandler

	startDate, _ := time.Parse("2006-01-02", request.StartDate)
	endDate, _ := time.Parse("2006-01-02", request.EndDate)

	calendars := utils.GetCalendars(handler.calendarService)
	handler.events = utils.GetEvents(startDate.Format(time.RFC3339), endDate.Format(time.RFC3339), handler.calendarService, calendars)

	response := Message{
		Type:    string(EventUpdated),
		Payload: nil,
	}

	c.Send <- response
}

// Function to calculate available times
func getAvailableTimesForDate(date time.Time, events []utils.EventData, config ScheduleConfig) []TimeSlot {
	// Determine the business hours for the given date
	hours, ok := config.WeekdayHours[date.Weekday()]
	if !ok {
		hours = config.DefaultHours
	}

	fmt.Println(hours)

	// Parse business start and end times
	businessStart, _ := time.Parse("15:04", hours.StartTime)
	businessEnd, _ := time.Parse("15:04", hours.EndTime)

	// Apply the date to the business start and end times
	businessStart = time.Date(date.Year(), date.Month(), date.Day(), businessStart.Hour(), businessStart.Minute(), 0, 0, date.Location())
	businessEnd = time.Date(date.Year(), date.Month(), date.Day(), businessEnd.Hour(), businessEnd.Minute(), 0, 0, date.Location())

	// Filter events for the day
	var dayEvents []utils.EventData
	for _, event := range events {
		if event.StartTime.Year() == date.Year() &&
			event.StartTime.Month() == date.Month() &&
			event.StartTime.Day() == date.Day() {
			dayEvents = append(dayEvents, event)
		}
	}
	// fmt.Printf("%# v\n", pretty.Formatter(events))
	// fmt.Printf("%# v\n", pretty.Formatter(dayEvents))

	// Sort events by start time
	sort.Slice(dayEvents, func(i, j int) bool {
		return dayEvents[i].StartTime.Before(dayEvents[j].StartTime)
	})

	// Find available slots
	var availableSlots []TimeSlot
	currentTime := businessStart

	for _, event := range dayEvents {
		// Apply buffer to event start and end times
		eventStart := event.StartTime.Add(-time.Duration(config.BufferMinutes) * time.Minute)
		eventEnd := event.EndTime.Add(time.Duration(config.BufferMinutes) * time.Minute)

		// Check for available slot before the event
		if currentTime.Before(eventStart) {
			availableSlots = append(availableSlots, TimeSlot{
				Start: currentTime.Format("15:04"),
				End:   eventStart.Format("15:04"),
			})
		}

		// Move current time to event end if it's later than the current time
		if eventEnd.After(currentTime) {
			currentTime = eventEnd
		}
	}

	// Add remaining time until business end if available
	if currentTime.Before(businessEnd) {
		availableSlots = append(availableSlots, TimeSlot{
			Start: currentTime.Format("15:04"),
			End:   businessEnd.Format("15:04"),
		})
	}

	return availableSlots
}

func (wsh *WebSocketHandler) refreshEvents() {
	ticker := time.NewTicker(15 * time.Minute)
	for {
		select {
		case <-ticker.C:
			wsh.updateEvents()
		}
	}
}

func (wsh *WebSocketHandler) updateEvents() {

	startDay := time.Now().AddDate(0, 0, -30).Format(time.RFC3339)
	endDay := time.Now().AddDate(0, 0, 60).Format(time.RFC3339)

	calendars := utils.GetCalendars(wsh.calendarService)
	wsh.events = utils.GetEvents(startDay, endDay, wsh.calendarService, calendars)
}

func (wsh *WebSocketHandler) HandleWS(ws *websocket.Conn) {
	client := &Client{
		ID:         ws.RemoteAddr().String(),
		Connection: ws,
		Hub:        wsh.hub,
		Send:       make(chan Message, 256),
	}

	wsh.hub.Register <- client

	go client.WritePump()
	client.ReadPump()
}

func (wsh *WebSocketHandler) Handler() http.Handler {
	return websocket.Handler(wsh.HandleWS)
}
