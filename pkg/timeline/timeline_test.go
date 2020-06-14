package timeline_test

import (
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/tlwr/operator-tools/pkg/timeline"
)

var _ = Describe("Timeline", func() {
	Context("when there are no entries", func() {
		It("renders just the spine", func() {
			width := 10
			start := time.Now()
			end := start.Add(1000 * time.Millisecond)

			t := timeline.NewTimeline(start, end, width)
			r, err := t.Render()

			Expect(err).NotTo(HaveOccurred())
			Expect(r).To(Equal("|========| total duration 1000ms\n"))
		})
	})

	Context("when there is one entry", func() {
		Context("and the single entry is an event", func() {
			It("renders the spine and the single event", func() {
				width := 10
				start := time.Now()
				end := start.Add(1000 * time.Millisecond)
				midpoint := start.Add(500 * time.Millisecond)

				t := timeline.NewTimeline(
					start, end, width,
					timeline.NewTimelineEntry(midpoint, nil, "midpoint"),
				)
				r, err := t.Render()

				Expect(err).NotTo(HaveOccurred())

				lines := strings.Split(strings.TrimSpace(r), "\n")
				Expect(lines).To(HaveLen(2))
				Expect(lines[0]).To(Equal("|========| total duration 1000ms"))
				Expect(lines[1]).To(Equal("|    x   | midpoint at 500ms"))
			})
		})
	})

	Context("when there are two entries", func() {
		Context("and both entries are events", func() {
			It("renders the spine and the events are rendered in asc order", func() {
				width := 10
				start := time.Now()
				end := start.Add(1000 * time.Millisecond)
				firstQuartile := start.Add(250 * time.Millisecond)
				thirdQuartile := start.Add(750 * time.Millisecond)

				t := timeline.NewTimeline(
					start, end, width,
					timeline.NewTimelineEntry(firstQuartile, nil, "first-quartile"),
					timeline.NewTimelineEntry(thirdQuartile, nil, "third-quartile"),
				)
				r, err := t.Render()

				Expect(err).NotTo(HaveOccurred())

				lines := strings.Split(strings.TrimSpace(r), "\n")
				Expect(lines).To(HaveLen(3))
				Expect(lines[0]).To(Equal("|========| total duration 1000ms"))
				Expect(lines[1]).To(Equal("|  x     | first-quartile at 250ms"))
				Expect(lines[2]).To(Equal("|      x | third-quartile at 750ms"))
			})
		})

		Context("and one entry is a window and another is an event", func() {
			It("renders the spine and the events are rendered in asc order", func() {
				width := 10
				start := time.Now()
				end := start.Add(1000 * time.Millisecond)
				firstQuartile := start.Add(250 * time.Millisecond)
				midpoint := start.Add(500 * time.Millisecond)
				thirdQuartile := start.Add(750 * time.Millisecond)

				t := timeline.NewTimeline(
					start, end, width,
					timeline.NewTimelineEntry(firstQuartile, &thirdQuartile, "middle"),
					timeline.NewTimelineEntry(midpoint, nil, "midpoint"),
				)
				r, err := t.Render()

				Expect(err).NotTo(HaveOccurred())

				lines := strings.Split(strings.TrimSpace(r), "\n")
				Expect(lines).To(HaveLen(3))
				Expect(lines[0]).To(Equal("|========| total duration 1000ms"))
				Expect(lines[1]).To(Equal("|  ~~~~  | middle from 250ms until 750ms duration 500ms"))
				Expect(lines[2]).To(Equal("|    x   | midpoint at 500ms"))
			})
		})
	})
})
