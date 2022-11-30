package linter

import (
	"go/ast"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/mvrahden/go-enumer/pkg/enumer"
	"github.com/mvrahden/go-enumer/pkg/utils"
)

func init() {
	enumer.ExtractCommentString = func(c *ast.Comment) string {
		// hint: strip any test artifacts
		els := strings.Split(c.Text, " // want ")
		return els[0]
	}
}

func Test_Linter_All(t *testing.T) {
	wd := utils.Must(os.Getwd())

	testdata := filepath.Join(wd, "testdata")
	analysistest.Run(t, testdata, New(&Config{}), "basic", "enums", "nonenums")
}

func Test_Linter_Csv(t *testing.T) {
	wd := utils.Must(os.Getwd())

	testdata := filepath.Join(wd, "testdata")
	analysistest.Run(t, testdata, New(&Config{}), "csv_no_file", "csv")
}
