package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/maxivhuber/monkeyi/lexer"
	"github.com/maxivhuber/monkeyi/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l, err := lexer.New(line)
		if err != nil {
			fmt.Fprint(out, err.Error())
		}
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
