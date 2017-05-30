package template

import (
	"fmt"
	"strings"

	"github.com/wallix/awless/template/internal/ast"
)

func (te *Template) Revert() (*Template, error) {
	var lines []string
	cmdsReverseIterator := te.CmdNodesReverseIterator()
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
			}

			switch cmd.Action {
			case "attach":
				switch cmd.Entity {
				case "routetable", "elasticip":
					params = append(params, fmt.Sprintf("association=%s", cmd.CmdResult))
				case "instance":
					for k, v := range cmd.Params {
						if k == "port" {
							continue
						}
						params = append(params, fmt.Sprintf("%s=%v", k, v))
					}
				default:
					for k, v := range cmd.Params {
						params = append(params, fmt.Sprintf("%s=%v", k, v))
					}
				}
			case "start", "stop", "detach":
				switch {
				case cmd.Entity == "routetable":
					params = append(params, fmt.Sprintf("association=%s", cmd.CmdResult))
				case cmd.Entity == "volume" && cmd.Action == "detach":
					for k, v := range cmd.Params {
						if k == "force" {
							continue
						}
						params = append(params, fmt.Sprintf("%s=%v", k, v))
					}
				default:
					for k, v := range cmd.Params {
						params = append(params, fmt.Sprintf("%s=%v", k, v))
					}
				}
			case "create":
				switch cmd.Entity {
				case "record", "tag":
					for k, v := range cmd.Params {
						params = append(params, fmt.Sprintf("%s=%v", k, v))
					}
				case "route":
					for k, v := range cmd.Params {
						if k == "gateway" {
							continue
						}
						params = append(params, fmt.Sprintf("%s=%v", k, v))
					}
				case "database":
					params = append(params, fmt.Sprintf("id=%s", cmd.CmdResult))
					params = append(params, "skipsnapshot=true")
				case "policy":
					params = append(params, fmt.Sprintf("arn=%s", cmd.CmdResult))
				case "queue":
					params = append(params, fmt.Sprintf("url=%s", cmd.CmdResult))
				case "s3object":
					params = append(params, fmt.Sprintf("name=%s", cmd.CmdResult))
					params = append(params, fmt.Sprintf("bucket=%s", cmd.Params["bucket"]))
				case "role", "group", "user", "stack":
					params = append(params, fmt.Sprintf("name=%s", cmd.Params["name"]))
				case "accesskey":
					params = append(params, fmt.Sprintf("id=%s", cmd.CmdResult))
					params = append(params, fmt.Sprintf("user=%s", cmd.Params["user"]))
				case "loginprofile":
					params = append(params, fmt.Sprintf("username=%s", cmd.Params["username"]))
				case "bucket", "launchconfiguration", "scalinggroup", "alarm", "dbsubnetgroup", "keypair":
					params = append(params, fmt.Sprintf("name=%s", cmd.CmdResult))
					if cmd.Entity == "scalinggroup" {
						params = append(params, "force=true")
					}
				default:
					params = append(params, fmt.Sprintf("id=%s", cmd.CmdResult))
				}
			case "delete":
				switch cmd.Entity {
				case "record":
					for k, v := range cmd.Params {
						params = append(params, fmt.Sprintf("%s=%v", k, v))
					}
				}
			case "copy":
				switch cmd.Entity {
				case "image":
					params = append(params, fmt.Sprintf("id=%s", cmd.CmdResult))
					params = append(params, "delete-snapshots=true")
				default:
					params = append(params, fmt.Sprintf("id=%s", cmd.CmdResult))
				}
			}

			// Prechecks
			if cmd.Action == "create" && cmd.Entity == "securitygroup" {
				lines = append(lines, fmt.Sprintf("check securitygroup id=%s state=unused timeout=180", cmd.CmdResult))
			}
			if cmd.Action == "create" && cmd.Entity == "scalinggroup" {
				lines = append(lines, fmt.Sprintf("update scalinggroup name=%s max-size=0 min-size=0", cmd.CmdResult))
				lines = append(lines, fmt.Sprintf("check scalinggroup count=0 name=%s timeout=180", cmd.CmdResult))
			}
			if cmd.Action == "start" && cmd.Entity == "instance" {
				lines = append(lines, fmt.Sprintf("check instance id=%s state=running timeout=180", cmd.Params["id"]))
			}
			if cmd.Action == "stop" && cmd.Entity == "instance" {
				lines = append(lines, fmt.Sprintf("check instance id=%s state=stopped timeout=180", cmd.Params["id"]))
			}

			lines = append(lines, fmt.Sprintf("%s %s %s", revertAction, cmd.Entity, strings.Join(params, " ")))

			// Postchecks
			if notLastCommand {
				if cmd.Action == "create" && cmd.Entity == "instance" {
					lines = append(lines, fmt.Sprintf("check instance id=%s state=terminated timeout=180", cmd.CmdResult))
				}
				if cmd.Action == "create" && cmd.Entity == "loadbalancer" {
					lines = append(lines, fmt.Sprintf("check loadbalancer id=%s state=not-found timeout=180", cmd.CmdResult))
				}
				if cmd.Action == "attach" && cmd.Entity == "volume" {
					lines = append(lines, fmt.Sprintf("check volume id=%s state=available timeout=180", cmd.Params["id"]))
				}
			}
		}
	}

	text := strings.Join(lines, "\n")
	tpl, err := Parse(text)
	if err != nil {
		return nil, fmt.Errorf("revert: \n%s\n%s", text, err)
	}

	return tpl, nil
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

	if cmd.Entity == "alarm" && (cmd.Action == "start" || cmd.Action == "stop") {
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
