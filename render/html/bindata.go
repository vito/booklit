// Code generated for package html by go-bindata DO NOT EDIT. (@generated)
// sources:
// render/html/aside.tmpl
// render/html/bold.tmpl
// render/html/definitions.tmpl
// render/html/image.tmpl
// render/html/inset.tmpl
// render/html/italic.tmpl
// render/html/larger.tmpl
// render/html/link.tmpl
// render/html/list.tmpl
// render/html/page.tmpl
// render/html/paragraph.tmpl
// render/html/preformatted.tmpl
// render/html/reference.tmpl
// render/html/section.tmpl
// render/html/sequence.tmpl
// render/html/smaller.tmpl
// render/html/strike.tmpl
// render/html/string.tmpl
// render/html/subscript.tmpl
// render/html/superscript.tmpl
// render/html/table.tmpl
// render/html/target.tmpl
// render/html/toc.tmpl
// render/html/verbatim.tmpl
package html

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _renderHtmlAsideTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb2\x49\xca\xc9\x4f\xce\x2e\x2c\xcd\x2f\x49\x55\x48\xce\x49\x2c\x2e\xb6\x55\x4a\x2c\xce\x4c\x49\x55\xb2\xab\xae\xd6\x73\xce\xcf\x2b\x49\xcd\x2b\x51\xa8\x51\x28\x4a\xcd\x4b\x49\x2d\xaa\xad\xb5\xd1\x47\xa8\xb7\xe3\x02\x04\x00\x00\xff\xff\xc9\x06\xb2\xf8\x3d\x00\x00\x00")

func renderHtmlAsideTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderHtmlAsideTmpl,
		"render/html/aside.tmpl",
	)
}

func renderHtmlAsideTmpl() (*asset, error) {
	bytes, err := renderHtmlAsideTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/html/aside.tmpl", size: 61, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderHtmlBoldTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb2\x29\x2e\x29\xca\xcf\x4b\xb7\xab\xae\xd6\x73\xce\xcf\x2b\x49\xcd\x2b\x51\xa8\x51\x28\x4a\xcd\x4b\x49\x2d\xaa\xad\xb5\xd1\x87\xca\x72\x01\x02\x00\x00\xff\xff\x6c\x89\x4d\x8f\x27\x00\x00\x00")

func renderHtmlBoldTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderHtmlBoldTmpl,
		"render/html/bold.tmpl",
	)
}

func renderHtmlBoldTmpl() (*asset, error) {
	bytes, err := renderHtmlBoldTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/html/bold.tmpl", size: 39, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderHtmlDefinitionsTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb2\x49\xc9\xb1\xe3\x52\x50\xa8\xae\x2e\x4a\xcc\x4b\x4f\x55\xd0\xab\xad\xe5\x52\x50\xb0\x49\x29\xb1\xab\xae\xd6\x0b\x2e\x4d\xca\x4a\x4d\x2e\x51\xa8\x51\x28\x4a\xcd\x4b\x49\x2d\xaa\xad\xb5\xd1\x4f\x29\x01\x29\x07\x29\x49\x01\x29\x71\x49\x4d\xcb\xcc\xcb\x2c\xc9\xcc\xcf\x43\x55\x95\x02\x31\x34\x35\x2f\xa5\xb6\x96\xcb\x46\x1f\x64\x09\x20\x00\x00\xff\xff\x3f\x6e\x87\x35\x6a\x00\x00\x00")

func renderHtmlDefinitionsTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderHtmlDefinitionsTmpl,
		"render/html/definitions.tmpl",
	)
}

func renderHtmlDefinitionsTmpl() (*asset, error) {
	bytes, err := renderHtmlDefinitionsTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/html/definitions.tmpl", size: 106, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderHtmlImageTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb2\xc9\xcc\x4d\x57\x28\x2e\x4a\xb6\x55\xaa\xae\xd6\x0b\x48\x2c\xc9\xa8\xad\x55\x52\x48\xcc\x29\x01\xf3\x5d\x52\x8b\x93\x8b\x32\x0b\x4a\x32\xf3\xf3\x40\xc2\xfa\x76\x5c\x80\x00\x00\x00\xff\xff\x6c\xbe\xfe\x41\x2f\x00\x00\x00")

func renderHtmlImageTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderHtmlImageTmpl,
		"render/html/image.tmpl",
	)
}

func renderHtmlImageTmpl() (*asset, error) {
	bytes, err := renderHtmlImageTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/html/image.tmpl", size: 47, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderHtmlInsetTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x04\xc0\x5d\x0a\x80\x20\x0c\x00\xe0\xf7\x4e\x31\x76\x80\xfe\x1e\x43\x7d\xe9\x24\x92\x23\x06\xba\xc0\x0d\x21\xcc\xbb\xf7\xb9\xc4\x0d\xd4\xde\x4c\x1e\x4b\xac\x37\xcb\x01\x2b\xec\x54\x60\xa3\x82\x70\xe5\xa8\xea\x91\x45\xc9\x30\xf4\x3e\x9f\x8f\x18\x89\xc1\x07\x95\x24\x51\x1d\xc3\x2d\x89\x5b\x98\xfe\x00\x00\x00\xff\xff\x05\xaa\x65\xcc\x49\x00\x00\x00")

func renderHtmlInsetTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderHtmlInsetTmpl,
		"render/html/inset.tmpl",
	)
}

func renderHtmlInsetTmpl() (*asset, error) {
	bytes, err := renderHtmlInsetTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/html/inset.tmpl", size: 73, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderHtmlItalicTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb2\x49\xcd\xb5\xab\xae\xd6\x73\xce\xcf\x2b\x49\xcd\x2b\x51\xa8\x51\x28\x4a\xcd\x4b\x49\x2d\xaa\xad\xb5\xd1\x4f\xcd\xb5\xe3\x02\x04\x00\x00\xff\xff\x47\x8b\x68\x8b\x1f\x00\x00\x00")

func renderHtmlItalicTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderHtmlItalicTmpl,
		"render/html/italic.tmpl",
	)
}

func renderHtmlItalicTmpl() (*asset, error) {
	bytes, err := renderHtmlItalicTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/html/italic.tmpl", size: 31, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderHtmlLargerTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb2\x29\x2e\x48\xcc\x53\x28\x2e\xa9\xcc\x49\xb5\x55\x4a\xcb\xcf\x2b\xd1\x2d\xce\xac\x4a\xb5\x52\x30\x34\x32\x50\x55\xb2\xab\xae\xd6\x73\xce\xcf\x2b\x49\xcd\x2b\x51\xa8\x51\x28\x4a\xcd\x4b\x49\x2d\xaa\xad\xb5\xd1\x07\xe9\xb1\xe3\x02\x04\x00\x00\xff\xff\x03\xa7\xd9\x06\x3b\x00\x00\x00")

func renderHtmlLargerTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderHtmlLargerTmpl,
		"render/html/larger.tmpl",
	)
}

func renderHtmlLargerTmpl() (*asset, error) {
	bytes, err := renderHtmlLargerTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/html/larger.tmpl", size: 59, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderHtmlLinkTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb2\x49\x54\xc8\x28\x4a\x4d\xb3\x55\xaa\xae\xd6\x0b\x49\x2c\x4a\x4f\x2d\xa9\xad\x55\xb2\xab\xae\xd6\x73\xce\xcf\x2b\x49\xcd\x2b\x51\xa8\x51\x28\x4a\xcd\x4b\x49\x2d\xaa\xad\xb5\xd1\x4f\xb4\xe3\x02\x04\x00\x00\xff\xff\x30\x04\xce\x2f\x30\x00\x00\x00")

func renderHtmlLinkTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderHtmlLinkTmpl,
		"render/html/link.tmpl",
	)
}

func renderHtmlLinkTmpl() (*asset, error) {
	bytes, err := renderHtmlLinkTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/html/link.tmpl", size: 48, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderHtmlListTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x54\xcc\x31\x0a\x82\x31\x10\x44\xe1\x3e\xa7\xd8\x13\x24\x17\x08\xe9\xad\x3c\x83\xb0\xa3\x04\xd6\x08\xfb\x6b\x35\xce\xdd\x05\x1b\x63\xf7\x9a\xf7\x91\xf3\x6a\xf5\x9c\x8e\x84\x4b\xfd\x11\x83\x44\x1c\x90\xfa\xeb\xdb\xcb\xa5\x42\xe6\x65\xdd\x60\xf5\xf4\xc4\xfd\x90\x8a\x59\x8f\x39\xc8\x6a\x6f\x4b\x2c\x47\x4a\xbd\xc5\x1c\xe5\xb7\xfc\xcb\x6d\xa7\xdb\x66\x7f\x02\x00\x00\xff\xff\x26\x88\xac\xef\x83\x00\x00\x00")

func renderHtmlListTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderHtmlListTmpl,
		"render/html/list.tmpl",
	)
}

func renderHtmlListTmpl() (*asset, error) {
	bytes, err := renderHtmlListTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/html/list.tmpl", size: 131, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderHtmlPageTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x34\x8e\x41\xaf\x82\x30\x10\x84\xef\xef\x57\xec\xeb\x1d\x7a\x7d\xc9\x5b\xb8\xa8\x67\x4d\xe4\xe2\x11\x61\xb5\x24\x50\xb0\x0e\x46\x52\xfb\xdf\x4d\x2d\x9e\x76\x36\xb3\x3b\xdf\xf0\xef\x76\xbf\xa9\x4e\x87\x1d\x19\x0c\x7d\xf9\xc3\x69\x10\xb1\x91\xba\x8d\x82\x88\x07\x41\x4d\x06\x98\x32\xb9\xcd\xdd\xa3\x50\xcd\x68\x21\x16\x19\x96\x49\x14\xad\x5b\xa1\x20\x4f\xe8\x18\xf0\x4f\x8d\xa9\xdd\x5d\x50\xcc\xb8\x64\x7f\x8a\xf4\x9a\x84\x0e\xbd\x94\xde\xe7\x55\x14\xf9\x11\xae\xb3\xd7\x10\x58\x27\x23\x72\xf5\x17\xcc\xe7\xb1\x5d\xd2\x9f\xf7\x39\xbd\xc8\x89\x6d\xc5\x85\xf0\xb9\x4a\x26\xeb\xd4\xf7\x1d\x00\x00\xff\xff\xec\xca\x7c\xf9\xc7\x00\x00\x00")

func renderHtmlPageTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderHtmlPageTmpl,
		"render/html/page.tmpl",
	)
}

func renderHtmlPageTmpl() (*asset, error) {
	bytes, err := renderHtmlPageTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/html/page.tmpl", size: 199, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderHtmlParagraphTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb2\x29\xb0\xab\xae\x2e\x4a\xcc\x4b\x4f\x55\x50\xc9\xcc\x4b\x49\xad\xd0\x51\x50\xc9\xc9\xcc\x4b\x55\xb0\xb2\x55\xd0\xab\xad\xad\xae\xce\x4c\x83\x4a\xd4\xd6\x2a\x54\x57\xa7\xe6\xa5\x80\x04\x21\x4a\x6a\x14\x8a\x52\xf3\x52\x52\x8b\x40\x22\x60\x09\x1b\xfd\x02\x3b\x2e\x40\x00\x00\x00\xff\xff\x70\x0e\xcc\x54\x52\x00\x00\x00")

func renderHtmlParagraphTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderHtmlParagraphTmpl,
		"render/html/paragraph.tmpl",
	)
}

func renderHtmlParagraphTmpl() (*asset, error) {
	bytes, err := renderHtmlParagraphTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/html/paragraph.tmpl", size: 82, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderHtmlPreformattedTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xaa\xae\x2e\x4a\xcc\x4b\x4f\x55\x50\xc9\xcc\x4b\x49\xad\xd0\x51\x50\xc9\xc9\xcc\x4b\x55\xb0\xb2\x55\xd0\xab\xad\xad\xae\xce\x4c\x83\x4a\xd4\xd6\x72\x55\x57\xa7\xe6\xa5\x80\x04\x21\x4a\x6a\x14\x8a\x52\xf3\x52\x52\x8b\x40\x22\x60\x09\x2e\x40\x00\x00\x00\xff\xff\x41\x45\x7e\xd9\x4b\x00\x00\x00")

func renderHtmlPreformattedTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderHtmlPreformattedTmpl,
		"render/html/preformatted.tmpl",
	)
}

