package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/deislabs/duffle/pkg/duffle/home"
	"github.com/deislabs/duffle/pkg/loader"
	"github.com/deislabs/duffle/pkg/packager"
)

const importDesc = `
Imports a Cloud Native Application Bundle (CNAB) from a gzipped tar file on the local file system.

The bundle is stored in local storage unless --destination is specified, in which case the file will be unzipped to the
specified path. Note: if --destination is specified as the empty string, the file will be unzipped to local storage.

Example:
	$ duffle import mybundle.tgz
	$ duffle import mybundle.tgz -d bundles/unzipped
`

type importCmd struct {
	source  string
	dest    string
	out     io.Writer
	home    home.Home
	verbose bool
}

func newImportCmd(w io.Writer) *cobra.Command {
	importc := &importCmd{
		out:  w,
		home: home.Home(homePath()),
	}

	cmd := &cobra.Command{
		Use:   "import [PATH]",
		Short: "extract CNAB bundle from gzipped tar file",
		Long:  importDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("this command requires the path to the packaged bundle")
			}
			importc.source = args[0]

			return importc.run()
		},
	}

	f := cmd.Flags()
	f.StringVarP(&importc.dest, "destination", "d", "", "Location to unpack bundle")
	f.BoolVarP(&importc.verbose, "verbose", "v", false, "Verbose output")

	return cmd
}

func (im *importCmd) run() error {
	source, err := filepath.Abs(im.source)
	if err != nil {
		return fmt.Errorf("Error in source: %s", err)
	}

	var dest string
	if im.dest == "" {
		dest, err = ioutil.TempDir("", "")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dest)
	} else {
		dest, err = filepath.Abs(im.dest)
		if err != nil {
			return fmt.Errorf("Error in destination: %s", err)
		}
	}

	l := loader.NewLoader()
	imp, err := packager.NewImporter(source, dest, l, im.verbose)
	if err != nil {
		return err
	}
	return imp.Import(im.home)
}
