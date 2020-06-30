// Copyright 2020 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package openapi3

import (
	kopenapi3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/gohugoio/hugo/cache/namedmemcache"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/resources/resource"
)

// New returns a new instance of the openapi3-namespaced template functions.
func New(deps *deps.Deps) *Namespace {
	// TODO1 consolidate when merging that "other branch"
	cache := namedmemcache.New()
	deps.BuildStartListeners.Add(
		func() {
			cache.Clear()
		})

	return &Namespace{
		cache: cache,
		deps:  deps,
	}
}

// Namespace provides template functions for the "openapi3".
type Namespace struct {
	cache *namedmemcache.Cache
	deps  *deps.Deps
}

func (ns *Namespace) Unmarshal(r resource.ContentResource) (*kopenapi3.Swagger, error) {
	c, err := r.Content()
	if err != nil {
		return nil, err
	}

	return kopenapi3.NewSwaggerLoader().LoadSwaggerFromData([]byte(c.(string)))

}
