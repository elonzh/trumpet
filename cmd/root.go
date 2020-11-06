/*
Copyright © 2020 elonzh

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
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfg = &Config{
		LogLevel: logrus.InfoLevel.String(),
	}
)

func newRootCmd(version string) *cobra.Command {
	var cfgFile string
	var rootCmd = &cobra.Command{
		Version: version,
		Use:     "trumpet",
		Short:   "🎺simple webhook transform server",
		Long:    ``,
	}
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ./config.yaml)")
	rootCmd.PersistentFlags().String("logLevel", "info", "")
	err := viper.BindPFlag("logLevel", rootCmd.PersistentFlags().Lookup("logLevel"))
	if err != nil {
		panic(err)
	}

	cobra.OnInitialize(func() {
		initConfig(cfgFile)
	})
	rootCmd.AddCommand(serveCmd)
	return rootCmd
}

func Execute(version string) {
	rootCmd := newRootCmd(version)
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatalln()
	}
}

func initConfig(cfgFile string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}
	var err error
	if err = viper.ReadInConfig(); os.IsNotExist(err) {
		logrus.WithError(err).Fatalln()
	}
	logrus.WithField("ConfigFile", viper.ConfigFileUsed()).Infoln("read in config")
	err = viper.Unmarshal(cfg)
	if err != nil {
		logrus.WithError(err).Fatalln("error when unmarshal config")
	}
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logrus.WithError(err).Fatalln()
	}
	logrus.SetLevel(level)
	if level >= logrus.DebugLevel {
		logrus.WithField("Config", cfg).Debug()
	}

	cfg.LoadAllTransformers()
}
