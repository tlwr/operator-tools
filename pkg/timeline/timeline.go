package timeline

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

type TimelineEntry interface {
	Start() time.Time
	End() (time.Time, bool)
	Label() string

	IsEvent() bool
	IsWindow() bool
}

type Timeline interface {
	Start() time.Time
	End() time.Time

	Entries() []TimelineEntry

	Render() (string, error)
}

type timelineEntry struct {
	start time.Time
	end   *time.Time
	label string
}

func (e timelineEntry) Start() time.Time {
	return e.start
}

func (e timelineEntry) End() (time.Time, bool) {
	if e.end == nil {
		return time.Time{}, false
	}
	return *e.end, true
}

func (e timelineEntry) Label() string {
	return e.label
}

func (e timelineEntry) IsEvent() bool {
	return e.end == nil
}

func (e timelineEntry) IsWindow() bool {
	return e.end != nil
}

func NewTimelineEntry(
	start time.Time,
	end *time.Time,
	label string,
) TimelineEntry {
	return timelineEntry{
		start: start,
		end:   end,
		label: label,
	}
}

type timeline struct {
	start time.Time
	end   time.Time

	entries []TimelineEntry

	width int
}

func (t timeline) Start() time.Time {
	return t.start
}

func (t timeline) End() time.Time {
	return t.end
}

func (t timeline) Entries() []TimelineEntry {
	return t.entries
}

func (t timeline) Render() (string, error) {
	var rendered string

	if t.End().Before(t.Start()) {
		return "", fmt.Errorf("timeline %v has end before start", t)
	}

	rendered += t.renderSpine() + "\n"

	for _, entry := range t.sortedEntries() {
		renderedEntry, err := t.renderEntry(entry)
		if err != nil {
			return "", err
		}
		rendered += renderedEntry + "\n"
	}

	return rendered, nil
}

func (t timeline) borderedWidth() int {
	return t.width - len("||")
}

func (t timeline) renderSpine() string {
	return fmt.Sprintf(
		"|%s| total duration %dms",
		strings.Repeat("=", t.borderedWidth()),
		t.End().Sub(t.Start()).Milliseconds(),
	)
}

func (t timeline) renderEntry(entry TimelineEntry) (string, error) {
	if entry.Start().Before(t.Start()) {
		return "", fmt.Errorf("entry %v starts before timeline %v", entry, t)
	}
	totalDistance := float64(t.End().Sub(t.Start()).Milliseconds())

	startDistance := float64(entry.Start().Sub(t.Start()).Milliseconds())
	startPerc := startDistance / totalDistance
	leftPad := int(startPerc * float64(t.borderedWidth()))

	if entry.IsEvent() {
		rightPad := int(math.Max(float64(t.borderedWidth()-leftPad-1), 0))
		eventAfter := entry.Start().Sub(t.Start()).Milliseconds()

		rendered := "|" + strings.Repeat(" ", leftPad) + "x"
		rendered += strings.Repeat(" ", rightPad) + "| "
		rendered += fmt.Sprintf("%s at %dms", entry.Label(), eventAfter)
		return rendered, nil
	}

	windowEnd, _ := entry.End() // checked with IsEvent, windowEnd not nil

	endDistance := float64(windowEnd.Sub(t.Start()).Milliseconds())
	windowDistance := endDistance - startDistance
	windowPerc := windowDistance / totalDistance
	windowWidth := int(windowPerc * float64(t.borderedWidth()))

	rightPad := t.borderedWidth() - windowWidth - leftPad
	windowStartAfter := entry.Start().Sub(t.Start()).Milliseconds()
	windowEndAfter := windowEnd.Sub(t.Start()).Milliseconds()
	duration := windowEndAfter - windowStartAfter

	rendered := "|" + strings.Repeat(" ", leftPad)
	rendered += strings.Repeat("~", windowWidth)
	rendered += strings.Repeat(" ", rightPad) + "| "
	rendered += entry.Label() + " "
	rendered += fmt.Sprintf("from %dms ", windowStartAfter)
	rendered += fmt.Sprintf("until %dms ", windowEndAfter)
	rendered += fmt.Sprintf("duration %dms", duration)
	return rendered, nil
}

func (t timeline) sortedEntries() []TimelineEntry {
	entries := append(make([]TimelineEntry, 0), t.Entries()...)
	sort.Slice(entries, func(index1 int, index2 int) bool {
		firstEntry := entries[index1]
		secondEntry := entries[index2]

		// if they start at the same time
		if firstEntry.Start().Equal(secondEntry.Start()) {
			// and if they are both events, it doesn't matter which is first
			if firstEntry.IsEvent() && secondEntry.IsEvent() {
				return true
			}

			// Windows come before events, if they start at the same time
			if secondEntry.IsWindow() {
				return false
			}

			return true
		}

		return firstEntry.Start().Before(secondEntry.Start())
	})
	return entries
}

func NewTimeline(
	start time.Time,
	end time.Time,
	width int,
	entries ...TimelineEntry,
) Timeline {
	return timeline{
		start: start,
		end:   end,

		entries: entries,

		width: width,
	}
}
