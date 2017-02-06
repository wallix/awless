package console

import (
	"fmt"
	"time"
)

func humanizeTime(t time.Time) string {
	d := time.Now().Sub(t)
	switch {
	case d.Seconds() <= time.Second.Seconds():
		return "now"
	case d.Seconds() <= 2*60*time.Second.Seconds():
		return fmt.Sprintf("%d seconds ago", int(d.Seconds()))
	case d.Seconds() <= 2*60*time.Minute.Seconds():
		return fmt.Sprintf("%d minutes ago", int(d.Minutes()))
	case d.Seconds() <= 2*24*time.Hour.Seconds():
		return fmt.Sprintf("%d hours ago", int(d.Hours()))
	case d.Seconds() <= 2*7*24*time.Hour.Seconds():
		return fmt.Sprintf("%d days ago", int(d.Hours()/24))
	case d.Seconds() <= 2*30*24*time.Hour.Seconds():
		return fmt.Sprintf("%d weeks ago", int(d.Hours()/(24*7)))
	case d.Seconds() <= 2*365*24*time.Hour.Seconds():
		return fmt.Sprintf("%d months ago", int(d.Hours()/(24*30)))
	default:
		return fmt.Sprintf("%d years ago", int(d.Hours()/(24*365)))
	}
}
