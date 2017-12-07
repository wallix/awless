package awsspec

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"

	"github.com/wallix/awless/cloud/graph"
	"github.com/wallix/awless/template/params"

	"github.com/fatih/color"
	"github.com/wallix/awless/aws/doc"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/logger"
)

const (
	dryRunOperation = "DryRunOperation"
	notFound        = "NotFound"
)

type BeforeRunner interface {
	BeforeRun(ctx map[string]interface{}) error
}

type AfterRunner interface {
	AfterRun(ctx map[string]interface{}, output interface{}) error
}

type ResultExtractor interface {
	ExtractResult(interface{}) string
}

type command interface {
	ParamsHelp() string
	Params() params.Rule
	ValidateCommand(map[string]interface{}, []string) []error
	inject(params map[string]interface{}) error
	Run(ctx map[string]interface{}, params map[string]interface{}) (interface{}, error)
	DryRun(ctx, params map[string]interface{}) (interface{}, error)
}

func implementsBeforeRun(i interface{}) (BeforeRunner, bool) {
	v, ok := i.(BeforeRunner)
	return v, ok
}

func implementsAfterRun(i interface{}) (AfterRunner, bool) {
	v, ok := i.(AfterRunner)
	return v, ok
}

func implementsResultExtractor(i interface{}) (ResultExtractor, bool) {
	v, ok := i.(ResultExtractor)
	return v, ok
}

func validateParams(cmd command, params []string) ([]string, error) {
	paramsDefinitions := structListParamsKeys(cmd)
	var missing []string
	for n, isRequired := range paramsDefinitions {
		if isRequired && !contains(params, n) {
			missing = append(missing, n)
		}
	}

	var unexpected []string
	for _, p := range params {
		_, ok := paramsDefinitions[p]
		if !ok {
			unexpected = append(unexpected, p)
		}
	}

	switch len(unexpected) {
	case 0:
		return missing, nil
	case 1:
		return missing, fmt.Errorf("unexpected '%s' param\n%s", unexpected[0], cmd.ParamsHelp())
	default:
		return missing, fmt.Errorf("unexpected '%s' params\n%s", strings.Join(unexpected, "', '"), cmd.ParamsHelp())
	}

}

func generateParamsHelp(commandKey string, params map[string]bool) string {
	var buff bytes.Buffer
	var extra, required []string
	for n, isRequired := range params {
		if isRequired {
			required = append(required, n)
		} else {
			extra = append(extra, n)
		}
	}
	var anyRequired bool
	if len(required) > 0 {
		anyRequired = true
		buff.WriteString("\tRequired params:")
		for _, req := range required {
			buff.WriteString(fmt.Sprintf("\n\t\t- %s", req))
			if d, ok := awsdoc.TemplateParamsDoc(commandKey, req); ok {
				buff.WriteString(fmt.Sprintf(": %s", d))
			}
		}
	}

	if len(extra) > 0 {
		if anyRequired {
			buff.WriteString("\n\tExtra params:")
		} else {
			buff.WriteString("\n\tParams:")
		}
		for _, ext := range extra {
			buff.WriteString(fmt.Sprintf("\n\t\t- %s", ext))
			if d, ok := awsdoc.TemplateParamsDoc(commandKey, ext); ok {
				buff.WriteString(fmt.Sprintf(": %s", d))
			}
		}
	}
	return buff.String()
}

func fakeDryRunId(entity string) string {
	suffix := rand.Intn(1e6)
	switch entity {
	case cloud.Instance:
		return fmt.Sprintf("i-%d", suffix)
	case cloud.Subnet:
		return fmt.Sprintf("subnet-%d", suffix)
	case cloud.Vpc:
		return fmt.Sprintf("vpc-%d", suffix)
	case cloud.Volume:
		return fmt.Sprintf("vol-%d", suffix)
	case cloud.SecurityGroup:
		return fmt.Sprintf("sg-%d", suffix)
	case cloud.InternetGateway:
		return fmt.Sprintf("igw-%d", suffix)
	case cloud.NatGateway:
		return fmt.Sprintf("nat-%d", suffix)
	case cloud.RouteTable:
		return fmt.Sprintf("rtb-%d", suffix)
	default:
		return fmt.Sprintf("dryrunid-%d", suffix)
	}
}

