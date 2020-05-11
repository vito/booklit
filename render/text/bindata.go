// Code generated for package text by go-bindata DO NOT EDIT. (@generated)
// sources:
// render/text/inset.tmpl
// render/text/page.tmpl
// render/text/paragraph.tmpl
// render/text/preformatted.tmpl
// render/text/section.tmpl
// render/text/sequence.tmpl
// render/text/string.tmpl
// render/text/toc.tmpl
package text

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

var _renderTextInsetTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x52\x50\x50\x50\xa8\xae\xd6\x73\xce\xcf\x2b\x49\xcd\x2b\x51\xa8\x51\x28\x4a\xcd\x4b\x49\x2d\x52\xa8\x51\xc8\xca\xcf\xcc\xf3\xc9\xcc\x4b\x2d\x56\x50\x02\x29\x52\xaa\xad\xe5\x02\x04\x00\x00\xff\xff\xa4\xfc\x3d\x4d\x2d\x00\x00\x00")

func renderTextInsetTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderTextInsetTmpl,
		"render/text/inset.tmpl",
	)
}

func renderTextInsetTmpl() (*asset, error) {
	bytes, err := renderTextInsetTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/text/inset.tmpl", size: 45, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderTextPageTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xaa\xae\xd6\x0b\xc9\x2c\xc9\x49\xd5\x0b\x2e\x29\xca\xcc\x4b\xaf\xad\xe5\xe2\xaa\xae\xd6\x53\xa8\x51\x28\x4a\xcd\x4b\x49\x2d\xaa\xad\x05\x04\x00\x00\xff\xff\xa4\xfc\x1e\xc9\x21\x00\x00\x00")

func renderTextPageTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderTextPageTmpl,
		"render/text/page.tmpl",
	)
}

func renderTextPageTmpl() (*asset, error) {
	bytes, err := renderTextPageTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/text/page.tmpl", size: 33, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderTextParagraphTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xaa\xae\x2e\x4a\xcc\x4b\x4f\x55\x50\xc9\xcc\x4b\x49\xad\xd0\x51\x50\xc9\xc9\xcc\x4b\x55\xb0\xb2\x55\xd0\xab\xad\xad\xae\xce\x4c\x83\x4a\xd4\xd6\x2a\x54\x57\xa7\xe6\xa5\x80\x04\x21\x4a\x6a\x14\x8a\x52\xf3\x52\x52\x8b\x40\x22\x60\x09\x2e\xae\xea\x6a\x25\xa5\xda\x5a\x2e\x40\x00\x00\x00\xff\xff\xc7\xb4\x03\x1c\x53\x00\x00\x00")

func renderTextParagraphTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderTextParagraphTmpl,
		"render/text/paragraph.tmpl",
	)
}

func renderTextParagraphTmpl() (*asset, error) {
	bytes, err := renderTextParagraphTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/text/paragraph.tmpl", size: 83, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderTextPreformattedTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xaa\xae\x2e\x4a\xcc\x4b\x4f\x55\x50\xc9\xcc\x4b\x49\xad\xd0\x51\x50\xc9\xc9\xcc\x4b\x55\xb0\xb2\x55\xd0\xab\xad\xad\xae\xce\x4c\x83\x4a\xd4\xd6\x72\x55\x57\xa7\xe6\xa5\x80\x04\x21\x4a\x6a\x14\x8a\x52\xf3\x52\x52\x8b\x40\x22\x60\x09\x2e\x40\x00\x00\x00\xff\xff\x41\x45\x7e\xd9\x4b\x00\x00\x00")

func renderTextPreformattedTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderTextPreformattedTmpl,
		"render/text/preformatted.tmpl",
	)
}

func renderTextPreformattedTmpl() (*asset, error) {
	bytes, err := renderTextPreformattedTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/text/preformatted.tmpl", size: 75, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderTextSectionTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x4c\xcd\xc1\x09\x02\x41\x0c\x85\xe1\xfb\x54\x91\x06\x4c\x11\x7a\x17\x61\x6d\x60\x75\xb2\x1a\x18\x33\x12\xe3\x41\x9e\xe9\x5d\x58\x65\xf0\xf4\xe0\x83\xc7\x0f\xe8\x42\x7c\x98\x5d\x2c\x32\x01\xde\x3f\x6f\x27\xf1\x4c\x02\xc4\xea\x4a\x47\x8d\x26\xf4\x26\x17\xab\xe2\xb4\xc9\x2c\x05\xe0\x6d\xaf\xaf\xa1\x5f\xd3\x85\xac\x07\xf1\x74\x6f\x1a\x93\x9c\x43\xbb\x3d\x32\x0b\xe0\xb3\x5d\x84\x78\x77\xd5\x56\x5d\x6c\x35\xfe\x7f\xff\x72\x63\x3f\x01\x00\x00\xff\xff\x34\xbb\x8c\x55\x99\x00\x00\x00")

func renderTextSectionTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderTextSectionTmpl,
		"render/text/section.tmpl",
	)
}

func renderTextSectionTmpl() (*asset, error) {
	bytes, err := renderTextSectionTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/text/section.tmpl", size: 153, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderTextSequenceTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xaa\xae\x2e\x4a\xcc\x4b\x4f\x55\xd0\x73\xce\xcf\x2b\x49\xcd\x2b\x29\xae\xad\xad\xae\xd6\x53\xa8\x51\x28\x4a\xcd\x4b\x49\x2d\x02\xf1\x52\xf3\x52\x6a\x6b\xb9\x00\x01\x00\x00\xff\xff\xfc\xbe\x50\x06\x29\x00\x00\x00")

func renderTextSequenceTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderTextSequenceTmpl,
		"render/text/sequence.tmpl",
	)
}

func renderTextSequenceTmpl() (*asset, error) {
	bytes, err := renderTextSequenceTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/text/sequence.tmpl", size: 41, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderTextStringTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xaa\xae\xd6\x0b\x2e\x29\xca\xcc\x4b\xaf\xad\xe5\x02\x04\x00\x00\xff\xff\x38\x6b\xcd\x33\x0c\x00\x00\x00")

func renderTextStringTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderTextStringTmpl,
		"render/text/string.tmpl",
	)
}

func renderTextStringTmpl() (*asset, error) {
	bytes, err := renderTextStringTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/text/string.tmpl", size: 12, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _renderTextTocTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x4c\xcc\xc1\x0a\xc2\x30\x10\x04\xd0\x7b\xbf\x62\xe8\x49\x2f\xf9\x07\x29\x78\xb4\x97\xfe\x40\x6a\xb6\x1a\x48\x36\x61\xbb\x05\x61\x9b\x7f\x17\x0b\xa2\xa7\x81\x37\xcc\x98\xc5\x05\x9e\x03\xdc\xf0\x8c\x29\x08\x31\x4e\x5c\x14\x6e\xcc\x51\xbf\x74\x95\x92\x27\x3f\x27\x1a\x97\xa1\xb0\x12\xeb\x7a\x6e\xad\x33\x13\xcf\x0f\xfa\x4d\x0f\x73\xb7\x2d\xcf\x24\xad\xc1\xcc\x4d\x51\x13\x61\xc7\xaa\x12\xeb\x65\x7b\x61\x87\x10\x87\x4f\xdd\x75\x66\x4a\xb9\x26\xaf\x84\x5e\xcb\xdd\x69\xae\xa9\x87\x3b\x5e\x88\xc3\x5f\xbe\x03\x00\x00\xff\xff\xb5\x34\x2d\x4c\xa6\x00\x00\x00")

func renderTextTocTmplBytes() ([]byte, error) {
	return bindataRead(
		_renderTextTocTmpl,
		"render/text/toc.tmpl",
	)
}

func renderTextTocTmpl() (*asset, error) {
	bytes, err := renderTextTocTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "render/text/toc.tmpl", size: 166, mode: os.FileMode(420), modTime: time.Unix(1589116514, 0)}
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
	"render/text/inset.tmpl":        renderTextInsetTmpl,
	"render/text/page.tmpl":         renderTextPageTmpl,
	"render/text/paragraph.tmpl":    renderTextParagraphTmpl,
	"render/text/preformatted.tmpl": renderTextPreformattedTmpl,
	"render/text/section.tmpl":      renderTextSectionTmpl,
	"render/text/sequence.tmpl":     renderTextSequenceTmpl,
	"render/text/string.tmpl":       renderTextStringTmpl,
	"render/text/toc.tmpl":          renderTextTocTmpl,
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
		"text": &bintree{nil, map[string]*bintree{
			"inset.tmpl":        &bintree{renderTextInsetTmpl, map[string]*bintree{}},
			"page.tmpl":         &bintree{renderTextPageTmpl, map[string]*bintree{}},
			"paragraph.tmpl":    &bintree{renderTextParagraphTmpl, map[string]*bintree{}},
			"preformatted.tmpl": &bintree{renderTextPreformattedTmpl, map[string]*bintree{}},
			"section.tmpl":      &bintree{renderTextSectionTmpl, map[string]*bintree{}},
			"sequence.tmpl":     &bintree{renderTextSequenceTmpl, map[string]*bintree{}},
			"string.tmpl":       &bintree{renderTextStringTmpl, map[string]*bintree{}},
			"toc.tmpl":          &bintree{renderTextTocTmpl, map[string]*bintree{}},
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
