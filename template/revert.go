package template

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/wallix/awless/template/internal/ast"
)

func (temp *Template) Revert() (*Template, error) {
	tpl, _, err := Compile(temp, new(noopCompileEnv), PreRevertCompileMode)
	if err != nil {
		return temp, err
	}

	var lines []string
	cmdsReverseIterator := tpl.CommandNodesReverseIterator()
	for i, cmd := range cmdsReverseIterator {
		notLastCommand := (i != len(cmdsReverseIterator)-1)
		if isRevertible(cmd) {
			var revertAction string
			var params []string

			switch cmd.Action {
			case "create", "copy":
				revertAction = "delete"
			case "start":
				revertAction = "stop"
			case "stop":
				revertAction = "start"
			case "detach":
				revertAction = "attach"
			case "attach":
				revertAction = "detach"
			case "delete":
				revertAction = "create"
			case "update":
				revertAction = "update"
			}

			switch cmd.Action {
			case "attach":
				switch cmd.Entity {
				case "routetable", "elasticip":
					params = append(params, fmt.Sprintf("association=%s", quoteParamIfNeeded(cmd.CmdResult)))
				case "instance":
					for k, v := range cmd.ParamNodes {
						if k == "port" {
							continue
						}
						params = append(params, fmt.Sprintf("%s=%s", k, printItem(v)))
					}
				case "containertask":
					params = append(params, fmt.Sprintf("name=%s", printItem(cmd.ParamNodes["name"])))
					params = append(params, fmt.Sprintf("container-name=%s", printItem(cmd.ParamNodes["container-name"])))
				case "networkinterface":
					params = append(params, fmt.Sprintf("attachment=%s", quoteParamIfNeeded(cmd.CmdResult)))
				case "mfadevice":
					params = append(params, fmt.Sprintf("id=%s", printItem(cmd.ParamNodes["id"])))
					params = append(params, fmt.Sprintf("user=%s", printItem(cmd.ParamNodes["user"])))
				default:
					for k, v := range cmd.ParamNodes {
						params = append(params, fmt.Sprintf("%s=%v", k, v))
					}
				}
			case "start", "stop", "detach":
				switch {
				case cmd.Entity == "routetable":
					params = append(params, fmt.Sprintf("association=%s", quoteParamIfNeeded(cmd.CmdResult)))
				case cmd.Entity == "volume" && cmd.Action == "detach":
					for k, v := range cmd.ParamNodes {
						if k == "force" {
							continue
						}
						params = append(params, fmt.Sprintf("%s=%v", k, printItem(v)))
					}
				case cmd.Entity == "containertask":
					params = append(params, fmt.Sprintf("cluster=%s", printItem(cmd.ParamNodes["cluster"])))
					params = append(params, fmt.Sprintf("type=%s", printItem(cmd.ParamNodes["type"])))
					switch fmt.Sprint(printItem(cmd.ParamNodes["type"])) {
					case "service":
						params = append(params, fmt.Sprintf("deployment-name=%s", printItem(cmd.ParamNodes["deployment-name"])))
					case "task":
						params = append(params, fmt.Sprintf("run-arn=%s", quoteParamIfNeeded(cmd.CmdResult)))
					default:
						return nil, fmt.Errorf("start containertask with type '%v' can not be reverted", printItem(cmd.ParamNodes["deployment-name"]))
					}
				default:
					for k, v := range cmd.ParamNodes {
						params = append(params, fmt.Sprintf("%s=%v", k, printItem(v)))
					}
				}
			case "create":
				switch cmd.Entity {
				case "tag":
					for k, v := range cmd.ParamNodes {
						params = append(params, fmt.Sprintf("%s=%v", k, printItem(v)))
					}
				case "record":
					for k, v := range cmd.ParamNodes {
						if k == "comment" {
							continue
						}
						params = append(params, fmt.Sprintf("%s=%v", k, printItem(v)))
					}
				case "route":
					for k, v := range cmd.ParamNodes {
						if k == "gateway" {
							continue
						}
						params = append(params, fmt.Sprintf("%s=%v", k, printItem(v)))
					}
				case "database":
					params = append(params, fmt.Sprintf("id=%s", quoteParamIfNeeded(cmd.CmdResult)))
					params = append(params, "skip-snapshot=true")
				case "certificate":
					params = append(params, fmt.Sprintf("arn=%s", quoteParamIfNeeded(cmd.CmdResult)))
				case "policy":
					params = append(params, fmt.Sprintf("arn=%s", quoteParamIfNeeded(cmd.CmdResult)))
					params = append(params, "all-versions=true")
				case "queue":
					params = append(params, fmt.Sprintf("url=%s", quoteParamIfNeeded(cmd.CmdResult)))
				case "s3object":
					params = append(params, fmt.Sprintf("name=%s", quoteParamIfNeeded(cmd.CmdResult)))
					params = append(params, fmt.Sprintf("bucket=%s", printItem(cmd.ParamNodes["bucket"])))
				case "role", "group", "user", "stack", "instanceprofile", "repository", "classicloadbalancer":
					params = append(params, fmt.Sprintf("name=%s", printItem(cmd.ParamNodes["name"])))
				case "accesskey":
					params = append(params, fmt.Sprintf("id=%s", quoteParamIfNeeded(cmd.CmdResult)))
					params = append(params, fmt.Sprintf("user=%s", printItem(cmd.ParamNodes["user"])))
				case "appscalingtarget":
					params = append(params, fmt.Sprintf("dimension=%s", printItem(cmd.ParamNodes["dimension"])))
					params = append(params, fmt.Sprintf("resource=%s", printItem(cmd.ParamNodes["resource"])))
					params = append(params, fmt.Sprintf("service-namespace=%s", printItem(cmd.ParamNodes["service-namespace"])))
				case "appscalingpolicy":
					params = append(params, fmt.Sprintf("dimension=%s", printItem(cmd.ParamNodes["dimension"])))
					params = append(params, fmt.Sprintf("name=%s", printItem(cmd.ParamNodes["name"])))
					params = append(params, fmt.Sprintf("resource=%s", printItem(cmd.ParamNodes["resource"])))
					params = append(params, fmt.Sprintf("service-namespace=%s", printItem(cmd.ParamNodes["service-namespace"])))
				case "loginprofile":
					params = append(params, fmt.Sprintf("username=%s", printItem(cmd.ParamNodes["username"])))
				case "bucket", "launchconfiguration", "scalinggroup", "alarm", "dbsubnetgroup", "keypair":
					params = append(params, fmt.Sprintf("name=%s", quoteParamIfNeeded(cmd.CmdResult)))
					if cmd.Entity == "scalinggroup" {
						params = append(params, "force=true")
					}
				default:
					params = append(params, fmt.Sprintf("id=%s", quoteParamIfNeeded(cmd.CmdResult)))
				}
			case "delete":
				switch cmd.Entity {
				case "record":
					for k, v := range cmd.ParamNodes {
						params = append(params, fmt.Sprintf("%s=%v", k, quoteParamIfNeeded(v)))
					}
				case "instanceprofile":
					params = append(params, fmt.Sprintf("name=%s", printItem(cmd.ParamNodes["name"])))
				}
			case "copy":
				switch cmd.Entity {
				case "image":
					params = append(params, fmt.Sprintf("id=%s", quoteParamIfNeeded(cmd.CmdResult)))
					params = append(params, "delete-snapshots=true")
				default:
					params = append(params, fmt.Sprintf("id=%s", quoteParamIfNeeded(cmd.CmdResult)))
				}
			case "update":
				switch cmd.Entity {
				case "securitygroup":
					for k, v := range cmd.ParamNodes {
						if k == "inbound" || k == "outbound" {
							if fmt.Sprint(v) == "authorize" {
								params = append(params, fmt.Sprintf("%s=revoke", k))
							} else if fmt.Sprint(v) == "revoke" {
								params = append(params, fmt.Sprintf("%s=authorize", k))
							}
							continue
						}
						params = append(params, fmt.Sprintf("%s=%v", k, printItem(v)))
					}
				}
			}

			// Prechecks
			if cmd.Action == "create" && cmd.Entity == "securitygroup" {
				lines = append(lines, fmt.Sprintf("check securitygroup id=%s state=unused timeout=300", quoteParamIfNeeded(cmd.CmdResult)))
			}
			if cmd.Action == "create" && cmd.Entity == "scalinggroup" {
				lines = append(lines, fmt.Sprintf("update scalinggroup name=%s max-size=0 min-size=0", quoteParamIfNeeded(cmd.CmdResult)))
				lines = append(lines, fmt.Sprintf("check scalinggroup count=0 name=%s timeout=600", quoteParamIfNeeded(cmd.CmdResult)))
			}
			if cmd.Action == "start" && cmd.Entity == "instance" {
				switch vv := cmd.ParamNodes["ids"].(type) {
				case string:
					lines = append(lines, fmt.Sprintf("check instance id=%s state=running timeout=180", printItem(vv)))
				case []interface{}:
					for _, s := range vv {
						lines = append(lines, fmt.Sprintf("check instance id=%v state=running timeout=180", printItem(s)))
					}
				default:
					return nil, fmt.Errorf("revert start instance: unexpected type of ids: %T", vv)
				}
			}
			if cmd.Action == "stop" && cmd.Entity == "instance" {
				switch vv := cmd.ParamNodes["ids"].(type) {
				case string:
					lines = append(lines, fmt.Sprintf("check instance id=%s state=stopped timeout=180", printItem(vv)))
				case []interface{}:
					for _, s := range vv {
						lines = append(lines, fmt.Sprintf("check instance id=%v state=stopped timeout=180", printItem(s)))
					}
				default:
					return nil, fmt.Errorf("revert stop instance: unexpected type of ids: %T", vv)
				}
			}
			if cmd.Action == "start" && cmd.Entity == "containertask" && fmt.Sprint(printItem(cmd.ParamNodes["type"])) == "service" {
				lines = append(lines, fmt.Sprintf("update containertask cluster=%s deployment-name=%s desired-count=0", printItem(cmd.ParamNodes["cluster"]), printItem(cmd.ParamNodes["deployment-name"])))
			}

			lines = append(lines, fmt.Sprintf("%s %s %s", revertAction, cmd.Entity, strings.Join(params, " ")))

			// Postchecks
			if notLastCommand {
				if cmd.Action == "create" && cmd.Entity == "instance" {
					lines = append(lines, fmt.Sprintf("check instance id=%s state=terminated timeout=180", quoteParamIfNeeded(cmd.CmdResult)))
				}
				if cmd.Action == "create" && cmd.Entity == "database" {
					lines = append(lines, fmt.Sprintf("check database id=%s state=not-found timeout=900", quoteParamIfNeeded(cmd.CmdResult)))
				}
				if cmd.Action == "create" && cmd.Entity == "loadbalancer" {
					lines = append(lines, fmt.Sprintf("check loadbalancer id=%s state=not-found timeout=180", quoteParamIfNeeded(cmd.CmdResult)))
				}
				if cmd.Action == "attach" && cmd.Entity == "volume" {
					lines = append(lines, fmt.Sprintf("check volume id=%s state=available timeout=180", printItem(cmd.ParamNodes["id"])))
				}
				if cmd.Action == "create" && cmd.Entity == "natgateway" {
					lines = append(lines, fmt.Sprintf("check natgateway id=%s state=deleted timeout=180", quoteParamIfNeeded(cmd.CmdResult)))
				}
			}
		}
	}

	text := strings.Join(lines, "\n")
	reverted, err := Parse(text)
	if err != nil {
		return nil, fmt.Errorf("revert: \n%s\n%s", text, err)
	}

	return reverted, nil
}