type paramRule struct {
	tree   ruleNode
	extras []string
}

func (p paramRule) help() string {
	if len(p.extras) == 0 {
		return p.tree.help()
	}

	return fmt.Sprintf("%s or extra params: \"%s\"", p.tree.help(), strings.Join(p.extras, "\", \""))
}

func (p paramRule) verify(keys []string) ([]string, error) {
	if p.tree == nil {
		return nil, nil
	}
	var unexpected []string
	for _, key := range keys {
		if p.tree.unexpected(key) && !contains(p.extras, key) {
			unexpected = append(unexpected, key)
		}
	}
	switch len(unexpected) {
	case 0:
		break
	case 1:
		return nil, fmt.Errorf("unexpected '%s' param\n%s", unexpected[0], p.help())
	default:
		return nil, fmt.Errorf("unexpected '%s' params\n%s", strings.Join(unexpected, "', '"), p.help())
	}

	missings, _, errs := p.tree.missings(keys)
	if len(errs) > 0 {
		var errStr bytes.Buffer
		for _, e := range errs {
			errStr.WriteString(e.Error())
		}
		return nil, errors.New(errStr.String())
	}
	return missings, nil
}

type ruleNode interface {
	help() string
	unexpected(string) bool
	missings([]string) ([]string, int, []error)
}

type oneOfNode struct {
	children []ruleNode
}

func (o oneOfNode) help() string {
	var childrenHint []string
	for _, c := range o.children {
		childrenHint = append(childrenHint, c.help())
	}
	return fmt.Sprintf("(%s)", strings.Join(childrenHint, " or "))
}

func (o oneOfNode) unexpected(s string) bool {
	for _, c := range o.children {
		if !c.unexpected(s) {
			return false
		}
	}
	return true
}

func (o oneOfNode) missings(keys []string) ([]string, int, []error) {
	var errs []error
	maxFound := -1
	var missings []string
	for _, child := range o.children {
		cMissings, nbFound, err := child.missings(keys)
		errs = append(errs, err...)
		if nbFound > maxFound {
			missings = cMissings
			maxFound = nbFound
		}
	}
	return missings, maxFound, nil
}

type oneOfNodeWithError struct {
	oneOfNode
}

func (o oneOfNodeWithError) missings(keys []string) (missings []string, found int, errs []error) {
	var hasFoundChild bool
	for _, child := range o.children {
		_, nbFound, _ := child.missings(keys)
		if nbFound > 0 {
			hasFoundChild = true
		}
	}
	missings, found, errs = o.oneOfNode.missings(keys)
	if !hasFoundChild {
		errs = append(errs, fmt.Errorf("expecting %s", o.help()))
	}
	return
}

type allOfNode struct {
	children []ruleNode
}

func (a allOfNode) help() string {
	var childrenHint []string
	for _, c := range a.children {
		childrenHint = append(childrenHint, c.help())
	}
	return fmt.Sprintf("(%s)", strings.Join(childrenHint, " and "))
}

func (a allOfNode) unexpected(s string) bool {
	for _, c := range a.children {
		if !c.unexpected(s) {
			return false
		}
	}
	return true
}

func (a allOfNode) missings(keys []string) (missings []string, nbFound int, errs []error) {
	for _, child := range a.children {
		cMissings, cFound, err := child.missings(keys)
		errs = append(errs, err...)
		if len(cMissings) > 0 {
			missings = append(missings, cMissings...)
		} else {
			nbFound += cFound
		}
	}
	return
}

type stringNode struct {
	key string
}

func (k stringNode) help() string {
	return fmt.Sprintf("\"%s\"", k.key)
}

func (k stringNode) unexpected(s string) bool {
	return k.key != s
}

