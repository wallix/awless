package awsspec

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/wallix/awless/template/env"
	"github.com/wallix/awless/template/params"

	"github.com/fatih/color"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/logger"
)

const (
	dryRunOperation = "DryRunOperation"
	notFound        = "NotFound"
)

type BeforeRunner interface {
	BeforeRun(env.Running) error
}

type AfterRunner interface {
	AfterRun(env.Running, interface{}) error
}

type ResultExtractor interface {
	ExtractResult(interface{}) string
}

type command interface {
	ParamsSpec() params.Spec
	inject(map[string]interface{}) error
	Run(env.Running, map[string]interface{}) (interface{}, error)
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

type awsCall struct {
	fnName  string
	fn      interface{}
	logger  *logger.Logger
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

func decorateAWSError(err error) error {
	if aerr, ok := err.(awserr.Error); ok {
		return fmt.Errorf("%s: %s", aerr.Code(), aerr.Message())
	}
	return err
}
