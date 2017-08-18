package main

import (
	"fmt"
	"hash/adler32"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/wallix/awless-scheduler/model"
	"github.com/wallix/awless/template"
	"github.com/wallix/awless/template/driver"
)

func New(filePath string) (tk *model.Task, err error) {
	tk = &model.Task{}

	var content []byte
	content, err = ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	tk.Content = string(content)
	fileName := filepath.Base(filePath)
	name := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	splits := strings.SplitN(name, "_", 4)
	checksum, err := strconv.ParseUint(splits[0], 10, 32)
	if err != nil {
		return
	}
	if cs := adler32.Checksum([]byte(tk.Content)); uint32(checksum) != cs {
		err = fmt.Errorf("unexpected checksum for file %s. Exepcted %d", name, cs)
		return
	}
	tk.RunAt, err = time.Parse(model.StampLayout, splits[1])
	if err != nil {
		return
	}
	tk.RevertAt, err = time.Parse(model.StampLayout, splits[2])
	if err != nil {
		return
	}
	tk.Region = splits[3]

	return
}

func executeTask(tk *model.Task, d driver.Driver, env *template.Env) (executed *template.Template, err error) {
	defer func() {
		id := tk.AsFilename()
		if err != nil {
			taskStore.MarkAsFailed(id)
		} else {
			err = taskStore.Remove(id)
		}
	}()

	var tpl, compiled, revertTmp *template.Template

	if tpl, err = template.Parse(tk.Content); err != nil {
		return
	}

	if compiled, _, err = template.Compile(tpl, env); err != nil {
		return
	}

	env.Driver = d

	if err = compiled.DryRun(env); err != nil {
		return
	}

	if executed, err = compiled.Run(env); err != nil {
		return
	}

	if executed.HasErrors() {
		var execErrors []string
		for _, cmd := range executed.CommandNodesIterator() {
			if cmd.CmdErr != nil {
				execErrors = append(execErrors, cmd.CmdErr.Error())
			}
		}
		err = fmt.Errorf(strings.Join(execErrors, ", "))
		return
	}

	if !tk.RevertAt.IsZero() && template.IsRevertible(executed) {
		if revertTmp, err = executed.Revert(); err != nil {
			return
		}
		revertTask := &model.Task{RunAt: tk.RevertAt, Region: tk.Region, Content: revertTmp.String()}
		if err = taskStore.Create(revertTask); err != nil {
			return
		}
	}
	return
}
