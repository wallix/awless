package params

type Reducer interface {
	Keys() []string
	Reduce(map[string]interface{}) (map[string]interface{}, error)
}

type reduceFunc func(map[string]interface{}) (map[string]interface{}, error)

func newReducer(fn reduceFunc, keys ...string) Reducer {
	return &reducer{reduce: fn, keys: keys}
}

type reducer struct {
	keys   []string
	reduce reduceFunc
}

func (r *reducer) Keys() []string {
	return r.keys
}

func (r *reducer) Reduce(all map[string]interface{}) (map[string]interface{}, error) {
	in := make(map[string]interface{})
	for _, k := range r.keys {
		if v, ok := all[k]; ok {
			in[k] = v
		}
	}
	return r.reduce(in)
}
