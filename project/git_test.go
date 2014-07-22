package project

import (
	"fmt"
	"io"
	"io/ioutil"
	"launchpad.net/gocheck"
	"os"
	"path"
	"sort"
	"testing"
)

var _ = gocheck.Suite(&GitSuite{})

func Test(t *testing.T) { gocheck.TestingT(t) }

type GitSuite struct {
}

func (g *GitSuite) SetUpTest(c *gocheck.C) {

}

func (g *GitSuite) Test_NewGitProjects(c *gocheck.C) {
	old := buildProject
	defer func() {
		buildProject = old
	}()

	buildProject = func(ref string) Project {
		return &DummyProject{ref: ref}
	}

	tmp := c.MkDir()

	expected := []string{}

	ftemp := path.Join(tmp, "test%d")
	for i := 0; i < 5; i++ {
		fname := fmt.Sprintf(ftemp, i)
		if err := ioutil.WriteFile(fname, []byte(""), os.ModePerm); err != nil {
			panic(err)
		}
		expected = append(expected, fname)
	}

	projects, err := NewGitProjects(tmp)
	c.Assert(err, gocheck.IsNil)

	obtained := []string{}
	for k, _ := range projects.cache {
		obtained = append(obtained, k)
	}

	sort.Strings(obtained)
	c.Check(obtained, gocheck.DeepEquals, expected)
}

type DummyProject struct {
	ref string
}

func (d *DummyProject) WriteTo(_ io.Writer) error { return nil }
func (d *DummyProject) Ref() (string, error)      { return "", nil }
func (d *DummyProject) Pull() error {
	fmt.Println("[TEST] Pull called for", d.ref)
	return nil
}
