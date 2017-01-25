package sync

type revsByDate []*Rev

func (r revsByDate) Len() int           { return len(r) }
func (r revsByDate) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r revsByDate) Less(i, j int) bool { return r[i].Date.Before(r[j].Date) }
