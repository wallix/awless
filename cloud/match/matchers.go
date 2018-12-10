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
package match

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/wallix/awless/cloud"
)

type and struct {
	matchers []cloud.Matcher
}

func (m and) Match(r cloud.Resource) bool {
	for _, match := range m.matchers {
		if !match.Match(r) {
			return false
		}
	}
	return len(m.matchers) > 0
}

func And(matchers ...cloud.Matcher) cloud.Matcher {
	return and{matchers: matchers}
}

type or struct {
	matchers []cloud.Matcher
}

func (m or) Match(r cloud.Resource) bool {
	for _, match := range m.matchers {
		if match.Match(r) {
			return true
		}
	}
	return false
}

func Or(matchers ...cloud.Matcher) cloud.Matcher {
	return or{matchers: matchers}
}

type propertyMatcher struct {
	name          string
	value         interface{}
	matchOnString bool
	ignoreCase    bool
	contains      bool
}

func (m propertyMatcher) Match(r cloud.Resource) bool {
	v, found := r.Property(m.name)
	if !found {
		return false
	}
	expectVal := m.value
	if m.matchOnString {
		v = fmt.Sprint(v)
		expectVal = fmt.Sprint(m.value)
	}
	if m.ignoreCase {
		if vv, vIsStr := v.(string); vIsStr {
			v = strings.ToLower(vv)
		}
		if expect, expectIsStr := expectVal.(string); expectIsStr {
			expectVal = strings.ToLower(expect)
		}
	}
	if m.contains {
		vv, vIsStr := v.(string)
		expect, expectIsStr := expectVal.(string)
		if vIsStr && expectIsStr {
			return strings.Contains(vv, expect)
		}
	}
	return reflect.DeepEqual(v, expectVal)
}

func Property(name string, val interface{}) propertyMatcher {
	return propertyMatcher{name: name, value: val}
}

func (p propertyMatcher) MatchString() propertyMatcher {
	p.matchOnString = true
	return p
}

func (p propertyMatcher) IgnoreCase() propertyMatcher {
	p.ignoreCase = true
	return p
}

func (p propertyMatcher) Contains() propertyMatcher {
	p.contains = true
	return p
}

type tagMatcher struct {
	tagRegexp *regexp.Regexp
}

func (m tagMatcher) Match(r cloud.Resource) bool {
	tags, ok := r.Properties()["Tags"].([]string)
	if !ok {
		return false
	}
	for _, t := range tags {
		if m.tagRegexp.MatchString(t) {
			return true
		}
	}
	return false
}

func Tag(key, val string) tagMatcher {
	tagQuoteRegexp := "^" + regexp.QuoteMeta(key + "=" + val) + "$"
	tagWildcardRegexp := regexp.MustCompile(strings.Replace(tagQuoteRegexp, "\\*", ".*", -1))

	return tagMatcher{tagRegexp: tagWildcardRegexp}
}

type tagKeyMatcher struct {
	key string
}

func (m tagKeyMatcher) Match(r cloud.Resource) bool {
	tags, ok := r.Properties()["Tags"].([]string)
	if !ok {
		return false
	}
	for _, t := range tags {
		splits := strings.Split(t, "=")
		if len(splits) > 0 {
			if splits[0] == m.key {
				return true
			}
		}
	}
	return false
}

func TagKey(key string) tagKeyMatcher {
	return tagKeyMatcher{key: key}
}

type tagValueMatcher struct {
	value string
}

func (m tagValueMatcher) Match(r cloud.Resource) bool {
	tags, ok := r.Properties()["Tags"].([]string)
	if !ok {
		return false
	}
	for _, t := range tags {
		splits := strings.Split(t, "=")
		if len(splits) > 1 {
			if splits[1] == m.value {
				return true
			}
		}
	}
	return false
}

func TagValue(value string) tagValueMatcher {
	return tagValueMatcher{value: value}
}
