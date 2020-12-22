package app

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// bindata_read reads the given file from disk. It returns an error on failure.
func bindata_read(path, name string) ([]byte, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset %s at %s: %v", name, path, err)
	}
	return buf, err
}

// assets_element_ui_lib_index_js reads file data from disk. It returns an error on failure.
func assets_element_ui_lib_index_js() ([]byte, error) {
	return bindata_read(
		"/root/go/src/github.com/victoryang/kubernetes-cicd/static/assets/element-ui/lib/index.js",
		"assets/element-ui/lib/index.js",
	)
}

// assets_element_ui_lib_theme_default_fonts_element_icons_ttf reads file data from disk. It returns an error on failure.
func assets_element_ui_lib_theme_default_fonts_element_icons_ttf() ([]byte, error) {
	return bindata_read(
		"/root/go/src/github.com/victoryang/kubernetes-cicd/static/assets/element-ui/lib/theme-default/fonts/element-icons.ttf",
		"assets/element-ui/lib/theme-default/fonts/element-icons.ttf",
	)
}

// assets_element_ui_lib_theme_default_fonts_element_icons_woff reads file data from disk. It returns an error on failure.
func assets_element_ui_lib_theme_default_fonts_element_icons_woff() ([]byte, error) {
	return bindata_read(
		"/root/go/src/github.com/victoryang/kubernetes-cicd/static/assets/element-ui/lib/theme-default/fonts/element-icons.woff",
		"assets/element-ui/lib/theme-default/fonts/element-icons.woff",
	)
}

// assets_element_ui_lib_theme_default_index_css reads file data from disk. It returns an error on failure.
func assets_element_ui_lib_theme_default_index_css() ([]byte, error) {
	return bindata_read(
		"/root/go/src/github.com/victoryang/kubernetes-cicd/static/assets/element-ui/lib/theme-default/index.css",
		"assets/element-ui/lib/theme-default/index.css",
	)
}

// assets_lockr_js reads file data from disk. It returns an error on failure.
func assets_lockr_js() ([]byte, error) {
	return bindata_read(
		"/root/go/src/github.com/victoryang/kubernetes-cicd/static/assets/lockr.js",
		"assets/lockr.js",
	)
}

// assets_login_background1_jpg reads file data from disk. It returns an error on failure.
func assets_login_background1_jpg() ([]byte, error) {
	return bindata_read(
		"/root/go/src/github.com/victoryang/kubernetes-cicd/static/assets/login-background1.jpg",
		"assets/login-background1.jpg",
	)
}

// assets_login_background2_jpg reads file data from disk. It returns an error on failure.
func assets_login_background2_jpg() ([]byte, error) {
	return bindata_read(
		"/root/go/src/github.com/victoryang/kubernetes-cicd/static/assets/login-background2.jpg",
		"assets/login-background2.jpg",
	)
}

// assets_logo_png reads file data from disk. It returns an error on failure.
func assets_logo_png() ([]byte, error) {
	return bindata_read(
		"/root/go/src/github.com/victoryang/kubernetes-cicd/static/assets/logo.png",
		"assets/logo.png",
	)
}

// assets_vue_resource_js reads file data from disk. It returns an error on failure.
func assets_vue_resource_js() ([]byte, error) {
	return bindata_read(
		"/root/go/src/github.com/victoryang/kubernetes-cicd/static/assets/vue-resource.js",
		"assets/vue-resource.js",
	)
}

// assets_vue_router_min_js reads file data from disk. It returns an error on failure.
func assets_vue_router_min_js() ([]byte, error) {
	return bindata_read(
		"/root/go/src/github.com/victoryang/kubernetes-cicd/static/assets/vue-router.min.js",
		"assets/vue-router.min.js",
	)
}

// assets_vue_min_js reads file data from disk. It returns an error on failure.
func assets_vue_min_js() ([]byte, error) {
	return bindata_read(
		"/root/go/src/github.com/victoryang/kubernetes-cicd/static/assets/vue.min.js",
		"assets/vue.min.js",
	)
}

// assets_vuex_js reads file data from disk. It returns an error on failure.
func assets_vuex_js() ([]byte, error) {
	return bindata_read(
		"/root/go/src/github.com/victoryang/kubernetes-cicd/static/assets/vuex.js",
		"assets/vuex.js",
	)
}

// config_project_html reads file data from disk. It returns an error on failure.
func config_project_html() ([]byte, error) {
	return bindata_read(
		"/root/go/src/github.com/victoryang/kubernetes-cicd/static/config-project.html",
		"config-project.html",
	)
}

// css_rolling_css reads file data from disk. It returns an error on failure.
func css_rolling_css() ([]byte, error) {
	return bindata_read(
		"/root/go/src/github.com/victoryang/kubernetes-cicd/static/css/rolling.css",
		"css/rolling.css",
	)
}

// deploy_html reads file data from disk. It returns an error on failure.
func deploy_html() ([]byte, error) {
	return bindata_read(
		"/root/go/src/github.com/victoryang/kubernetes-cicd/static/deploy.html",
		"deploy.html",
	)
}

