package main

import (
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/tomogoma/go-commons/errors"
)

const (
	microRepo = "github.com/micro/micro"
)

var (
	workingDir = ""
	buildFile  = ""

	unitTemplate = template.Must(template.New("unit").Parse(
		`[Unit]
Description=Micro {{.Command}} Auto Starter
After=consul.service
Requires=consul.service

[Install]
WantedBy=multi-user.target

[Service]
ExecStart=/usr/local/bin/micro {{.Command}}
SyslogIdentifier=micro {{- .Command}}
RestartSec=30
Restart=always`,
	))
)

func main() {

	if len(os.Args) < 2 {
		log.Fatalf("usage: installer [micro commands]\n" +
			"Where micro commands is a space separated list of " + microRepo + " commands.")
	}
	if err := initVars(); err != nil {
		log.Fatalf("Initializing: %v", err)
	}
	if err := fetchMicro(); err != nil {
		log.Fatalf("fetching micro: %v", err)
	}
	if err := buildMicro(); err != nil {
		log.Fatalf("building micro: %v", err)
	}
	for _, command := range os.Args[1:] {
		if err := buildUnit(command); err != nil {
			log.Fatalf("build unit file: %v", err)
		}
	}
}

func initVars() error {
	var err error
	workingDir, err = filepath.Abs("./")
	if err != nil {
		return errors.Newf("get working dir: %v", err)
	}
	buildFile = path.Join(workingDir, "bin", "micro")
	return err
}

func fetchMicro() error {
	//cmd := exec.Command("go", "get", "-u", microRepo)
	//if out, err := cmd.CombinedOutput(); err != nil {
	//	return errors.Newf("%v: %s", err, out)
	//}
	return nil
}

func buildMicro() error {

	repoDir := path.Join(os.Getenv("GOPATH"), "src", microRepo)

	cmd := exec.Command("git", "branch")
	cmd.Dir = repoDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Newf("listing micro's git release tags: %v: %s", err, out)
	}
	origiBranch := strings.TrimSpace(strings.TrimPrefix(string(out), "* "))

	cmd = exec.Command("git", "tag", "-l")
	cmd.Dir = repoDir
	out, err = cmd.CombinedOutput()
	if err != nil {
		return errors.Newf("listing micro's git release tags: %v: %s", err, out)
	}
	tags := strings.Split(string(out), "\n")
	if len(tags) == 0 {
		return errors.New("unexpected: no tags found in the micro repository")
	}
	latest := tags[0]
	for _, tag := range tags[1:] {
		if strings.Compare(latest, tag) > 0 {
			continue
		}
		latest = tag
	}

	cmd = exec.Command("git", "checkout", "tags/"+latest)
	cmd.Dir = repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return errors.Newf("checkout latest release: %v: %s", err, out)
	}

	cmd = exec.Command("go", "build", "-o", buildFile)
	cmd.Dir = repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return errors.Newf("%v: %s", err, out)
	}

	cmd = exec.Command("git", "checkout", origiBranch)
	cmd.Dir = repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return errors.Newf("%v: %s", err, out)
	}

	return nil
}

func buildUnit(forCommand string) error {

	if forCommand == "" {
		return errors.Newf("found empty micro command")
	}

	unitsDir := "unit"
	unitFile := path.Join(unitsDir, "micro"+forCommand+".service")

	if err := os.MkdirAll(unitsDir, 0755); err != nil {
		return errors.Newf("create unit files dir: %v", err)
	}

	f, err := os.Create(unitFile)
	if err != nil {
		return errors.Newf("create unit file: %v", err)
	}
	defer f.Close()

	data := struct {
		Command string
	}{Command: forCommand}
	if err := unitTemplate.Execute(f, data); err != nil {
		return errors.Newf("Executing template on '%s': %v", forCommand, err)
	}

	return nil
}
