// Code generated by go-bindata.
// sources:
// web/templates/ab_views/html/layout.html
// web/templates/ab_views/html/register.html
// DO NOT EDIT!

package handler

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

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _webTemplatesAb_viewsHtmlLayoutHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x55\xc1\x8e\xdb\x46\x0c\xbd\xef\x57\xb0\xd3\x43\x2f\x2b\x6b\x9d\xc4\x4d\x03\x48\x06\x16\x2d\x5a\xf4\xd2\xc3\x76\x7f\x80\x1e\x51\x16\xbb\xa3\x19\x61\x48\xc9\x15\x82\x00\xfd\x88\x1e\xfb\x01\xfd\xae\x7c\x49\x31\xd2\xda\x96\xd7\xc9\xa2\x27\x8f\x38\xe4\x23\xf9\xe6\x91\x2e\xbe\xa9\x82\xd5\xb1\x23\x68\xb4\x75\xdb\x9b\x22\xfd\x80\x43\xbf\x2f\x0d\x79\xb3\xbd\x01\x28\x1a\xc2\x2a\x1d\x00\x8a\x96\x14\xc1\x36\x18\x85\xb4\x34\xbd\xd6\xd9\x0f\x66\x79\xe5\xb1\xa5\xd2\x3c\xdc\x3f\xfe\xfa\xdb\x2f\x06\x6c\xf0\x4a\x5e\x4b\xf3\xf0\x78\x9f\x6d\xee\xde\xbd\xc9\xd6\x1f\x3e\x7c\x9f\xad\xdf\xdd\xdd\x65\xeb\xcd\xfb\xf7\xd9\xc3\xe3\xfd\x17\xc2\x07\xa6\x43\x17\xa2\x2e\x00\x0e\x5c\x69\x53\x56\x34\xb0\xa5\x6c\xfa\xb8\x05\xf6\xac\x8c\x2e\x13\x8b\x8e\xca\xf5\x2d\x48\x13\xd9\x3f\x65\x1a\xb2\x9a\xb5\xf4\xe1\x0b\xd0\x15\x89\x8d\xdc\x29\x07\xbf\x40\xff\x39\x12\xc1\x81\x76\x16\x5b\x88\x64\x43\xac\x04\xea\x18\xda\xd4\xa9\xf6\x71\x87\x4a\xd0\x51\xac\x43\x6c\x29\xca\x11\x56\x59\x1d\x6d\x7f\xc7\x81\x66\x8f\x0c\xea\x6b\x9c\x22\x9f\xdd\x6e\xae\x4a\xd1\x86\x5a\xca\x6c\x70\x21\x2e\x4a\xf9\xf6\x4d\xfd\x76\xfd\x76\x63\x8e\x01\x8e\xfd\x13\x44\x72\xa5\x11\x1d\x1d\x49\x43\xa4\x06\x9a\x48\x75\x69\x72\x51\x54\xb6\x79\xc5\xa2\x79\x8b\xec\x57\x56\x52\x75\x45\x3e\xbf\x58\xb1\x0b\xd5\x38\xbd\x60\xc5\x03\x58\x87\x22\xa5\x49\x99\x90\x3d\xc5\x63\x1b\x8b\xbb\x18\x0e\xcf\xd6\x97\x31\xee\x64\xbf\xbc\x11\x56\xca\x52\xb6\x13\xde\xb3\x4f\xb3\xbe\x70\x71\x61\x1f\xcc\xb6\x90\x0e\xfd\xc9\xae\x18\xcd\xf6\xf3\x3f\xff\x16\x79\x32\x6f\xe1\xcc\xe4\xab\x7e\x45\xde\xac\x2f\x53\x6d\xae\x52\x65\xd2\xef\x26\xde\xcd\xf6\xa7\xe0\xbf\x53\x68\x59\x04\xb4\x21\x50\x14\x65\x12\x2d\xf2\x66\xb3\x68\x29\xaf\x78\x38\x75\x7e\xfe\x38\x1d\x4f\x87\x8f\x1f\x95\xda\xce\xa5\x2a\xcd\xf3\xa3\x19\x58\x7d\xfa\x94\x9e\xab\xa8\x43\x50\x8a\xd7\xbc\xbe\xe4\xfc\xeb\xac\xbf\xc6\x3b\x40\x21\x2d\x3a\xb7\xb4\xc0\x82\xb5\xcf\x7f\xfd\x0d\x8f\x0d\xc1\x8f\x67\xd1\xde\x47\xdb\xf0\x40\xab\x8b\x88\xc5\x3d\x0b\xa0\x07\xac\x7a\xa7\x49\xb7\x89\x3f\xe8\x62\x18\xb8\x62\xbf\x07\xc7\xc3\x51\xce\x17\x00\xcf\xa3\x80\xde\x92\xc0\x6e\x04\x6c\x51\xa9\x8f\x60\xb1\xdd\x73\x74\x72\x9b\x4e\xbb\x30\x26\xf0\x0a\x6c\xe8\x3b\x47\x02\x3a\x76\x6c\xd1\xb9\x11\x6a\x4a\x05\xa4\x0c\xbe\xaf\x58\xc7\xc9\x4d\xe8\xcf\x1e\x1d\xa0\x55\x1e\x58\xc7\x8b\x84\x11\xfd\x3e\xb9\x4f\x33\x29\x9a\x06\x98\x50\x68\x8a\xab\x38\xea\x08\x8a\xee\x09\x34\x40\x8b\x32\xf7\xc6\xc1\xc3\x81\xb5\x49\xb8\xa0\x61\x94\x4b\x0e\xcc\x99\x04\x33\xb1\x00\x69\xe1\xb4\xe8\x95\xb0\x87\x50\x83\x49\xa3\x6f\xa6\x0c\xe6\x04\x4a\x66\xf5\x15\xee\x8f\x5b\x63\x0c\x7d\x84\x1a\x87\x10\x13\x95\x13\x81\x27\x76\xd3\x56\xd8\xc5\x80\x95\x45\x51\xb9\x00\x6a\xf1\x29\xf5\xb7\x1b\x67\x00\x17\x06\x72\xe3\x62\xe5\xcc\x9d\x9f\x6b\x5e\xd9\xd0\xa6\x76\x0f\xa8\xb6\x01\x56\x48\x82\x8c\xcb\xe2\x8a\xfc\x85\x56\xfe\x87\xc4\x8f\xf2\x7d\x21\xf3\x3f\x70\xc0\x79\x6b\xca\x2c\xf5\x22\x9f\x17\x4b\x91\xcf\xff\x19\xff\x05\x00\x00\xff\xff\x92\x4b\xe1\xe4\x44\x06\x00\x00")

