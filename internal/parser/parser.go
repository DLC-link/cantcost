package parser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/DLC-link/cantcost/internal/env"
)

// Recipient represents either a MemberRecipient or MediatorGroupRecipient
type Recipient struct {
	Type    string `json:"type"`
	Member  string `json:"member"`
	GroupID int    `json:"group_id"`
}

// EnvelopeCostDetails represents the cost details for an envelope
type EnvelopeCostDetails struct {
	WriteCost  int         `json:"write_cost"`
	ReadCost   int         `json:"read_cost"`
	FinalCost  int         `json:"final_cost"`
	Recipients []Recipient `json:"recipients"`
}

// EventCostDetails represents the parsed cost details from the log message
type EventCostDetails struct {
	EventCost          int                   `json:"event_cost"`
	CostMultiplier     int                   `json:"cost_multiplier"`
	GroupToMembersSize map[int]int           `json:"group_to_members_size"`
	EnvelopesCost      []EnvelopeCostDetails `json:"envelopes_cost"`
}

type Line struct {
	// DockerTimestamp is the timestamp from the Docker log prefix
	DockerTimestamp time.Time `json:"-"`

	// Fields from the JSON payload
	Timestamp    time.Time `json:"@timestamp"`
	Message      string    `json:"message"`
	LoggerName   string    `json:"logger_name"`
	ThreadName   string    `json:"thread_name"`
	Level        string    `json:"level"`
	SpanID       string    `json:"span-id"`
	SpanParentID string    `json:"span-parent-id"`
	TraceID      string    `json:"trace-id"`
	SpanName     string    `json:"span-name"`

	// Parsed from Message
	CostDetails *EventCostDetails `json:"-"`
}

type MessageLine struct {
	// DockerTimestamp is the timestamp from the Docker log prefix
	DockerTimestamp time.Time `json:"-"`

	// Fields from the JSON payload
	Timestamp    time.Time `json:"@timestamp"`
	Message      string    `json:"message"`
	LoggerName   string    `json:"logger_name"`
	ThreadName   string    `json:"thread_name"`
	Level        string    `json:"level"`
	SpanID       string    `json:"span_id"`
	SpanParentID string    `json:"span_parent_id"`
	TraceID      string    `json:"trace_id"`
	SpanName     string    `json:"span_name"`

	// Parsed from Message
	CostDetails *EventCostDetails `json:"cost_details"`
}

func ProcessLine(line string) (Line, error) {
	// Find the first space which separates the Docker timestamp from the JSON payload
	spaceIdx := strings.Index(line, " ")
	if spaceIdx == -1 {
		return Line{}, fmt.Errorf("invalid line format: no space separator found")
	}

	dockerTimestampStr := line[:spaceIdx]
	jsonPayload := line[spaceIdx+1:]

	// Parse the Docker timestamp (RFC3339Nano format)
	dockerTimestamp, err := time.Parse(time.RFC3339Nano, dockerTimestampStr)
	if err != nil {
		return Line{}, fmt.Errorf("failed to parse docker timestamp: %w", err)
	}

	// Parse the JSON payload
	var l Line
	if err := json.Unmarshal([]byte(jsonPayload), &l); err != nil {
		return Line{}, fmt.Errorf("failed to parse JSON payload: %w", err)
	}

	l.DockerTimestamp = dockerTimestamp

	// Parse EventCostDetails from the message if present
	if strings.Contains(l.Message, "EventCostDetails(") {
		costDetails, err := parseEventCostDetails(l.Message)
		if err != nil {
			return Line{}, fmt.Errorf("failed to parse EventCostDetails: %w", err)
		}
		l.CostDetails = costDetails
	}

	return l, nil
}

func (l *Line) ToMessageLine() *MessageLine {
	message := &MessageLine{
		DockerTimestamp: l.DockerTimestamp,
		Timestamp:       l.Timestamp,
		LoggerName:      l.LoggerName,
		ThreadName:      l.ThreadName,
		Level:           l.Level,
		SpanID:          l.SpanID,
		SpanParentID:    l.SpanParentID,
		TraceID:         l.TraceID,
		SpanName:        l.SpanName,
		CostDetails:     l.CostDetails,
	}
	if env.GetIncludeMessage() {
		message.Message = l.Message
	}
	return message
}

