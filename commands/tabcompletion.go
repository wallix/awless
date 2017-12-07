package commands

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/fxaguessy/readline"
	"github.com/wallix/awless/aws/services"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/template"
)

func enumCompletionFunc(enum []string) readline.AutoCompleter {
	var items []readline.PrefixCompleterInterface
	for _, e := range enum {
		items = append(items, readline.PcItem(e))
	}
	return readline.NewPrefixCompleter(items...)
}

func typedParamCompletionFunc(g *graph.Graph, resourceType, propName string) readline.AutoCompleter {
	var items []readline.PrefixCompleterInterface
	resources, _ := g.GetAllResources(resourceType)
	for _, res := range resources {
		if val, ok := res.Properties()[propName]; ok {
			items = append(items, readline.PcItem(fmt.Sprint(val)))
		}
	}

	return readline.NewPrefixCompleter(items...)
}
func holeAutoCompletion(g *graph.Graph, hole string) readline.AutoCompleter {
	completeFunc := func(string) []string { return []string{} }

	if entityTypes, entityProp := guessEntityTypeFromHoleQuestion(hole); len(entityTypes) > 0 {
		var resources []*graph.Resource
		res, err := g.GetAllResources(entityTypes...)
		resources = append(resources, res...)
		exitOn(err)

		if len(resources) == 0 {
			return &prefixCompleter{callback: completeFunc}
		}

		if entityProp != "" {
			var validPropName string
			for _, r := range resources {
				for propName := range r.Properties() {
					if keyCorrespondsToProperty(entityProp, propName) {
						validPropName = propName
					}
				}
			}

			if validPropName != "" {
				completeFunc = func(s string) (suggest []string) {
					for _, res := range resources {
						if v, ok := res.Properties()[validPropName]; ok {
							switch prop := v.(type) {
							case string, float64, int, bool:
								suggest = appendIfContains(suggest, fmt.Sprint(prop), s)
								if validPropName == "ID" {
									suggest = appendWithNameAliases(suggest, res, s)
								}
							case []string:
								for _, str := range prop {
									suggest = appendIfContains(suggest, str, s)
								}
							case []*graph.KeyValue:
								for _, kv := range prop {
									suggest = appendIfContains(suggest, fmt.Sprintf("%s:%s", kv.KeyName, kv.Value), s)
								}
							}
						}
					}
					suggest = quotedSortedSet(suggest)
					return
				}
				return &prefixCompleter{callback: completeFunc, splitChar: ","}
			}
		}

		completeFunc = func(s string) (suggest []string) {
			s = splitKeepLast(s, ",")
			s = strings.TrimLeft(s, "'@\"")
			for _, res := range resources {
				suggest = appendIfContains(suggest, res.Id(), s)
				suggest = appendWithNameAliases(suggest, res, s)
			}
			suggest = quotedSortedSet(suggest)
			return
		}
	}

	return &prefixCompleter{callback: completeFunc, splitChar: ","}
}

type prefixCompleter struct {
	callback  readline.DynamicCompleteFunc
	splitChar string
}

func (p *prefixCompleter) Do(line []rune, pos int) (newLine [][]rune, offset int) {
	var lines []string
	lines, offset = doInternal(p, string(line), pos, line)
	for _, l := range lines {
		newLine = append(newLine, []rune(l))
	}
	return
}

func doInternal(p *prefixCompleter, line string, pos int, origLine []rune) (newLine []string, offset int) {
	strings.TrimLeftFunc(line[:pos], func(r rune) bool {
		return unicode.IsSpace(r)
	})
	if p.splitChar != "" {
		line = splitKeepLast(line, p.splitChar)
	}
	for _, suggest := range p.callback(line) {
		line = strings.TrimLeft(line, "[")
		if len(line) >= len(suggest) {
			if strings.HasPrefix(line, suggest) {
				if len(line) != len(suggest) {
					newLine = append(newLine, suggest)
				}
				offset = len(suggest)
			}
		} else {
			if strings.HasPrefix(suggest, line) {
				newLine = append(newLine, suggest[len(line):])
				offset = len(line)
			}
		}
	}
	return
}

func splitKeepLast(s, sep string) (last string) {
	if !strings.Contains(s, sep) {
		last = s
		return
	}
	offset := strings.LastIndex(s, sep)
	if offset+1 < len(s) {
		last = s[offset+1:]
	}
	return
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

// Return potential resource types and a prop
// according to given holes questions.
// See corresponding unit test for logic
func guessEntityTypeFromHoleQuestion(hole string) (resolved []string, prop string) {
	tokens := strings.Split(strings.TrimSpace(hole), ".")
	if len(tokens) == 0 {
		return
	}

	var types []string
	for _, t := range tokens {
		for _, r := range resourcesTypesWithPlural {
			if t == r {
				types = append(types, cloud.SingularizeResource(r))
				break
			}
		}
	}

	if l := len(types); l > 0 {
		if len(tokens) == 2 {
			prop = tokens[1]
		}
		if l > 1 {
			prop = ""
		}
		resolved = []string{types[l-1]}
	} else if len(tokens) > 1 {
		for i := len(tokens) - 1; i >= 0; i-- {
			if len(tokens[i]) < 4 {
				continue
			}
			for _, r := range awsservices.ResourceTypes {
				if strings.Contains(r, tokens[i]) {
					resolved = append(resolved, r)
				}
			}
			if len(resolved) > 0 {
				return
			}
		}
	}
	return
}

func keyCorrespondsToProperty(holekey, prop string) bool {
	holekey = strings.ToLower(holekey)
	prop = strings.ToLower(prop)
	if holekey == prop {
		return true
	}
	if strings.Replace(holekey, "-", "", -1) == prop {
		return true
	}
	return false
}

func appendIfContains(slice []string, value, subst string) []string {
	subst = strings.TrimLeft(subst, "[")

	if strings.Contains(value, subst) && value != "" {
		return append(slice, value)
	}
	return slice
}

func appendWithNameAliases(slice []string, res *graph.Resource, s string) []string {
	if val, ok := res.Properties()["Name"]; ok {
		switch val.(type) {
		case string:
			name := val.(string)
			if name != "" {
				slice = appendIfContains(slice, fmt.Sprintf("@%s", name), s)
			}
		}
	}
	return slice
}

var resourcesTypesWithPlural []string

func init() {
	for _, r := range awsservices.ResourceTypes {
		resourcesTypesWithPlural = append(resourcesTypesWithPlural, r, cloud.PluralizeResource(r))
	}
}
