package awsservices

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/wallix/awless/console"
)

var DefaultNetworkMonitor = &NetworkMonitor{requests: make(map[*request.Request]*req)}

type NetworkMonitor struct {
	requests map[*request.Request]*req
	l        sync.Mutex
}

type req struct {
	*request.Request
	from time.Time
	to   time.Time
}

func (n *NetworkMonitor) DisplayStats(w io.Writer) {
	fmt.Fprintf(w, "\n%d requests sent:\n", len(n.requests))

	var sorted []*req

	var min, max time.Time
	var maxFunctionNameLength int
	for _, r := range n.requests {
		if min.IsZero() || r.from.Before(min) {
			min = r.from
		}
		if max.IsZero() || r.to.After(max) {
			max = r.to
		}
		if len(r.Operation.Name) > maxFunctionNameLength {
			maxFunctionNameLength = len(r.Operation.Name)
		}
		sorted = append(sorted, r)
	}
	sort.Slice(sorted, func(i int, j int) bool {
		if sorted[i].from.Equal(sorted[j].from) {
			return sorted[i].to.Before(sorted[j].to)
		}
		return sorted[i].from.Before(sorted[j].from)
	})

	maxwidth := uint(console.GetTerminalWidth() - maxFunctionNameLength - 11) // 11 = '['+']'+' '+'('+4+'m'+'s'+')'
	maxDuration := max.Sub(min)

	for _, r := range sorted {
		duration := r.to.Sub(r.from)
		width := uint(duration) * maxwidth / uint(maxDuration)
		before := uint(r.from.Sub(min)) * maxwidth / uint(maxDuration)
		after := maxwidth - width - before
		fmt.Fprintf(w, "%s[%s]%s %s(%dms)\n", strings.Repeat(" ", int(before)), strings.Repeat("-", int(width)), strings.Repeat(" ", int(after)), r.Operation.Name, duration/(1000*1000))
	}
}

func (n *NetworkMonitor) addRequest(r *request.Request) {
	n.l.Lock()
	defer n.l.Unlock()
	n.requests[r] = &req{Request: r, from: time.Now().UTC()}
}

func (n *NetworkMonitor) setRequestEnd(r *request.Request) {
	n.l.Lock()
	defer n.l.Unlock()
	request, ok := n.requests[r]
	if !ok {
		fmt.Fprintf(os.Stderr, "request '%s' not found\n", r.RequestID)
		return
	}
	request.to = time.Now().UTC()
}
