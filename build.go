package main

import (
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/tomogoma/go-typed-errors"
	"sort"
	"strconv"
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
	cmd := exec.Command("go", "get", "-u", microRepo)
	if out, err := cmd.CombinedOutput(); err != nil {
		return errors.Newf("%v: %s", err, out)
	}
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
	defer func() {
		cmd = exec.Command("git", "checkout", origiBranch)
		cmd.Dir = repoDir
		if out, err := cmd.CombinedOutput(); err != nil {
			log.Printf("Unable to checkout original %s branch: %v: %s", repoDir, err, out)
		}
	}()

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
	// start with the latest release
	sort.Slice(tags, func(i, j int) bool {
		as, err := decomposeSemVer(tags[i])
		if err != nil {
			return false
		}
		bs, err := decomposeSemVer(tags[j])
		if err != nil {
			return true
		}
		for k := 0; k < 3; k++ {
			if as[k] == bs[k] {
				continue
			}
			return as[k] > bs[k]
		}
		return false
	})
	// format for tags git checkout
	for i, tag := range tags {
		tags[i] = "tags/" + tag
	}
	// fallback to original branch in case all tags fail to build
	tags = append(tags, origiBranch)

	// Keep trying building tags until a successful build
	for _, tag := range tags {

		cmd = exec.Command("git", "checkout", tag)
		cmd.Dir = repoDir
		if out, err = cmd.CombinedOutput(); err != nil {
			err = errors.Newf("checkout release at %s: %v: %s", tag, err, out)
			log.Printf("Unable to checkout release at %s", tag)
			continue
		}

		cmd = exec.Command("go", "build", "-o", buildFile)
		cmd.Dir = repoDir
		if out, err = cmd.CombinedOutput(); err != nil {
			err = errors.Newf("build release at %s: %v: %s", tag, err, out)
			log.Printf("Unable to build release at %s", tag)
			continue
		}
		return nil
	}

	return err
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

// decomposeSemVer breaks a semantic version into a slice of its
// constituent version numbers. The slice always has 3 values
// major version, minor version and patch version in this order.
// If the version missed one, it is filled by a zero.
// e.g. all these will return the slice [0,1,0]:
//     v0.1.0
//     v0.1
//     0.1.0
func decomposeSemVer(ver string) ([]int, error) {
	ver = strings.TrimPrefix(ver, "v")
	aStrs := strings.Split(ver, ".")
	var as []int
	for _, aStr := range aStrs {
		a, err := strconv.Atoi(aStr)
		if err != nil {
			return nil, errors.Newf("found none int in semantic version")
		}
		as = append(as, a)
	}
	if len(as) > 3 {
		return nil, errors.Newf("semantic version too long")
	}
	for i := len(as); i < 3; i++ {
		as = append(as, 0)
	}
	return as, nil
}
