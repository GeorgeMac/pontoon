package project

import (
	"bytes"
	"github.com/GeorgeMac/pontoon/archive"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

var (
	Stdout, Stderr *os.File = os.Stdout, os.Stderr
)

var buildProject func(string) Project = func(dir string) Project {
	return &GitProject{
		git: git(dir),
		dir: dir,
	}
}

type GitProjects struct {
	localDir string
	cache    map[string]Project
}

func NewGitProjects(localDir string) (g *GitProjects, err error) {
	g = &GitProjects{
		localDir: localDir,
		cache:    map[string]Project{},
	}

	stats, err := ioutil.ReadDir(localDir)
	for _, stat := range stats {
		fname := path.Join(localDir, stat.Name())
		p := buildProject(fname)
		if err = p.Pull(); err != nil {
			return
		}
		g.cache[fname] = p
	}

	return
}

func (g *GitProjects) Get(url string) (p Project, err error) {
	var ok bool
	_, name := path.Split(url)
	if p, ok = g.cache[path.Join(g.localDir, name)]; ok {
		p.Pull()
		return
	}

	return NewGitProject(g.localDir, url)
}

type GitProject struct {
	git GitCmdBuilder
	dir string
}

func NewGitProject(local, remote string) (g *GitProject, err error) {
	_, name := path.Split(remote)
	dir := path.Join(local, name)
	g = &GitProject{
		git: git(dir),
	}

	if err = git(local)(Stdout, Stderr, "clone", remote).Run(); err != nil {
		return
	}
	return
}

func (g *GitProject) WriteTo(wr io.Writer) error {
	return archive.Dir(g.dir, wr)
}

func (g *GitProject) Pull() error {
	return g.git(Stdout, Stderr, "pull").Run()
}

func (g *GitProject) Ref() (sha string, err error) {
	output := bytes.NewBuffer(nil)
	if err = g.git(output, Stderr, "rev-parse", "HEAD").Run(); err != nil {
		return
	}
	sha = string(output.Next(40))
	return
}

type GitCmdBuilder func(out, err io.Writer, args ...string) *exec.Cmd

func git(dir string) GitCmdBuilder {
	return func(out, err io.Writer, args ...string) (cmd *exec.Cmd) {
		cmd = exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Stdout = out
		cmd.Stderr = err
		return
	}

}