// home_html reads file data from disk. It returns an error on failure.
func home_html() ([]byte, error) {
	return bindata_read(
		"/root/go/src/github.com/victoryang/kubernetes-cicd/static/home.html",
		"home.html",
	)
}

// index_html reads file data from disk. It returns an error on failure.
func index_html() ([]byte, error) {
	return bindata_read(
		"/root/go/src/github.com/victoryang/kubernetes-cicd/static/index.html",
		"index.html",
	)
}

// login_html reads file data from disk. It returns an error on failure.
func login_html() ([]byte, error) {
	return bindata_read(
		"/root/go/src/github.com/victoryang/kubernetes-cicd/static/login.html",
		"login.html",
	)
}

// main_js reads file data from disk. It returns an error on failure.
func main_js() ([]byte, error) {
	return bindata_read(
		"/root/go/src/github.com/victoryang/kubernetes-cicd/static/main.js",
		"main.js",
	)
}

// project_html reads file data from disk. It returns an error on failure.
func project_html() ([]byte, error) {
	return bindata_read(
		"/root/go/src/github.com/victoryang/kubernetes-cicd/static/project.html",
		"project.html",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
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
var _bindata = map[string]func() ([]byte, error){
	"assets/element-ui/lib/index.js": assets_element_ui_lib_index_js,
	"assets/element-ui/lib/theme-default/fonts/element-icons.ttf": assets_element_ui_lib_theme_default_fonts_element_icons_ttf,
	"assets/element-ui/lib/theme-default/fonts/element-icons.woff": assets_element_ui_lib_theme_default_fonts_element_icons_woff,
	"assets/element-ui/lib/theme-default/index.css": assets_element_ui_lib_theme_default_index_css,
	"assets/lockr.js": assets_lockr_js,
	"assets/login-background1.jpg": assets_login_background1_jpg,
	"assets/login-background2.jpg": assets_login_background2_jpg,
	"assets/logo.png": assets_logo_png,
	"assets/vue-resource.js": assets_vue_resource_js,
	"assets/vue-router.min.js": assets_vue_router_min_js,
	"assets/vue.min.js": assets_vue_min_js,
	"assets/vuex.js": assets_vuex_js,
	"config-project.html": config_project_html,
	"css/rolling.css": css_rolling_css,
	"deploy.html": deploy_html,
	"home.html": home_html,
	"index.html": index_html,
	"login.html": login_html,
	"main.js": main_js,
	"project.html": project_html,
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
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() ([]byte, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"assets": &_bintree_t{nil, map[string]*_bintree_t{
		"element-ui": &_bintree_t{nil, map[string]*_bintree_t{
			"lib": &_bintree_t{nil, map[string]*_bintree_t{
				"index.js": &_bintree_t{assets_element_ui_lib_index_js, map[string]*_bintree_t{
				}},
				"theme-default": &_bintree_t{nil, map[string]*_bintree_t{
					"fonts": &_bintree_t{nil, map[string]*_bintree_t{
						"element-icons.ttf": &_bintree_t{assets_element_ui_lib_theme_default_fonts_element_icons_ttf, map[string]*_bintree_t{
						}},
						"element-icons.woff": &_bintree_t{assets_element_ui_lib_theme_default_fonts_element_icons_woff, map[string]*_bintree_t{
						}},
					}},
					"index.css": &_bintree_t{assets_element_ui_lib_theme_default_index_css, map[string]*_bintree_t{
					}},
				}},
			}},
		}},
		"lockr.js": &_bintree_t{assets_lockr_js, map[string]*_bintree_t{
		}},
		"login-background1.jpg": &_bintree_t{assets_login_background1_jpg, map[string]*_bintree_t{
		}},
		"login-background2.jpg": &_bintree_t{assets_login_background2_jpg, map[string]*_bintree_t{
		}},
		"logo.png": &_bintree_t{assets_logo_png, map[string]*_bintree_t{
		}},
		"vue-resource.js": &_bintree_t{assets_vue_resource_js, map[string]*_bintree_t{
		}},
		"vue-router.min.js": &_bintree_t{assets_vue_router_min_js, map[string]*_bintree_t{
		}},
		"vue.min.js": &_bintree_t{assets_vue_min_js, map[string]*_bintree_t{
		}},
		"vuex.js": &_bintree_t{assets_vuex_js, map[string]*_bintree_t{
		}},
	}},
	"config-project.html": &_bintree_t{config_project_html, map[string]*_bintree_t{
	}},
	"css": &_bintree_t{nil, map[string]*_bintree_t{
		"rolling.css": &_bintree_t{css_rolling_css, map[string]*_bintree_t{
		}},
	}},
	"deploy.html": &_bintree_t{deploy_html, map[string]*_bintree_t{
	}},
	"home.html": &_bintree_t{home_html, map[string]*_bintree_t{
	}},
	"index.html": &_bintree_t{index_html, map[string]*_bintree_t{
	}},
	"login.html": &_bintree_t{login_html, map[string]*_bintree_t{
	}},
	"main.js": &_bintree_t{main_js, map[string]*_bintree_t{
	}},
	"project.html": &_bintree_t{project_html, map[string]*_bintree_t{
	}},
}}