func IsRevertible(t *Template) bool {
	revertible := false
	t.visitCommandNodes(func(cmd *ast.CommandNode) {
		if isRevertible(cmd) {
			revertible = true
		}
	})
	return revertible
}

func isRevertible(cmd *ast.CommandNode) bool {
	if cmd.CmdErr != nil {
		return false
	}

	if cmd.Action == "check" {
		return false
	}

	if cmd.Action == "detach" && cmd.Entity == "routetable" {
		return false
	}

	if cmd.Entity == "record" && (cmd.Action == "create" || cmd.Action == "delete") {
		return true
	}

	if cmd.Entity == "instanceprofile" && (cmd.Action == "create" || cmd.Action == "delete") {
		return true
	}

	if cmd.Entity == "alarm" && (cmd.Action == "start" || cmd.Action == "stop") {
		return true
	}

	if cmd.Entity == "database" && (cmd.Action == "start" || cmd.Action == "stop") {
		return true
	}

	if cmd.Entity == "containertask" && cmd.Action == "start" {
		t, ok := cmd.ToDriverParams()["type"].(string)
		return ok && (t == "service" || t == "task")
	}

	if cmd.Entity == "container" && cmd.Action == "create" {
		return true
	}

	if cmd.Entity == "appscalingtarget" && cmd.Action == "create" {
		return true
	}

	if cmd.Entity == "securitygroup" && cmd.Action == "update" {
		return true
	}

	if cmd.Entity == "appscalingpolicy" && cmd.Action == "create" {
		return true
	}

	if v, ok := cmd.CmdResult.(string); ok && v != "" {
		if cmd.Action == "create" || cmd.Action == "start" || cmd.Action == "stop" || cmd.Action == "copy" {
			return true
		}
	}

	return cmd.Action == "attach" || cmd.Action == "detach" || cmd.Action == "check" ||
		(cmd.Action == "create" && cmd.Entity == "tag") || (cmd.Action == "create" && cmd.Entity == "route")
}

func printItem(i interface{}) string {
	switch ii := i.(type) {
	case string:
		if _, err := strconv.Atoi(ii); err == nil {
			return "'" + ii + "'"
		}
		if _, err := strconv.ParseFloat(ii, 64); err == nil {
			return "'" + ii + "'"
		}
		return quoteParamIfNeeded(i)
	case []interface{}:
		var out []string
		for _, e := range ii {
			out = append(out, printItem(e))
		}
		return "[" + strings.Join(out, ",") + "]"
	default:
		return fmt.Sprint(ii)
	}
}

func quoteParamIfNeeded(param interface{}) string {
	input := fmt.Sprint(param)
	if ast.SimpleStringValue.MatchString(input) {
		return input
	} else {
		if strings.ContainsRune(input, '\'') {
			return "\"" + input + "\""
		} else {
			return "'" + input + "'"
		}
	}
}
