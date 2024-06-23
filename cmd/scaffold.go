package cmd

import (
	"bytes"
	"context"
	"fmt"
	"go/format"
	"os"
	"os/exec"
	"path/filepath"
	"vimana/cmd/utils"
	"vimana/config"
	"vimana/log"
	"vimana/scaffold"

	"github.com/spf13/cobra"
	// "gorm.io/gorm/logger"
)

type Project struct {
	ProjectName string `survey:"project_name"`
}

var (
	repoURL string
)
var ScaffoldNew = &cobra.Command{
	Use:   "scaffold",
	Short: "Scaffold a new Spacecore",
	Run:   run,
}

func NewProject() *Project {
	return &Project{}
}

// go-nunu/nunu new and create does it well, used here as a reference
func run(cmd *cobra.Command, args []string) {
	logger := log.GetLogger(cmd.Context())
	logger.Info("scaffolding repo")
	var spacecore string
	cmd.Flags().StringVarP(&spacecore, "spacecore-name", "s", "", "Name of the spacecore")

	p := NewProject()
	if repoURL != "" {
		p.ProjectName = spacecore
	}
	prompter := utils.NewPrompter()

	if len(args) == 0 {
		spacecore, err := prompter.InputString(
			"Name your Spacecore:",
			"fancy-spacecore",
			"Please enter a valid name for your Spacecore",
			func(s string) error {
				return nil
			},
		)
		if err != nil {
			return
		}
		p.ProjectName = spacecore

		stat, _ := os.Stat(p.ProjectName)
		if stat != nil {
			overwrite, err := prompter.Confirm("Do you want to overwrite the existing spacecore project?")
			if err != nil {
				logger.Error(err)
				return
			}
			if !overwrite {
				return
			}
			err = os.RemoveAll(p.ProjectName)
			if err != nil {
				fmt.Println("remove old project error: ", err)
				return
			}
		}

	} else {
		p.ProjectName = args[0]
	}

	yes, err := p.cloneTemplate()
	if err != nil {
		logger.Error(err)
		return
	}

	if !yes {
		fmt.Printf("Folder %s already exists, do you want to overwrite it?\n", p.ProjectName)
		return
	}
	err = p.replacePackageName()
	if err != nil || !yes {
		return
	}

	// TODO: pass file names dynamically as more templates get added
	applyTemplate("grpc", p.ProjectName)
	// applyTemplate("hac.toml", p.ProjectName)

	err = p.modTidy()
	if err != nil || !yes {
		return
	}

	p.rmGit()

	fmt.Printf("\n")
	fmt.Printf("\x1B[38;2;255;105;180mV\x1B[39m")
	fmt.Printf("\x1B[38;2;255;182;193mi\x1B[39m")
	fmt.Printf("\x1B[38;2;135;206;235mm\x1B[39m")
	fmt.Printf("\x1B[38;2;144;238;144ma\x1B[39m")
	fmt.Printf("\x1B[38;2;255;255;102mn\x1B[39m")
	fmt.Printf("\x1B[38;2;255;165;0ma\x1B[39m")
	fmt.Printf("\n")
	fmt.Printf("🎉 Project \u001B[36m%s\u001B[0m created successfully!\n\n", p.ProjectName)
	fmt.Printf("Done. Now run:\n\n")
	fmt.Printf("› \033[36mcd %s \033[0m\n", p.ProjectName)
	// go build -v -o bin/spacecods .
	fmt.Printf("› \033[36mgo build -o bin/%s . \033[0m\n", p.ProjectName)
	fmt.Printf("› \033[36mvimana plugin %s/bin/%s $PLUGIN_PATH start \033[0m\n\n", p.ProjectName, p.ProjectName)
}

func applyTemplate(filename string, spacecore string) {
	logger := log.GetLogger(context.Background())
	spacecoreBytes, err := scaffold.SpacecoreBytes(filename, spacecore)
	if err != nil {
		logger.Error(err)
		return
	}
	formattedScBytes, err := format.Source(spacecoreBytes)
	if err != nil {
		logger.Error(err)
		return
	}

	// get filepath extension and remove .gotmpl or .tomltmpl, only remove tmpl
	// filename = filename[:len(filename)-5]
	filePath := fmt.Sprintf("%s/%s", spacecore, filename)
	fmt.Printf("Creating %s\n", filePath)
	err = scaffold.WriteBytes(filePath, formattedScBytes)
	if err != nil {
		logger.Error(err)
		return
	}

}

func (p *Project) cloneTemplate() (bool, error) {
	// stat, _ := os.Stat(p.ProjectName)
	repo := config.RepoBase
	cmd := exec.Command("git", "clone", repo, p.ProjectName)
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("git clone %s error: %s\n", repo, err)
		return false, err
	}
	return true, nil
}

func GetProjectName(dir string) string {
	modFile, err := os.Open(dir + "/go.mod")
	if err != nil {
		fmt.Println("go.mod does not exist", err)
		return ""
	}
	defer modFile.Close()

	var moduleName string
	_, err = fmt.Fscanf(modFile, "module %s", &moduleName)
	if err != nil {
		fmt.Println("read go mod error: ", err)
		return ""
	}
	return moduleName
}
func (p *Project) replacePackageName() error {
	packageName := GetProjectName(p.ProjectName)

	err := p.replaceFiles(packageName)
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "mod", "edit", "-module", p.ProjectName)
	cmd.Dir = p.ProjectName
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("go mod edit error: ", err)
		return err
	}
	return nil
}
func (p *Project) modTidy() error {
	fmt.Println("go mod tidy")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = p.ProjectName
	if err := cmd.Run(); err != nil {
		fmt.Println("go mod tidy error: ", err)
		return err
	}
	return nil
}

func (p *Project) rmGit() {
	os.RemoveAll(p.ProjectName + "/.git")
}

func (p *Project) replaceFiles(packageName string) error {
	err := filepath.Walk(p.ProjectName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		newData := bytes.ReplaceAll(data, []byte(packageName), []byte(p.ProjectName))
		if err := os.WriteFile(path, newData, 0644); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Println("walk file error: ", err)
		return err
	}
	return nil
}