func renderHtmlPreformattedTmpl() (*asset, error) {
	bytes, err := renderHtmlPreformattedTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/html/preformatted.tmpl", size: 75, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderHtmlReferenceTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb2\x49\x54\xc8\x28\x4a\x4d\xb3\x55\xaa\xae\xd6\x0b\x49\x4c\x57\xa8\x51\x28\x2d\xca\xa9\xad\x55\xb2\xab\xae\xd6\x73\xc9\x2c\x2e\xc8\x49\xac\x54\xa8\x51\x28\x2e\x29\xca\x2c\x70\x2c\xad\x50\xa8\x51\x28\x4a\xcd\x4b\x49\x2d\xaa\xad\xb5\xd1\x4f\xb4\xe3\x02\x04\x00\x00\xff\xff\x6c\xcd\x19\x93\x3e\x00\x00\x00")

func renderHtmlReferenceTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderHtmlReferenceTmpl,
		"render/html/reference.tmpl",
	)
}

func renderHtmlReferenceTmpl() (*asset, error) {
	bytes, err := renderHtmlReferenceTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/html/reference.tmpl", size: 62, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderHtmlSectionTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x64\xce\x41\x6a\xc3\x30\x10\x05\xd0\xbd\x4f\xf1\xf1\xde\xe3\x0b\x28\x5e\xb4\x5d\x87\x40\x72\x81\x69\x34\x89\x05\xf6\xd8\xc8\xea\x22\x4c\xe7\xee\xc5\x89\x53\x4a\xb3\x92\xf8\xf3\x79\xfc\xd0\x9b\xf5\xc2\x51\xf2\x87\xcc\xa5\x07\xb9\xe3\x3c\xf0\xb2\xec\xea\x45\xce\x25\x4d\xda\x3c\xce\x75\x17\x18\x29\xee\x6a\x33\x3a\xe4\x34\x72\xbe\x9d\xf8\x4a\x7b\x1e\xc5\xbd\xee\x42\xcb\x5d\x05\x98\x35\x48\x17\xd0\x81\xb3\x68\x41\xe3\x5e\x01\x40\x58\x66\xd6\xff\xae\x7e\x8d\x9f\xab\x6b\x46\xfb\xfb\xd7\x1d\xa1\x5d\x9b\x4f\x49\x34\x6e\x84\x19\x9d\x52\x19\x04\xdf\xc8\xa2\x51\xf2\x3d\x0f\xed\xeb\xfa\xae\xaa\xcc\xe8\x6d\x8a\xb7\xdf\xae\xfb\x9a\xa5\x0b\x74\x2a\xa0\xe3\x3c\xa4\x72\x7c\x6c\x58\x36\x3c\xb3\x5e\x05\xf4\xde\xa7\x21\x66\xd1\x6d\xb5\x19\xfd\x35\xd6\x40\x34\xba\x57\xcf\xf7\x27\x00\x00\xff\xff\xa7\xb5\xd8\xf7\x3d\x01\x00\x00")

func renderHtmlSectionTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderHtmlSectionTmpl,
		"render/html/section.tmpl",
	)
}

func renderHtmlSectionTmpl() (*asset, error) {
	bytes, err := renderHtmlSectionTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/html/section.tmpl", size: 317, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderHtmlSequenceTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xaa\xae\x2e\x4a\xcc\x4b\x4f\x55\xd0\x73\xce\xcf\x2b\x49\xcd\x2b\x29\xae\xad\xad\xae\xd6\x53\xa8\x51\x28\x4a\xcd\x4b\x49\x2d\x02\xf1\x52\xf3\x52\x6a\x6b\xb9\x00\x01\x00\x00\xff\xff\xfc\xbe\x50\x06\x29\x00\x00\x00")

func renderHtmlSequenceTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderHtmlSequenceTmpl,
		"render/html/sequence.tmpl",
	)
}

func renderHtmlSequenceTmpl() (*asset, error) {
	bytes, err := renderHtmlSequenceTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/html/sequence.tmpl", size: 41, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderHtmlSmallerTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb2\x29\x2e\x48\xcc\x53\x28\x2e\xa9\xcc\x49\xb5\x55\x4a\xcb\xcf\x2b\xd1\x2d\xce\xac\x4a\xb5\x52\xb0\x30\x50\x55\xb2\xab\xae\xd6\x73\xce\xcf\x2b\x49\xcd\x2b\x51\xa8\x51\x28\x4a\xcd\x4b\x49\x2d\xaa\xad\xb5\xd1\x07\x69\xb1\xe3\x02\x04\x00\x00\xff\xff\x0c\x76\xc9\xfa\x3a\x00\x00\x00")

func renderHtmlSmallerTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderHtmlSmallerTmpl,
		"render/html/smaller.tmpl",
	)
}

func renderHtmlSmallerTmpl() (*asset, error) {
	bytes, err := renderHtmlSmallerTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/html/smaller.tmpl", size: 58, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderHtmlStrikeTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb2\x29\x2e\x48\xcc\x53\x28\x2e\xa9\xcc\x49\xb5\x55\x2a\x49\xad\x28\xd1\x4d\x49\x4d\xce\x2f\x4a\x2c\xc9\xcc\xcf\xb3\x52\xc8\xc9\xcc\x4b\xd5\x2d\xc9\x28\xca\x2f\x4d\xcf\x50\xb2\xab\xae\xd6\x73\xce\xcf\x2b\x49\xcd\x2b\x51\xa8\x51\x28\x4a\xcd\x4b\x49\x2d\xaa\xad\xb5\xd1\x07\x99\x60\xc7\x05\x08\x00\x00\xff\xff\x1b\xce\x7a\x76\x49\x00\x00\x00")

func renderHtmlStrikeTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderHtmlStrikeTmpl,
		"render/html/strike.tmpl",
	)
}

func renderHtmlStrikeTmpl() (*asset, error) {
	bytes, err := renderHtmlStrikeTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/html/strike.tmpl", size: 73, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderHtmlStringTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xaa\xae\xd6\x0b\x2e\x29\xca\xcc\x4b\xaf\xad\xe5\x02\x04\x00\x00\xff\xff\x38\x6b\xcd\x33\x0c\x00\x00\x00")

func renderHtmlStringTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderHtmlStringTmpl,
		"render/html/string.tmpl",
	)
}

func renderHtmlStringTmpl() (*asset, error) {
	bytes, err := renderHtmlStringTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/html/string.tmpl", size: 12, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderHtmlSubscriptTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb2\x29\x2e\x4d\xb2\xab\xae\xd6\x73\xce\xcf\x2b\x49\xcd\x2b\x51\xa8\x51\x28\x4a\xcd\x4b\x49\x2d\xaa\xad\xb5\xd1\x07\x49\x71\x01\x02\x00\x00\xff\xff\xb4\x8b\xde\x56\x21\x00\x00\x00")

func renderHtmlSubscriptTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderHtmlSubscriptTmpl,
		"render/html/subscript.tmpl",
	)
}

func renderHtmlSubscriptTmpl() (*asset, error) {
	bytes, err := renderHtmlSubscriptTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/html/subscript.tmpl", size: 33, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderHtmlSuperscriptTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb2\x29\x2e\x2d\xb0\xab\xae\xd6\x73\xce\xcf\x2b\x49\xcd\x2b\x51\xa8\x51\x28\x4a\xcd\x4b\x49\x2d\xaa\xad\xb5\xd1\x07\x49\x71\x01\x02\x00\x00\xff\xff\xb8\x06\x51\x41\x21\x00\x00\x00")

func renderHtmlSuperscriptTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderHtmlSuperscriptTmpl,
		"render/html/superscript.tmpl",
	)
}

func renderHtmlSuperscriptTmpl() (*asset, error) {
	bytes, err := renderHtmlSuperscriptTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/html/superscript.tmpl", size: 33, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderHtmlTableTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb2\x29\x49\x4c\xca\x49\xb5\xe3\x52\x50\xa8\xae\x2e\x4a\xcc\x4b\x4f\x55\xd0\x0b\xca\x2f\x2f\xae\xad\xe5\x52\x50\xb0\x29\x29\x02\x49\x20\x49\x81\x85\x41\x12\x29\x76\xd5\xd5\x7a\x0a\x35\x0a\x45\xa9\x79\x29\xa9\x45\xb5\xb5\x36\xfa\x25\x29\x30\xb5\xa9\x79\x29\x10\xed\xfa\x10\xfd\x30\x11\x1b\x7d\xa8\x5d\x80\x00\x00\x00\xff\xff\xa9\x69\xd9\x6b\x74\x00\x00\x00")

func renderHtmlTableTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderHtmlTableTmpl,
		"render/html/table.tmpl",
	)
}

func renderHtmlTableTmpl() (*asset, error) {
	bytes, err := renderHtmlTableTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/html/table.tmpl", size: 116, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderHtmlTargetTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb2\x49\x54\xc8\x4c\xb1\x55\xaa\xae\xd6\x0b\x49\x4c\xf7\x4b\xcc\x4d\xad\xad\x55\xb2\xb3\xd1\x4f\xb4\xe3\x02\x04\x00\x00\xff\xff\x60\x82\x12\x78\x1a\x00\x00\x00")

func renderHtmlTargetTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderHtmlTargetTmpl,
		"render/html/target.tmpl",
	)
}

func renderHtmlTargetTmpl() (*asset, error) {
	bytes, err := renderHtmlTargetTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/html/target.tmpl", size: 26, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderHtmlTocTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x44\x8e\x41\x6a\xc3\x30\x10\x45\xf7\x3e\xc5\xc7\xab\x76\x23\x5f\x40\x35\x94\x40\x97\x4d\x17\xbe\x80\x52\x8f\x13\x81\x34\x36\x93\x51\x69\x19\xeb\xee\xc5\x6e\x4d\x56\x33\xbc\x0f\xff\x3f\xb3\x38\x21\xf0\x08\x77\xba\xc5\x34\x0a\x31\x9e\x78\x56\xb8\x73\x8e\x7a\xa0\x37\x99\xf3\x10\x2e\x89\xce\xd3\x69\x66\x25\xd6\xfb\x73\xad\x8d\xe7\xf0\xd5\x37\x80\x2f\x69\x3b\x66\x12\xf8\x4a\x8f\xa6\x5a\x1b\x00\xf0\x29\xf6\xfb\x03\xf8\x80\x9b\xd0\xf4\xd2\x9a\xb9\x0f\x89\x39\xc8\xcf\x10\xae\x58\x51\x24\xd5\xda\xf6\x66\xee\xbd\xe4\x0b\x49\xad\x30\x73\x43\xd4\x44\x58\x71\x57\x89\xcb\x6b\xf9\xc6\x0a\x21\x1e\xb7\xd8\x77\xa1\x6f\xfe\x5b\xcd\x94\xf2\x92\x82\x12\x5a\x9d\x3f\x9d\xe6\x25\xb5\x70\xc7\x7c\xf7\xb7\x6f\x46\x3c\xee\xcc\x77\x9b\xb0\xef\x76\xfd\x03\xff\x06\x00\x00\xff\xff\xc5\x44\xc0\xe1\x09\x01\x00\x00")

func renderHtmlTocTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderHtmlTocTmpl,
		"render/html/toc.tmpl",
	)
}

func renderHtmlTocTmpl() (*asset, error) {
	bytes, err := renderHtmlTocTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/html/toc.tmpl", size: 265, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderHtmlVerbatimTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xaa\xae\xce\x4c\x53\xd0\xf3\x2c\x76\xcb\xc9\x2f\xaf\xad\xb5\x49\xce\x4f\x49\xb5\xab\xae\xd6\x73\xce\xcf\x2b\x49\xcd\x2b\x51\xa8\x51\x28\x4a\xcd\x4b\x49\x2d\xaa\xad\xb5\xd1\x87\xca\xa5\xe6\x14\xa7\xd6\xd6\xda\x14\x14\xe1\x54\x08\x91\x4a\xcd\x4b\xa9\xad\xe5\x02\x04\x00\x00\xff\xff\xe0\xa0\x44\x54\x60\x00\x00\x00")

func renderHtmlVerbatimTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderHtmlVerbatimTmpl,
		"render/html/verbatim.tmpl",
	)
}