func (k stringNode) missings(keys []string) (missings []string, nbFound int, errs []error) {
	if contains(keys, k.key) {
		nbFound++
		return
	}
	missings = append(missings, k.key)
	return
}

func oneOf(nodes ...ruleNode) ruleNode {
	return oneOfNode{children: nodes}
}

func oneOfE(nodes ...ruleNode) ruleNode {
	return oneOfNodeWithError{oneOfNode{children: nodes}}
}

func allOf(nodes ...ruleNode) ruleNode {
	return allOfNode{children: nodes}
}

func node(key string) ruleNode {
	return stringNode{key}
}

type awsCall struct {
	fnName  string
	fn      interface{}
	logger  *logger.Logger
	graph   cloudgraph.GraphAPI
	setters []setter
}

func (dc *awsCall) execute(input interface{}) (output interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			output = nil
			err = fmt.Errorf("%s", e)
		}
	}()

	for _, s := range dc.setters {
		if err = s.set(input); err != nil {
			return nil, err
		}
	}

	fnVal := reflect.ValueOf(dc.fn)
	values := []reflect.Value{reflect.ValueOf(input)}

	start := time.Now()
	results := fnVal.Call(values)

	if err, ok := results[1].Interface().(error); ok && err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	dc.logger.ExtraVerbosef("%s call took %s", dc.fnName, time.Since(start))

	output = results[0].Interface()

	return
}

type checker struct {
	description string
	timeout     time.Duration
	frequency   time.Duration
	fetchFunc   func() (string, error)
	expect      string
	logger      *logger.Logger
	graph       cloudgraph.GraphAPI
	checkName   string
}

func (c *checker) check() error {
	now := time.Now().UTC()
	timer := time.NewTimer(c.timeout)
	if c.checkName == "" {
		c.checkName = "status"
	}
	defer timer.Stop()
	defer c.logger.Println()
	for {
		select {
		case <-timer.C:
			return fmt.Errorf("timeout of %s expired", c.timeout)
		default:
		}
		got, err := c.fetchFunc()
		if err != nil {
			return fmt.Errorf("check %s: %s", c.description, err)
		}
		if strings.ToLower(got) == strings.ToLower(c.expect) {
			c.logger.InteractiveInfof("check %s %s '%s' done", c.description, c.checkName, c.expect)
			return nil
		}
		elapsed := time.Since(now)
		c.logger.InteractiveInfof("%s %s '%s', expect '%s', timeout in %s (retry in %s)", c.description, c.checkName, got, c.expect, color.New(color.FgGreen).Sprint(c.timeout-elapsed.Round(time.Second)), c.frequency)
		time.Sleep(c.frequency)
	}
}

type enumValidator struct {
	expected []string
}

func NewEnumValidator(expected ...string) *enumValidator {
	return &enumValidator{expected: expected}
}

func (v *enumValidator) Validate(in *string) error {
	val := strings.ToLower(StringValue(in))
	for _, e := range v.expected {
		if val == strings.ToLower(e) {
			return nil
		}
	}
	var expString string
	switch len(v.expected) {
	case 0:
		return errors.New("empty enumeration")
	case 1:
		expString = fmt.Sprintf("'%s'", v.expected[0])
	case 2:
		expString = fmt.Sprintf("'%s' or '%s'", v.expected[0], v.expected[1])
	default:
		expString = fmt.Sprintf("'%s' or '%s'", strings.Join(v.expected[0:len(v.expected)-1], "', '"), v.expected[len(v.expected)-1])
	}
	return fmt.Errorf("invalid value '%s' expect %s", StringValue(in), expString)
}

func String(v string) *string {
	return &v
}

func StringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

func Int64(v int64) *int64 {
	return &v
}

func Int64AsIntValue(v *int64) int {
	if v != nil {
		return int(*v)
	}
	return 0
}

func Bool(v bool) *bool {
	return &v
}

func BoolValue(v *bool) bool {
	if v != nil {
		return *v
	}
	return false
}
