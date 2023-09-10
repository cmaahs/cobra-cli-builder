package main

import (
	"embed"
	"fmt"
	"log"
	"new-cli-tool/version"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/maahsome/vault-view/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	//go:embed template template/cmd template/cmd/get template/common template/config template/objects
	fs embed.FS

	templateDir string
)

type ProgramInfo struct {
	Program string
}

func main() {
	program := os.Args[1]
	appVer := version.ExpressVersion()
	semVer := version.GetSemVer()
	logrus.Infof("Version: %s", appVer)
	logrus.Infof("semVer: %s", semVer)
	initHomes(semVer)
	logrus.Infof("WorkDir: %s", templateDir)
	logrus.Infof("Program: %s", program)

	err := filepath.Walk(templateDir,
		func(path string, info os.FileInfo, err error) error {
			shortPath := strings.TrimPrefix(path, fmt.Sprintf("%s", templateDir))
			shortPath = strings.TrimPrefix(shortPath, "/")
			if err != nil {
				return err
			}
			if len(shortPath) > 0 {
				if stat, err := os.Stat(path); err == nil && stat.IsDir() {
					logrus.Infof("Directory: %s", shortPath)
					if _, err := os.Stat(fmt.Sprintf("./%s", shortPath)); err != nil {
						if os.IsNotExist(err) {
							mkerr := os.MkdirAll(fmt.Sprintf("./%s", shortPath), os.ModePerm)
							if mkerr != nil {
								common.Logger.Fatal(fmt.Sprintf("Error creating %s directory", shortPath), mkerr)
							}
						}
					}
				} else {
					newFile := strings.TrimSuffix(shortPath, ".tpl")
					logrus.Infof("File: %s", newFile)
					var temp *template.Template
					temp = template.Must(template.ParseFiles(path))
					tplVariables := ProgramInfo{Program: program}

					f, ferr := os.Create(newFile)
					if ferr != nil {
						logrus.Fatal("create file: ", ferr)
					}

					err := temp.Execute(f, tplVariables)
					if err != nil {
						log.Fatalln(err)
					}
					f.Close()
				}
			}
			return nil
		})
	if err != nil {
		logrus.Error(err)
	}
}

func saveTemplateFile(dir string, filename string, workDir string) {
	data, _ := fs.ReadFile(fmt.Sprintf("%s/%s", dir, filename))
	shortDir := strings.Replace(dir, "template", "", 1)
	if strings.HasPrefix(shortDir, "/") {
		shortDir = strings.Replace(shortDir, "/", "", 1)
	}
	writeDir := fmt.Sprintf("%s/%s", workDir, shortDir)
	if _, err := os.Stat(writeDir); err != nil {
		if os.IsNotExist(err) {
			mkerr := os.MkdirAll(writeDir, os.ModePerm)
			if mkerr != nil {
				common.Logger.Fatalf("Error creating %s", writeDir, mkerr)
			}
		}
	}
	logrus.Warnf("Write file to: %s/%s", writeDir, filename)
	werr := os.WriteFile(fmt.Sprintf("%s/%s", writeDir, filename), data, os.ModePerm)
	if werr != nil {
		logrus.Errorf("Error creating file: %s/%s", writeDir, filename)
	}
}

func populateTemplates(dir string, workDir string) {

	logrus.Warnf("Template Embedded Files: ", dir)
	files, _ := fs.ReadDir(dir)

	for _, f := range files {
		if f.IsDir() {
			logrus.Warnf("%s/", f.Name())
			populateTemplates(fmt.Sprintf("%s/%s", dir, f.Name()), workDir)
		} else {
			saveTemplateFile(dir, f.Name(), workDir)
		}
	}
}

func initHomes(semVer string) {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	templateDir = fmt.Sprintf("%s/.config/new-cli-tool/templates/%s", home, semVer)
	workDir := fmt.Sprintf("%s/.config/new-cli-tool/templates/%s", home, semVer)
	logrus.Warnf("workDir: %s", workDir)
	if _, err := os.Stat(workDir); err != nil {
		if os.IsNotExist(err) {
			mkerr := os.MkdirAll(workDir, os.ModePerm)
			if mkerr != nil {
				common.Logger.Fatal("Error creating ~/.config/new-cli-tool directory", mkerr)
			}
			populateTemplates("template", workDir)
		}
	}
	workDir = fmt.Sprintf("%s/.config/new-cli-tool", home)
	if stat, err := os.Stat(workDir); err == nil && stat.IsDir() {
		configFile := fmt.Sprintf("%s/%s", workDir, "config.yaml")
		createRestrictedConfigFile(configFile)
		viper.SetConfigFile(configFile)
	} else {
		common.Logger.Info("The ~/.config/new-cli-tool path is a file and not a directory, please remove the 'new-cli-tool' file.")
		os.Exit(1)
	}
}

func createRestrictedConfigFile(fileName string) {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			file, ferr := os.Create(fileName)
			if ferr != nil {
				common.Logger.Info("Unable to create the configfile.")
				os.Exit(1)
			}
			mode := int(0600)
			if cherr := file.Chmod(os.FileMode(mode)); cherr != nil {
				common.Logger.Info("Chmod for config file failed, please set the mode to 0600.")
			}
		}
	}
}
