package ansibleutils

import (
	"bytes"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/logger"
	"io"
	"runtime"

	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
)

type Ansible struct {
	Playbook string
}

type AnsibleParams struct {
	Playbooks []string `json:"playbooks"`
	Inventory string   `json:"inventory"`
	Roles     string   `json:"roles"`
	Packages  string   `json:"packages"`
	Base      string   `json:"base"`
}

func New() *Ansible {
	return &Ansible{}
}

func (a *Ansible) InstallPackages(packages string) error {
	if packages != "" {
		out, _, err := command([]string{"sh", "-c", "apt update && apt install -y " + packages})
		if err != nil {
			return fmt.Errorf("apt install error %w (%s)", err, out)
		}
	}
	return nil
}

func (a *Ansible) InstallRoles(roles string) error {
	if roles != "" {
		out, _, err := command([]string{"sh", "-c", "ansible-galaxy install " + roles})
		if err != nil {
			return fmt.Errorf("ansible-galaxy install error %w (%s)", err, out)
		}
	}
	return nil
}

// RunPlaybook runs one playbook
func (a *Ansible) RunPlaybook(ctx context.Context, base, playbookUrl, inventory string) (string, error) {
	var err error

	log := logger.FromContext(ctx, true)

	buff := new(bytes.Buffer)

	httpFiles := NewHttpFiles("/tmp")
	log.WithField("url", base+playbookUrl).Trace("get playbook file")
	locakPlaybookFile, err := httpFiles.GetFile(base + playbookUrl)
	if err != nil {
		return "", fmt.Errorf("get playbook file : %w", err)
	}

	log.WithField("url", base+inventory).Trace("get inventory file")
	localInventoryFile, err := httpFiles.GetFile(base + inventory)
	if err != nil {
		return "", fmt.Errorf("get inventory file : %w", err)
	}

	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		Connection: "local",
		//		User:       "root",
	}

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: localInventoryFile,
	}

	execute := execute.NewDefaultExecute(
		execute.WithWrite(io.Writer(buff)),
	)

	playbook := &playbook.AnsiblePlaybookCmd{
		Playbooks:         []string{locakPlaybookFile},
		Exec:              execute,
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		Options:           ansiblePlaybookOptions,
		StdoutCallback:    "json",
	}
	// TODO : search ansible-playbook on OSX
	if runtime.GOOS == "darwin" {
		playbook.Binary = "/Users/telemac/Library/Python/3.9/bin/ansible-playbook"
	}

	log.WithField("playbook", locakPlaybookFile).Trace("run playbook")

	err = playbook.Run(ctx)
	if err != nil {
		//log.WithError(err).Warn("playbook run")
		return "", err
	}

	log.WithFields(logrus.Fields{
		"result":   buff.String(),
		"playbook": locakPlaybookFile,
		"error":    err,
	}).Debug("playbook result")

	return buff.String(), nil
}

func (a *Ansible) RunPlaybooks(ctx context.Context, params AnsibleParams) ([]string, error) {
	log := logger.FromContext(ctx, true)
	// Install specified packages
	log.Debug("install packages")
	err := a.InstallPackages(params.Packages)
	if err != nil {
		return nil, err
	}
	// install specified roles
	log.Debug("install roles")
	err = a.InstallRoles(params.Roles)
	if err != nil {
		return nil, err
	}

	var _results []string
	for _, pb := range params.Playbooks {
		res, err := a.RunPlaybook(ctx, params.Base, pb, params.Inventory)
		_results = append(_results, res)
		if err != nil {
			return _results, err
		}
	}
	return _results, nil
}
