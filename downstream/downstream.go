package downstream

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Package struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

type Downstream struct {
	Name           string `json:"name"`
	Version        string `json:"version"`
	DependsOn      string `json:"depends_on"`
	GitOriginFetch string `json:"git_origin_fetch"`
	DevDep         bool   `json:"dev_dep"`
}

type NodeDependency struct {
	Module  string
	Version string
	DevDep  bool
}

func (d *Downstream) clone(buildPath string) error {
	cmd := exec.Command("git", "clone", d.GitOriginFetch)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Dir = buildPath
	err := cmd.Run()
	errorStr := stderr.String()
	if err != nil && !strings.Contains(errorStr, "already exists") {
		return errors.New(errorStr)
	}
	return nil
}

func (d *Downstream) npmInstall(buildPath string, modulePath string) error {
	dirName, err := gitRemoteToDirName(d.GitOriginFetch)
	if err != nil {
		return err
	}
	repoDir := path.Join(buildPath, dirName)

	var cmd *exec.Cmd
	if modulePath != "" {
		cmd = exec.Command("npm", "install", modulePath)
	} else {
		cmd = exec.Command("npm", "install")
	}
	cmd.Dir = repoDir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return errors.New(stderr.String())
	}
	return nil
}

func (d *Downstream) Build(basePath string) error {
	createBuildPath(basePath)
	buildPath := path.Join(basePath, ".downstream")
	log.Printf("Cloning into %s...", d.GitOriginFetch)
	if err := d.clone(buildPath); err != nil {
		return err
	}

	log.Println("Npm installing...")
	if err := d.npmInstall(buildPath, basePath); err != nil {
		return err
	}

	if err := d.npmInstall(buildPath, ""); err != nil {
		return err
	}

	return nil
}

func (d *Downstream) Test(basePath string, verbose bool) error {
	log.Printf("Running tests for %s...", d.Name)
	buildPath := path.Join(basePath, ".downstream")
	dirName, err := gitRemoteToDirName(d.GitOriginFetch)
	if err != nil {
		return err
	}
	repoDir := path.Join(buildPath, dirName)

	cmd := exec.Command("make", "test")
	cmd.Dir = repoDir

	if verbose {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}

	err = cmd.Run()
	if err != nil {
		log.Printf("Tests failed: %s", err)
	} else {
		log.Print("Tests pass!")
	}
	return nil
}

func gitRemoteToDirName(gitRemote string) (string, error) {
	parsed, err := url.Parse(gitRemote)
	if err != nil {
		return "", err
	}
	p := parsed.Path
	dirName := path.Base(p)
	dirName = strings.TrimSuffix(dirName, ".git")
	return dirName, nil
}

func createBuildPath(basePath string) error {
	err := os.Mkdir(path.Join(basePath, ".downstream"), 0777)
	if !os.IsExist(err) {
		return err
	}
	return nil
}

func LoadPackage(filePath string) (Package, error) {
	var p Package

	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		return p, err
	}

	if err := json.Unmarshal(raw, &p); err != nil {
		return p, err
	}

	return p, nil
}

func hasUpstreamDepOn(pack Package, moduleName string) (bool, NodeDependency) {
	for name, version := range pack.Dependencies {
		if name == moduleName {
			dep := NodeDependency{
				Module:  name,
				Version: version,
				DevDep:  false,
			}
			return true, dep
		}
	}

	for name, version := range pack.DevDependencies {
		if name == moduleName {
			dep := NodeDependency{
				Module:  name,
				Version: version,
				DevDep:  true,
			}
			return true, dep
		}
	}
	return false, NodeDependency{}
}

func listSubDirectories(basePath string) ([]string, error) {
	subdirs := []string{}
	infos, err := ioutil.ReadDir(basePath)
	if err != nil {
		return subdirs, err
	}
	for _, info := range infos {
		if info.IsDir() {
			subdirs = append(subdirs, path.Join(basePath, info.Name()))
		}
	}
	return subdirs, nil
}

func getGitOrigin(basePath string, ds *Downstream) error {
	cmd := exec.Command("git", "remote", "-v")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Dir = basePath
	err := cmd.Run()
	if err != nil {
		return err
	}
	total := out.String()
	lines := strings.Split(total, "\n")

	for _, remote := range lines {
		tokens := strings.Split(remote, "\t")
		if tokens[0] == "origin" && strings.Contains(tokens[1], "(fetch)") {
			parts := strings.Split(tokens[1], " ")
			ds.GitOriginFetch = parts[0]
			return nil
		}
	}
	return errors.New("Could not determine git origin fetch url")
}

func IsNodeDir(basePath string) bool {
	if _, err := os.Stat(path.Join(basePath, "package.json")); err != nil {
		return false
	}
	return true
}

func List(basePath string, moduleName string) ([]Downstream, error) {
	downstreams := []Downstream{}

	subdirs, err := listSubDirectories(basePath)
	if err != nil {
		return downstreams, err
	}

	for _, subdir := range subdirs {
		if !IsNodeDir(subdir) {
			continue
		}

		pack, err := LoadPackage(path.Join(subdir, "package.json"))
		if err != nil {
			return downstreams, err
		}

		hasDep, dep := hasUpstreamDepOn(pack, moduleName)
		if !hasDep {
			continue
		}

		ds := Downstream{
			Name:      pack.Name,
			Version:   pack.Version,
			DependsOn: dep.Version,
			DevDep:    dep.DevDep,
		}

		if err := getGitOrigin(subdir, &ds); err != nil {
			return downstreams, err
		}

		downstreams = append(downstreams, ds)
	}

	return downstreams, nil
}
