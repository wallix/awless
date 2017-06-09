package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/chzyer/readline"
	"github.com/wallix/awless/aws"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/template"
)

func holeAutoCompletion(g *graph.Graph, hole string) readline.AutoCompleter {
	completeFunc := func(string) []string { return []string{} }

	if entityType, entityProp := guessEntityTypeFromHoleQuestion(hole); entityType != "" {
		resources, err := g.GetAllResources(entityType)
		exitOn(err)

		if len(resources) == 0 {
			return readline.NewPrefixCompleter(readline.PcItemDynamic(completeFunc))
		}

		if entityProp != "" {
			var validPropName string
			for _, r := range resources {
				for key := range r.Properties {
					if strings.ToLower(key) == strings.ToLower(entityProp) {
						validPropName = key
					}
				}
			}
			if validPropName != "" {
				completeFunc = func(s string) (suggest []string) {
					for _, res := range resources {
						if val, ok := res.Properties[validPropName].(string); ok {
							suggest = append(suggest, val)
						}
					}
					suggest = quotedSortedSet(suggest)
					return
				}
			}

			return readline.NewPrefixCompleter(readline.PcItemDynamic(completeFunc))
		}

		completeFunc = func(s string) (suggest []string) {
			for _, res := range resources {
				id := res.Id()
				if strings.Contains(id, s) {
					suggest = append(suggest, id)
				}
				if val, ok := res.Properties["Name"]; ok {
					switch val.(type) {
					case string:
						name := val.(string)
						prefixed := fmt.Sprintf("@%s", name)
						if strings.Contains(prefixed, s) && name != "" {
							suggest = append(suggest, prefixed)
						}
					}
				}
			}

			suggest = quotedSortedSet(suggest)
			return
		}
	}

	return readline.NewPrefixCompleter(readline.PcItemDynamic(completeFunc))
}

func quotedSortedSet(list []string) (out []string) {
	unique := make(map[string]bool)
	for _, l := range list {
		unique[l] = true
	}

	for k := range unique {
		if !template.MatchStringParamValue(k) {
			k = "'" + k + "'"
		}

		out = append(out, k)
	}

	sort.Strings(out)
	return
}

// Return 2 strings according to holes questions:
// notmatching? -> ["", ""]
// entitytype? -> [entitytype, ""]
// entitytype.anything? -> [entitytype, anything]
// entitytype.otherentitytype? -> [otherentitytype, ""]
func guessEntityTypeFromHoleQuestion(hole string) (string, string) {
	var types []string
	splits := strings.Split(hole, ".")
	for _, t := range splits {
		for _, r := range aws.ResourceTypes {
			if t == r {
				types = append(types, r)
				break
			}
		}
	}

	var prop string
	if len(splits) == 2 {
		prop = splits[1]
	}

	if l := len(types); l == 1 {
		return types[0], prop
	} else if l == 2 {
		return types[1], ""
	}

	return "", ""
}
