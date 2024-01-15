// Copyright 2024 Alexey Lavrenchenko. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package views

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"unsafe"
)

type view struct {
	content []byte
	tmpl    *template.Template
}

type ViewManager struct {
	dir   string
	views map[string]*view
	mu    sync.RWMutex
}

func NewViewManager(dir string) *ViewManager {
	return &ViewManager{
		dir:   filepath.Clean(dir),
		views: map[string]*view{},
	}
}

// Render renders a view.
func (m *ViewManager) Render(w http.ResponseWriter, viewName string, viewData any) error {
	v, err := m.get(viewName)
	if err != nil {
		return fmt.Errorf("[views.ViewManager.Render] get a view: %w", err)
	}

	// w.WriteHeader(http.StatusOK)
	// val := reflect.ValueOf(viewData)
	// if viewData != nil && (val.Kind() != reflect.Ptr || !val.IsNil()) {
	// 	if err = v.tmpl.Execute(w, viewData); err != nil {
	// 		return fmt.Errorf("[views.ViewManager.Render] execute a view template: %w", err)
	// 	}
	// } else if _, err = w.Write(v.content); err != nil {
	// 	return fmt.Errorf("[views.ViewManager.Render] write the view content: %w", err)
	// }

	var b []byte
	val := reflect.ValueOf(viewData)
	if viewData != nil && (val.Kind() != reflect.Ptr || !val.IsNil()) {
		buf := new(bytes.Buffer)
		buf.Grow(len(v.content))

		if err = v.tmpl.Execute(buf, viewData); err != nil {
			return fmt.Errorf("[views.ViewManager.Render] execute a view template: %w", err)
		}
		b = buf.Bytes()
	} else {
		b = v.content
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(b); err != nil {
		return fmt.Errorf("[views.ViewManager.Render] write data: %w", err)
	}
	return nil
}

func (m *ViewManager) get(name string) (*view, error) {
	m.mu.RLock()
	v := m.views[name]
	m.mu.RUnlock()
	if v != nil {
		return v, nil
	}

	m.mu.Lock()
	if v = m.views[name]; v != nil {
		m.mu.Unlock()
	} else {
		defer m.mu.Unlock()
		var err error
		if v, err = m.load(name); err != nil {
			return nil, fmt.Errorf("[views.ViewManager.get] load a view: %w", err)
		}

		m.views[name] = v
	}
	return v, nil
}

func (m *ViewManager) load(name string) (*view, error) {
	c, err := os.ReadFile(filepath.Join(m.dir, name))
	if err != nil {
		return nil, fmt.Errorf("[views.ViewManager.load] read a file: %w", err)
	}

	t, err := template.New(name).Parse(unsafe.String(unsafe.SliceData(c), len(c)))
	if err != nil {
		return nil, fmt.Errorf("[views.ViewManager.load] parse the view content: %w", err)
	}

	return &view{
		content: c,
		tmpl:    t,
	}, nil
}
