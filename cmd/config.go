/*
Copyright © 2020 elonzh <elonzh@qq.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"

	"github.com/elonzh/trumpet/transformers"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type Config struct {
	LogLevel        string
	TransformersDir string
	Transformers    []*transformers.Transformer

	allTransformers []*transformers.Transformer
	m               map[string]*transformers.Transformer
}

func (c *Config) LoadAllTransformers() {
	allTransformers := make([]*transformers.Transformer, 0, len(BuiltinTransformers))
	allTransformers = append(allTransformers, BuiltinTransformers...)
	allTransformers = append(allTransformers, c.Transformers...)
	if c.TransformersDir != "" {
		trans, err := transformers.Load(c.TransformersDir)
		if err != nil {
			logrus.WithError(err).Fatalln()
		}
		allTransformers = append(allTransformers, trans...)
	}
	c.allTransformers = allTransformers

	m := make(map[string]*transformers.Transformer, len(c.Transformers))
	for _, t := range allTransformers {
		if err := t.InitThread(); err != nil {
			logrus.WithError(err).Fatalln()
		}
		if _, exists := m[t.Slug]; exists {
			logrus.WithField("Slug", t.Slug).Warnln("Transformer already exists")
		}
		m[t.Slug] = t
	}
	c.m = m
}

func (c *Config) GetTransformer(slug string) (*transformers.Transformer, bool) {
	t, ok := c.m[slug]
	return t, ok
}

func (c *Config) GetAllTransformers() []*transformers.Transformer {
	return c.allTransformers
}

func newConfigCmd(cfg *Config) *cobra.Command {
	configCmd := &cobra.Command{
		Use: "config",
	}
	showCmd := &cobra.Command{
		Use: "show",
		RunE: func(cmd *cobra.Command, args []string) error {
			bs, err := yaml.Marshal(cfg)
			if err != nil {
				return err
			}
			fmt.Println(string(bs))
			return nil
		},
	}
	configCmd.AddCommand(showCmd)
	return configCmd
}