func parseEventCostDetails(message string) (*EventCostDetails, error) {
	details := &EventCostDetails{
		GroupToMembersSize: make(map[int]int),
	}

	// Extract event cost
	eventCostRe := regexp.MustCompile(`event cost = (\d+)`)
	if match := eventCostRe.FindStringSubmatch(message); len(match) > 1 {
		details.EventCost, _ = strconv.Atoi(match[1])
	}

	// Extract cost multiplier
	multiplierRe := regexp.MustCompile(`cost multiplier = (\d+)`)
	if match := multiplierRe.FindStringSubmatch(message); len(match) > 1 {
		details.CostMultiplier, _ = strconv.Atoi(match[1])
	}

	// Extract group to members size: MediatorGroupRecipient(group = 0) -> 14
	groupSizeRe := regexp.MustCompile(`MediatorGroupRecipient\(group = (\d+)\) -> (\d+)`)
	if match := groupSizeRe.FindStringSubmatch(message); len(match) > 2 {
		groupID, _ := strconv.Atoi(match[1])
		size, _ := strconv.Atoi(match[2])
		details.GroupToMembersSize[groupID] = size
	}

	// Extract envelope cost details
	details.EnvelopesCost = parseEnvelopeCostDetails(message)

	return details, nil
}

func parseEnvelopeCostDetails(message string) []EnvelopeCostDetails {
	var envelopes []EnvelopeCostDetails

	// Find the start of envelopes cost details section
	envelopesStart := strings.Index(message, "envelopes cost details = ")
	if envelopesStart == -1 {
		return envelopes
	}

	envelopesSection := message[envelopesStart:]

	// Pattern to match EnvelopeCostDetails
	// Can be either a single one or a Seq(...)
	envelopeRe := regexp.MustCompile(`EnvelopeCostDetails\(`)
	matches := envelopeRe.FindAllStringIndex(envelopesSection, -1)

	for _, match := range matches {
		startIdx := match[0]
		envelope := extractEnvelopeCostDetails(envelopesSection[startIdx:])
		if envelope != nil {
			envelopes = append(envelopes, *envelope)
		}
	}

	return envelopes
}

func extractEnvelopeCostDetails(s string) *EnvelopeCostDetails {
	envelope := &EnvelopeCostDetails{}

	// Extract write cost
	writeCostRe := regexp.MustCompile(`write cost = (\d+)`)
	if match := writeCostRe.FindStringSubmatch(s); len(match) > 1 {
		envelope.WriteCost, _ = strconv.Atoi(match[1])
	}

	// Extract read cost
	readCostRe := regexp.MustCompile(`read cost = (\d+)`)
	if match := readCostRe.FindStringSubmatch(s); len(match) > 1 {
		envelope.ReadCost, _ = strconv.Atoi(match[1])
	}

	// Extract final cost
	finalCostRe := regexp.MustCompile(`final cost = (\d+)`)
	if match := finalCostRe.FindStringSubmatch(s); len(match) > 1 {
		envelope.FinalCost, _ = strconv.Atoi(match[1])
	}

	// Extract recipients - find the recipients section for this envelope
	recipientsStart := strings.Index(s, "recipients = ")
	if recipientsStart == -1 {
		return envelope
	}

	// Find where this envelope's recipients end (before the next EnvelopeCostDetails or end of envelope)
	recipientsSection := s[recipientsStart:]

	// Check if it's a Seq or a single recipient
	if strings.HasPrefix(recipientsSection, "recipients = Seq(") {
		// Find the matching closing paren for the Seq
		seqContent := extractBalancedParens(recipientsSection[len("recipients = Seq"):])
		envelope.Recipients = parseRecipients(seqContent)
	} else {
		// Single recipient - extract until the closing paren of the envelope
		singleRecipient := recipientsSection[len("recipients = "):]
		// Find end - either ) for closing envelope or ,
		endIdx := strings.Index(singleRecipient, ")")
		if endIdx != -1 {
			// Include the closing paren if it's part of the recipient
			recipientStr := singleRecipient[:endIdx+1]
			envelope.Recipients = parseRecipients(recipientStr)
		}
	}

	return envelope
}

func extractBalancedParens(s string) string {
	if len(s) == 0 || s[0] != '(' {
		return ""
	}

	depth := 0
	for i, c := range s {
		if c == '(' {
			depth++
		} else if c == ')' {
			depth--
			if depth == 0 {
				return s[1:i] // Return content inside the parens
			}
		}
	}
	return s[1:] // Return everything if unbalanced
}

func parseRecipients(s string) []Recipient {
	var recipients []Recipient

	// Match MemberRecipient(PAR::name::hash...)
	memberRe := regexp.MustCompile(`MemberRecipient\(([^)]+)\)`)
	memberMatches := memberRe.FindAllStringSubmatch(s, -1)
	for _, match := range memberMatches {
		if len(match) > 1 {
			recipients = append(recipients, Recipient{
				Type:   "MemberRecipient",
				Member: match[1],
			})
		}
	}

	// Match MediatorGroupRecipient(group = N)
	mediatorRe := regexp.MustCompile(`MediatorGroupRecipient\(group = (\d+)\)`)
	mediatorMatches := mediatorRe.FindAllStringSubmatch(s, -1)
	for _, match := range mediatorMatches {
		if len(match) > 1 {
			groupID, _ := strconv.Atoi(match[1])
			recipients = append(recipients, Recipient{
				Type:    "MediatorGroupRecipient",
				GroupID: groupID,
			})
		}
	}

	return recipients
}
