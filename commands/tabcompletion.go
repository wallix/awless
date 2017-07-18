package commands

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/chzyer/readline"
	"github.com/wallix/awless/aws/services"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/template"
)

func holeAutoCompletion(g *graph.Graph, hole string) readline.AutoCompleter {
	completeFunc := func(string) []string { return []string{} }

	if entityType, entityProp := guessEntityTypeFromHoleQuestion(hole); entityType != "" {
		resources, err := g.GetAllResources(entityType)
		exitOn(err)

		if len(resources) == 0 {
			return &prefixCompleter{callback: completeFunc}
		}

		if entityProp != "" {
			var validPropName string
			for _, r := range resources {
				for propName := range r.Properties {
					if keyCorrespondsToProperty(entityProp, propName) {
						validPropName = propName
					}
				}
			}

			if validPropName != "" {
				completeFunc = func(s string) (suggest []string) {
					for _, res := range resources {
						if v, ok := res.Properties[validPropName]; ok {
							switch prop := v.(type) {
							case string, float64, int, bool:
								suggest = appendIfContains(suggest, fmt.Sprint(prop), s)
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
			}
			return &prefixCompleter{callback: completeFunc, splitChar: ","}
		}

		completeFunc = func(s string) (suggest []string) {
			s = splitKeepLast(s, ",")
			s = strings.TrimLeft(s, "'@\"")
			for _, res := range resources {
				suggest = appendIfContains(suggest, res.Id(), s)
				if val, ok := res.Properties["Name"]; ok {
					switch val.(type) {
					case string:
						name := val.(string)
						if name != "" {
							suggest = appendIfContains(suggest, fmt.Sprintf("@%s", name), s)
						}
					}
				}
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
		last = s[offset+1 : len(s)]
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

// Return 2 strings according to holes questions:
// notmatching? -> ["", ""]
// entitytype? -> [entitytype, ""]
// entitytype.anything? -> [entitytype, anything]
// entitytype.otherentitytype? -> [otherentitytype, ""]
func guessEntityTypeFromHoleQuestion(hole string) (string, string) {
	var types []string
	splits := strings.Split(hole, ".")
	for _, t := range splits {
		for _, r := range resourcesTypesWithPlural {
			if t == r {
				types = append(types, cloud.SingularizeResource(r))
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
	if strings.Contains(value, subst) && value != "" {
		return append(slice, value)
	}
	return slice
}

var resourcesTypesWithPlural []string

func init() {
	for _, r := range awsservices.ResourceTypes {
		resourcesTypesWithPlural = append(resourcesTypesWithPlural, r, cloud.PluralizeResource(r))
	}
}
