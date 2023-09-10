package new

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"new-cli-subcmd/common"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	//go:embed template
	fs embed.FS

	templateDir string
)

type ProgramInfo struct {
	Program      string
	Cmd          string
	Subcmd       string
	SubcmdPublic string
}

// subcmdCmd represents the version command
var subcmdCmd = &cobra.Command{
	Use:   "subcmd",
	Short: "Add a sub command to an exiting project",
	Run: func(cmd *cobra.Command, args []string) {
		program := os.Args[3]
		command := os.Args[4]
		subcommand := os.Args[5]
		if len(program) == 0 || len(command) == 0 || len(subcommand) == 0 {
			logrus.Fatal("Must provide 'program' 'command' 'subcommand' positional parameters")
		}
		initHomes(c.VersionDetail.SemVer)
		createSubCommand(program, command, subcommand)
	},
}

func createSubCommand(program string, command string, subcommand string) {
	logrus.Infof("WorkDir: %s", templateDir)
	logrus.Infof("Program: %s", program)
	logrus.Infof("Command: %s", command)
	logrus.Infof("SubCommand: %s", subcommand)

	// Setup the directories
	cmdDir := fmt.Sprintf("cmd/%s", command)
	logrus.Infof("Directory: %s", cmdDir)
	if _, err := os.Stat(fmt.Sprintf("./%s", cmdDir)); err != nil {
		if os.IsNotExist(err) {
			mkerr := os.MkdirAll(fmt.Sprintf("./%s", cmdDir), os.ModePerm)
			if mkerr != nil {
				common.Logger.Fatal(fmt.Sprintf("Error creating %s directory", cmdDir), mkerr)
			}
		}
	}

	err := filepath.Walk(templateDir,
		func(path string, info os.FileInfo, err error) error {
			shortPath := strings.TrimPrefix(path, templateDir)
			shortPath = strings.TrimPrefix(shortPath, "/")
			if err != nil {
				return err
			}
			if len(shortPath) > 0 {
				newFile := strings.TrimSuffix(shortPath, ".tpl")
				newFile = strings.Replace(newFile, "cmd", command, 1)
				newFile = strings.Replace(newFile, "subcmd", subcommand, 1)
				newFile = fmt.Sprintf("cmd/%s/%s", command, newFile)
				logrus.Infof("File: %s", newFile)
				temp := template.Must(template.ParseFiles(path))
				tplVariables := ProgramInfo{
					Program:      program,
					Cmd:          command,
					Subcmd:       subcommand,
					SubcmdPublic: cases.Title(language.English, cases.NoLower).String(subcommand),
				}

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
			return nil
		})
	if err != nil {
		logrus.Error(err)
	}
	fmt.Printf("Add '%s/cmd/%s' to root.go import section\n", program, command)
	fmt.Printf("Add '%s.InitSubCommands(c)' to RootCmd.AddCommand()\n", command)

}

func init() {
	newCmd.AddCommand(subcmdCmd)
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

	templateDir = fmt.Sprintf("%s/.config/new-cli-subcmd/templates/%s", home, semVer)
	workDir := fmt.Sprintf("%s/.config/new-cli-subcmd/templates/%s", home, semVer)
	logrus.Warnf("workDir: %s", workDir)
	if _, err := os.Stat(workDir); err != nil {
		if os.IsNotExist(err) {
			mkerr := os.MkdirAll(workDir, os.ModePerm)
			if mkerr != nil {
				common.Logger.Fatal("Error creating ~/.config/new-cli-subcmd directory", mkerr)
			}
			populateTemplates("template", workDir)
		}
	}
	workDir = fmt.Sprintf("%s/.config/new-cli-subcmd", home)
	if stat, err := os.Stat(workDir); err == nil && stat.IsDir() {
		configFile := fmt.Sprintf("%s/%s", workDir, "config.yaml")
		createRestrictedConfigFile(configFile)
		viper.SetConfigFile(configFile)
	} else {
		common.Logger.Info("The ~/.config/new-cli-subcmd path is a file and not a directory, please remove the 'new-cli-subcmd' file.")
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
