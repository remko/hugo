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

package js

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/gohugoio/hugo/hugolib/filesystems"
	"github.com/gohugoio/hugo/resources/internal"

	"github.com/mitchellh/mapstructure"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
)

type Options struct {
	Minify    bool
	Externals []string
	Target    string
	Defines   map[string]string
}

func DecodeOptions(m map[string]interface{}) (opts Options, err error) {
	if m == nil {
		return
	}
	err = mapstructure.WeakDecode(m, &opts)
	return
}

type Client struct {
	rs  *resources.Spec
	sfs *filesystems.SourceFilesystem
}

func New(fs *filesystems.SourceFilesystem, rs *resources.Spec) *Client {
	return &Client{rs: rs, sfs: fs}
}

type jsTransformation struct {
	options Options
	rs      *resources.Spec
	sfs     *filesystems.SourceFilesystem
}

func (t *jsTransformation) Key() internal.ResourceTransformationKey {
	return internal.NewResourceTransformationKey("js", t.options)
}

func (t *jsTransformation) Transform(ctx *resources.ResourceTransformationCtx) error {
	var target api.Target
	switch t.options.Target {
	case "", "esnext":
		target = api.ESNext
	case "es6", "es2015":
		target = api.ES2015
	case "es2016":
		target = api.ES2016
	case "es2017":
		target = api.ES2017
	case "es2018":
		target = api.ES2018
	case "es2019":
		target = api.ES2019
	case "es2020":
		target = api.ES2020
	default:
		return fmt.Errorf("invalid target: %q", t.options.Target)
	}

	src, err := ioutil.ReadAll(ctx.From)
	if err != nil {
		return err
	}

	sdir, sfile := path.Split(ctx.SourcePath)
	sdir = t.sfs.RealFilename(sdir)

	buildOptions := api.BuildOptions{
		Outfile: "",
		Bundle:  true,

		Target: target,

		MinifyWhitespace:  t.options.Minify,
		MinifyIdentifiers: t.options.Minify,
		MinifySyntax:      t.options.Minify,

		Defines: t.options.Defines,

		Externals: t.options.Externals,

		Stdin: &api.StdinOptions{
			Contents:   string(src),
			Sourcefile: sfile,
			ResolveDir: sdir,
		},
	}
	result := api.Build(buildOptions)
	if len(result.Errors) > 0 {
		return fmt.Errorf("%s", result.Errors[0].Text)
	}
	if len(result.OutputFiles) != 1 {
		return fmt.Errorf("unexpected output count: %d", len(result.OutputFiles))
	}

	ctx.To.Write(result.OutputFiles[0].Contents)
	return nil
}

func (c *Client) Process(res resources.ResourceTransformer, options Options) (resource.Resource, error) {
	return res.Transform(
		&jsTransformation{rs: c.rs, sfs: c.sfs, options: options},
	)
}
