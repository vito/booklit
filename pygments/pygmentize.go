package pygments

import (
	"bytes"
	"encoding/xml"
	"os"
	"os/exec"
)

type highlightedHTML struct {
	Pre struct {
		Code string `xml:",innerxml"`
	} `xml:"pre"`
}

func pygmentize(language string, content string) (string, error) {
	html := new(bytes.Buffer)

	pygmentize := exec.Command("pygmentize", "-l", language, "-f", "html", "-O", "encoding=utf-8")
	pygmentize.Stdin = bytes.NewBufferString(content)
	pygmentize.Stdout = html
	pygmentize.Stderr = os.Stderr
	err := pygmentize.Run()
	if err != nil {
		return "", err
	}

	var hl highlightedHTML
	err = xml.Unmarshal(html.Bytes(), &hl)
	if err != nil {
		return "", err
	}

	return hl.Pre.Code, nil
}
