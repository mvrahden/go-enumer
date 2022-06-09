package linter

import (
	"go/ast"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/mvrahden/go-enumer/pkg/common"
	"github.com/mvrahden/go-enumer/pkg/utils"
)

func init() {
	common.ExtractCommentString = func(c *ast.Comment) string {
		// hint: strip any test artifacts
		els := strings.Split(c.Text, " // want ")
		return els[0]
	}
}

func TestAll(t *testing.T) {
	wd := utils.Must(os.Getwd())

	testdata := filepath.Join(wd, "testdata")
	analysistest.Run(t, testdata, New(&Config{}), "basic", "enums", "nonenums")
}

func TestCsv(t *testing.T) {
	wd := utils.Must(os.Getwd())

	testdata := filepath.Join(wd, "testdata")
	analysistest.Run(t, testdata, New(&Config{}), "csv_no_file", "csv")
}