func webTemplatesAb_viewsHtmlLayoutHtmlBytes() ([]byte, error) {
	return bindataRead(
		_webTemplatesAb_viewsHtmlLayoutHtml,
		"web/templates/ab_views/html/layout.html",
	)
}

func webTemplatesAb_viewsHtmlLayoutHtml() (*asset, error) {
	bytes, err := webTemplatesAb_viewsHtmlLayoutHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "web/templates/ab_views/html/layout.html", size: 1604, mode: os.FileMode(420), modTime: time.Unix(1607595886, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _webTemplatesAb_viewsHtmlRegisterHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x95\xcf\x6e\x9c\x3e\x10\xc7\xef\x79\x8a\xd1\x1c\x7e\xfa\xf5\x00\x5c\xab\x16\xb8\x44\xbd\x46\x51\x5f\x20\x32\x30\x2c\x6e\xc1\xb6\xc6\x66\x37\x11\xf2\xbb\x57\x06\x2f\x64\x37\xdd\x2a\xd9\x6d\x4f\xb6\x87\xf9\xf3\x9d\x8f\x65\x66\x9a\x1a\x6a\xa5\x22\xc0\x5a\x2b\x47\xca\xa1\xf7\x77\x79\x23\xf7\x50\xf7\xc2\xda\x02\x59\x1f\xb0\xbc\x03\x78\x6d\xab\x75\x9f\xd8\x01\xcb\xff\x54\x65\xcd\xd7\x3c\x6b\xe4\xfe\xf7\x2e\xc9\xe7\x39\x16\x20\x6f\x35\x0f\x20\x6a\x27\xb5\x2a\x70\x9a\x06\x3d\x2a\x67\x84\xeb\xa8\x01\x64\xda\x49\xeb\x88\xd1\x7b\x84\x81\x5c\xa7\x9b\x02\x8d\xb6\x2e\x06\x9f\x66\x0e\x99\x92\x1d\xeb\xd1\xac\x9f\x01\xa6\xe9\x20\x5d\x07\x29\x31\x6b\xb6\xde\xc7\xf3\xff\x52\x35\xf4\x0c\x29\x20\x7e\x0a\x46\x16\x6a\x47\x90\x7a\x9f\x5b\x23\x54\x39\x4d\x61\x9b\xcd\xfb\xbc\x62\xc8\xca\x69\x22\xd5\x04\xcf\x6d\x81\xc4\xfb\xb5\x4e\xde\x8b\x8a\x7a\x68\x35\x17\xa8\xc4\x40\x58\x3e\x88\x81\xbe\xe4\xd9\x6c\xdf\xf4\xe4\x52\x99\xd1\x9d\x48\x0e\x7c\x59\xf7\x08\x21\x2e\x46\x83\x7b\x31\x54\xa0\xa3\x67\x87\xb0\x17\xfd\x48\x01\xce\xd2\x8a\x61\xb2\xc4\x7b\x5a\x9b\x49\x43\x48\x38\xa5\xa7\x02\x03\x35\xd3\x8b\x9a\x3a\xdd\x37\xc4\x05\x3e\xcc\xa9\xb3\x15\xde\x72\x3f\xd7\xa3\x8c\xd4\x96\xf2\xef\x21\x77\x11\x19\x0d\x42\xf6\x58\x7e\x4b\xc2\x7a\x0d\xb5\x25\xc1\xc7\xb0\xcd\x31\xef\xe1\xb6\xc8\xfa\x07\xe4\xa2\x82\x9b\xd0\x19\x61\xed\x41\x73\x83\xe5\x63\xdc\x5d\xc3\x6f\xcd\x12\x11\x6e\xe7\x13\x12\x8f\xab\xf9\xaf\xb3\x38\x56\xbc\x11\x47\xad\x55\x2b\x79\x78\xda\xb0\xdc\x2f\x16\xb8\x05\xcf\x9b\xac\x7f\xc6\x74\x5e\xf2\x22\xae\x8b\x34\xce\x0b\x7e\x98\x4a\x5e\x8d\xce\x69\x15\x75\xda\xb1\x1a\xa4\xc3\x63\x8b\x95\x53\x50\x39\x95\x18\x96\x83\xe0\x17\x2c\xbf\xc7\x3f\x6d\x9e\x2d\x61\xdb\x7d\x0a\xe8\x98\xda\x02\x33\x2c\xef\x85\xaa\xa9\xcf\x33\xf1\x46\x7e\x6d\xb9\x7d\x72\xfa\x27\x29\xef\x23\xcc\xa5\x6e\x27\x9b\x86\xd4\x4a\x71\x75\x7b\xf5\x3a\xd3\xf0\xdc\xd6\x36\x96\xa1\x90\x85\x3b\x98\x27\xc7\xe5\x11\x72\x36\x65\xe2\x72\x4c\x73\xb7\xcd\xaf\x1f\x62\x2f\x6c\xcd\xd2\x38\x1b\x66\xd8\xd1\xe3\x57\x00\x00\x00\xff\xff\x19\x79\xe6\xb0\xe1\x06\x00\x00")

func webTemplatesAb_viewsHtmlRegisterHtmlBytes() ([]byte, error) {
	return bindataRead(
		_webTemplatesAb_viewsHtmlRegisterHtml,
		"web/templates/ab_views/html/register.html",
	)
}

func webTemplatesAb_viewsHtmlRegisterHtml() (*asset, error) {
	bytes, err := webTemplatesAb_viewsHtmlRegisterHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "web/templates/ab_views/html/register.html", size: 1761, mode: os.FileMode(420), modTime: time.Unix(1607595886, 0)}
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
	"web/templates/ab_views/html/layout.html":   webTemplatesAb_viewsHtmlLayoutHtml,
	"web/templates/ab_views/html/register.html": webTemplatesAb_viewsHtmlRegisterHtml,
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
	"web": &bintree{nil, map[string]*bintree{
		"templates": &bintree{nil, map[string]*bintree{
			"ab_views": &bintree{nil, map[string]*bintree{
				"html": &bintree{nil, map[string]*bintree{
					"layout.html":   &bintree{webTemplatesAb_viewsHtmlLayoutHtml, map[string]*bintree{}},
					"register.html": &bintree{webTemplatesAb_viewsHtmlRegisterHtml, map[string]*bintree{}},
				}},
			}},
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