func renderHtmlVerbatimTmpl() (*asset, error) {
	bytes, err := renderHtmlVerbatimTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/html/verbatim.tmpl", size: 96, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"render/html/aside.tmpl":        renderHtmlAsideTmpl,
	"render/html/bold.tmpl":         renderHtmlBoldTmpl,
	"render/html/definitions.tmpl":  renderHtmlDefinitionsTmpl,
	"render/html/image.tmpl":        renderHtmlImageTmpl,
	"render/html/inset.tmpl":        renderHtmlInsetTmpl,
	"render/html/italic.tmpl":       renderHtmlItalicTmpl,
	"render/html/larger.tmpl":       renderHtmlLargerTmpl,
	"render/html/link.tmpl":         renderHtmlLinkTmpl,
	"render/html/list.tmpl":         renderHtmlListTmpl,
	"render/html/page.tmpl":         renderHtmlPageTmpl,
	"render/html/paragraph.tmpl":    renderHtmlParagraphTmpl,
	"render/html/preformatted.tmpl": renderHtmlPreformattedTmpl,
	"render/html/reference.tmpl":    renderHtmlReferenceTmpl,
	"render/html/section.tmpl":      renderHtmlSectionTmpl,
	"render/html/sequence.tmpl":     renderHtmlSequenceTmpl,
	"render/html/smaller.tmpl":      renderHtmlSmallerTmpl,
	"render/html/strike.tmpl":       renderHtmlStrikeTmpl,
	"render/html/string.tmpl":       renderHtmlStringTmpl,
	"render/html/subscript.tmpl":    renderHtmlSubscriptTmpl,
	"render/html/superscript.tmpl":  renderHtmlSuperscriptTmpl,
	"render/html/table.tmpl":        renderHtmlTableTmpl,
	"render/html/target.tmpl":       renderHtmlTargetTmpl,
	"render/html/toc.tmpl":          renderHtmlTocTmpl,
	"render/html/verbatim.tmpl":     renderHtmlVerbatimTmpl,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"render": &bintree{nil, map[string]*bintree{
		"html": &bintree{nil, map[string]*bintree{
			"aside.tmpl":        &bintree{renderHtmlAsideTmpl, map[string]*bintree{}},
			"bold.tmpl":         &bintree{renderHtmlBoldTmpl, map[string]*bintree{}},
			"definitions.tmpl":  &bintree{renderHtmlDefinitionsTmpl, map[string]*bintree{}},
			"image.tmpl":        &bintree{renderHtmlImageTmpl, map[string]*bintree{}},
			"inset.tmpl":        &bintree{renderHtmlInsetTmpl, map[string]*bintree{}},
			"italic.tmpl":       &bintree{renderHtmlItalicTmpl, map[string]*bintree{}},
			"larger.tmpl":       &bintree{renderHtmlLargerTmpl, map[string]*bintree{}},
			"link.tmpl":         &bintree{renderHtmlLinkTmpl, map[string]*bintree{}},
			"list.tmpl":         &bintree{renderHtmlListTmpl, map[string]*bintree{}},
			"page.tmpl":         &bintree{renderHtmlPageTmpl, map[string]*bintree{}},
			"paragraph.tmpl":    &bintree{renderHtmlParagraphTmpl, map[string]*bintree{}},
			"preformatted.tmpl": &bintree{renderHtmlPreformattedTmpl, map[string]*bintree{}},
			"reference.tmpl":    &bintree{renderHtmlReferenceTmpl, map[string]*bintree{}},
			"section.tmpl":      &bintree{renderHtmlSectionTmpl, map[string]*bintree{}},
			"sequence.tmpl":     &bintree{renderHtmlSequenceTmpl, map[string]*bintree{}},
			"smaller.tmpl":      &bintree{renderHtmlSmallerTmpl, map[string]*bintree{}},
			"strike.tmpl":       &bintree{renderHtmlStrikeTmpl, map[string]*bintree{}},
			"string.tmpl":       &bintree{renderHtmlStringTmpl, map[string]*bintree{}},
			"subscript.tmpl":    &bintree{renderHtmlSubscriptTmpl, map[string]*bintree{}},
			"superscript.tmpl":  &bintree{renderHtmlSuperscriptTmpl, map[string]*bintree{}},
			"table.tmpl":        &bintree{renderHtmlTableTmpl, map[string]*bintree{}},
			"target.tmpl":       &bintree{renderHtmlTargetTmpl, map[string]*bintree{}},
			"toc.tmpl":          &bintree{renderHtmlTocTmpl, map[string]*bintree{}},
			"verbatim.tmpl":     &bintree{renderHtmlVerbatimTmpl, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
