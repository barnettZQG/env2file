// RAINBOND, application Management Platform
// Copyright (C) 2014-2017 Goodrain Co., Ltd.

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version. For any non-GPL usage of Rainbond,
// one or multiple Commercial Licenses authorized by Goodrain Co., Ltd.
// must be obtained first.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/urfave/cli"
)

var version string

func main() {
	app := cli.NewApp()
	app.Version = version
	app.Flags = []cli.Flag{}
	app.Commands = getCommands()
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("env2config run error:%s \n", err.Error())
		os.Exit(1)
	}
}

var reg = regexp.MustCompile(`(?U)\$\{.*\}`)

func getCommands() (cmds []cli.Command) {
	cmds = append(cmds, cli.Command{
		Name:      "conversion",
		ShortName: "conv",
		Usage:     "Renders the specified file template by reading the environment variable",
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:  "file,f",
				Usage: "The configuration file that needs to be processed",
			},
		},
		Action: func(c *cli.Context) error {
			if files := c.StringSlice("file"); files != nil {
				for _, f := range files {
					if err := handleConfigFile(f); err != nil {
						return err
					}
				}
				return nil
			}
			file := c.Args().First()
			if file == "" {
				return fmt.Errorf("config file not define")
			}
			return handleConfigFile(file)
		},
	})
	cmds = append(cmds, cli.Command{
		Name:      "create",
		ShortName: "cre",
		Usage:     "Generates a configuration file to the specified path by reading the environment variable and following the specified format",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "path",
				Usage: "The path of configuration file ",
			},
			cli.IntFlag{
				Name:  "perm",
				Value: 0755,
				Usage: "The perm of configuration file ",
			},
			cli.StringFlag{
				Name:  "format,f",
				Value: "default",
				Usage: "Some common configuration file formats.eg mysql redis",
			},
		},
		Action: func(c *cli.Context) error {
			format := c.String("format")
			if format != "default" && format != "mysql" && format != "redis" {
				return fmt.Errorf("this format (%s) can not support", format)
			}
			filepath := c.String("path")
			fmt.Printf("Will write config file %s", filepath)
			if err := os.MkdirAll(path.Dir(filepath), 0755); err != nil {
				if !strings.Contains(err.Error(), "") {
					return fmt.Errorf("create config dir error %s", err.Error())
				}
			}
			//if config file exist,remove it
			os.Remove(filepath)
			perm := c.Int("perm")
			if perm < 400 || perm > 777 {
				return fmt.Errorf("set file perm is %d ivide", perm)
			}
			switch format {
			case "mysql":
				if err := createMysqlConfigFile(filepath, 000000000000|os.FileMode(perm)); err != nil {
					return err
				}
			case "redis":
				if err := createRedisConfigFile(filepath, 000000000000|os.FileMode(perm)); err != nil {
					return err
				}
			default:
				if err := createDefaultConfigFile(filepath, 000000000000|os.FileMode(perm)); err != nil {
					return err
				}
			}
			return nil
		},
	})
	return
}

//GetConfigKey get really key
func GetConfigKey(rk string) string {
	if len(rk) < 4 {
		return ""
	}
	left := strings.Index(rk, "{")
	right := strings.Index(rk, "}")
	return rk[left+1 : right]
}

func handleConfigFile(file string) error {
	f, err := os.Stat(file)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	bstr := string(body)
	resultKey := reg.FindAllString(bstr, -1)
	for _, rk := range resultKey {
		value := os.Getenv(GetConfigKey(rk))
		fmt.Printf("key %s value %s \n", GetConfigKey(rk), value)
		bstr = strings.Replace(bstr, rk, value, -1)
	}
	if len(resultKey) == 0 {
		fmt.Println("")
	}
	return ioutil.WriteFile(file, []byte(bstr), f.Mode())
}
func getEnvs() map[string]string {
	envs := os.Environ()
	envMap := make(map[string]string, len(envs))
	for _, env := range envs {
		index := strings.Index(env, "=")
		if index != -1 {
			key := env[:index]
			value := env[index+1:]
			envMap[key] = value
		}
	}
	return envMap
}

//createMysqlConfigFile create mysql config file by env
//eg. MYSQLC_MYSQLD_PORT => [mysqld] port
//eg. MYSQLC_MYSQLD_DATADIR => [mysqld] datadir
func createMysqlConfigFile(file string, perm os.FileMode) error {
	envs := getEnvs()
	var configs = make(map[string]map[string]string)
	for k, v := range envs {
		k1 := strings.ToLower(k)
		if strings.HasPrefix(k1, "mysqlc_") {
			k2 := k1[7:]
			index := strings.Index(k2, "_")
			if index != -1 {
				module := k2[:index]
				if module != "" {
					k3 := k2[index+1:]
					if k3 != "" {
						if c, ok := configs[module]; ok {
							c[k3] = handleConvValue(v)
						} else {
							configs[module] = map[string]string{k3: handleConvValue(v)}
						}
					}
				}
			}
		}
	}
	writer := bytes.NewBuffer(nil)
	for k, v := range configs {
		writer.WriteString(fmt.Sprintf("[%s]\n", k))
		for k2, v2 := range v {
			writer.WriteString(fmt.Sprintf("%s = %s\n", k2, v2))
		}
		writer.WriteString("\n")
	}
	ioutil.WriteFile(file, writer.Bytes(), perm)
	return nil
}

//createRedisConfigFile create redis config file by env
//eg. REDISC_PORT=6379 => port 6379
func createRedisConfigFile(file string, perm os.FileMode) error {
	envs := getEnvs()
	var configs = make(map[string]string)
	for k, v := range envs {
		k1 := strings.ToLower(k)
		if strings.HasPrefix(k1, "redisc_") {
			k2 := k1[7:]
			if k2 != "" {
				configs[k2] = handleConvValue(v)
			}
		}
	}
	writer := bytes.NewBuffer(nil)
	for k, v := range configs {
		writer.WriteString(fmt.Sprintf("%s %s \n", k, v))
	}
	ioutil.WriteFile(file, writer.Bytes(), perm)
	return nil
}

//createDefaultConfigFile create default format config
//eg. C_MYNAME=barnett >  myname = barnett
func createDefaultConfigFile(file string, perm os.FileMode) error {
	envs := getEnvs()
	var configs = make(map[string]string)
	for k, v := range envs {
		k1 := strings.ToLower(k)
		if strings.HasPrefix(k1, "c_") {
			k2 := k1[7:]
			if k2 != "" {
				configs[k2] = handleConvValue(v)
			}
		}
	}
	writer := bytes.NewBuffer(nil)
	for k, v := range configs {
		writer.WriteString(fmt.Sprintf("%s=%s \n", k, v))
	}
	ioutil.WriteFile(file, writer.Bytes(), perm)
	return nil
}

//handleConvValue conv env value.
func handleConvValue(source string) string {
	resultKey := reg.FindAllString(source, -1)
	for _, rk := range resultKey {
		value := os.Getenv(GetConfigKey(rk))
		if value != "" {
			source = strings.Replace(source, rk, value, -1)
		}
	}
	return source
}
