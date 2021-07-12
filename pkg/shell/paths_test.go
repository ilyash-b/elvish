package shell

import (
	"os"
	"path/filepath"
	"testing"

	"src.elv.sh/pkg/env"
	"src.elv.sh/pkg/testutil"
)

var j = filepath.Join

func TestMakePaths_PopulatesUnsetSubPaths(t *testing.T) {
	paths := MakePaths(os.Stderr, Paths{RunDir: "run", DataDir: "data"})
	wantPaths := Paths{
		RunDir: "run",
		Sock:   j("run", "sock"),

		DataDir: "data",
		Db:      j("data", "db"),
		Rc:      j("data", "rc.elv"),
		LibDir:  j("data", "lib"),
	}
	if paths != wantPaths {
		t.Errorf("got paths %v, want %v", paths, wantPaths)
	}
}

func TestMakePaths_RespectsSetSubPaths(t *testing.T) {
	sock := "sock-override"
	paths := MakePaths(os.Stderr, Paths{
		RunDir: "run", DataDir: "data",
		Sock: sock,
	})
	if paths.Sock != sock {
		t.Errorf("got paths.Sock = %q, want %q", paths.Sock, sock)
	}
}

func TestMakePaths_SetsAndCreatesDataDir(t *testing.T) {
	home, cleanupDir := testutil.TestDir()
	defer cleanupDir()
	cleanupEnv := testutil.WithTempEnv(env.HOME, home)
	defer cleanupEnv()

	paths := MakePaths(os.Stderr, Paths{
		RunDir: "run",
	})

	wantDataDir := home + "/.elvish"
	if paths.DataDir != wantDataDir {
		t.Errorf("paths.DataDir = %q, want %q", paths.DataDir, wantDataDir)
	}

	stat, err := os.Stat(paths.DataDir)
	if err != nil {
		t.Errorf("could not stat %q: %v", paths.DataDir, err)
	}
	if !stat.IsDir() {
		t.Errorf("data dir %q is not dir", paths.DataDir)
	}
}
