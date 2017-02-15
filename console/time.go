/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package console

import (
	"fmt"
	"time"
)

func humanizeTime(t time.Time) string {
	d := time.Now().UTC().Sub(t)
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
