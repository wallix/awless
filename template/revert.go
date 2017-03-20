package template

import (
	"fmt"
	"strings"

	"github.com/wallix/awless/template/ast"
)

func (te *Template) Revert() (*Template, error) {
	var lines []string

	for _, cmd := range te.CmdNodesReverseIterator() {
		if isRevertible(cmd) {
			var revertAction string
			var params []string

			switch cmd.Action {
			case "create":
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
			case "start", "stop", "attach", "detach":
				for k, v := range cmd.Params {
					params = append(params, fmt.Sprintf("%s=%v", k, v))
				}
			case "create":
				switch cmd.Entity {
				case "record", "tag":
					for k, v := range cmd.Params {
						params = append(params, fmt.Sprintf("%s=%v", k, v))
					}
				case "database":
					params = append(params, fmt.Sprintf("id=%s", cmd.CmdResult))
					params = append(params, "skipsnapshot=true")
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
			}

			lines = append(lines, fmt.Sprintf("%s %s %s", revertAction, cmd.Entity, strings.Join(params, " ")))

			if cmd.Action == "create" && cmd.Entity == "instance" {
				lines = append(lines, fmt.Sprintf("check instance id=%s state=terminated timeout=180", cmd.CmdResult))
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
	if cmd.Entity == "record" && (cmd.Action == "create" || cmd.Action == "delete") {
		return true
	}

	if v, ok := cmd.CmdResult.(string); ok && v != "" {
		if cmd.Action == "create" || cmd.Action == "start" || cmd.Action == "stop" {
			return true
		}
	}

	return cmd.Action == "attach" || cmd.Action == "detach" || cmd.Action == "check" || (cmd.Action == "create" && cmd.Entity == "tag")
}
