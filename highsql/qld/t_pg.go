package qld

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/araddon/qlbridge/expr"
	"github.com/araddon/qlbridge/rel"
	"github.com/k0kubun/pp"
)

func pgTransform(tenantId, groupId, query string) (string, error) {

	pp.Println("Before")

	st, err := rel.ParseSqlSelect(query)
	if err != nil {
		return "", err
	}
	pp.Println("After")

	wctx := expr.NewDefaultWriter()

	pp.Println("...", query)

	st.WriteDialect(wctx)

	return wctx.String(), nil
}

// executeECPGCommand executes a sql statement against the ecpg tool
func executeECPGCommand(stmt string) error {

	defer pp.Println("done", stmt)

	stmt = fmt.Sprintf("EXEC SQL %s;", stmt)

	args := []string{"-o", "-", "-"}

	pp.Println("Executing ecpg", args)

	cmd := exec.Command("ecpg", args...)

	pp.Println("Execute done")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, stmt)
		pp.Println("wrote =>", stmt)
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(pullErrorData(string(out)))
	}
	return nil
}

// pullErrorData grabs the first line of the out text (the actual error)
// then removes the common, unhelpful location indicator and returns
func pullErrorData(out string) string {
	outs := strings.Split(out, "\n")

	for _, v := range outs {
		if strings.HasPrefix(v, "stdin:1:") {
			// remove `stdin:1:` & return
			return strings.TrimSpace(strings.Replace(v, "stdin:1:", "", 1))
		}
	}

	return ""
}
