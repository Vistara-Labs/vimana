package scaffold

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

// Scaffold is a struct that holds the template and the data to be used to generate the code
type Scaffold struct {
	Template         string
	Data             interface{}
	PluginName       string
	MagicCookieKey   string
	MagicCookieValue string
}

// NewScaffold returns a new Scaffold object
func NewScaffold(template string, data interface{}) *Scaffold {
	return &Scaffold{
		Template: template,
		Data:     data,
	}
}

// Execute generates the code using the template and the data
func (s *Scaffold) Execute() (string, error) {
	tmpl, err := template.New("scaffold").Parse(s.Template)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, s.Data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// SpacecoreBytes generates the spacecore bytes
func SpacecoreBytes(filename string, spacecore string) ([]byte, error) {
	var tmpl *template.Template
	var err error

	tmpl, err = template.ParseFS(createTemplateFS, fmt.Sprintf("templates/%s.gotmpl", filename))
	if err != nil {
		return nil, err
	}
	spacecoreData := struct {
		PackageName      string
		Data             string
		PluginName       string
		MagicCookieKey   string
		MagicCookieValue string
	}{
		PackageName: spacecore,
		Data:        spacecore,
		PluginName:  spacecore,
		// if this needs to be coming from flags, we also need to change plugins/plugin.go
		MagicCookieKey:   "SPACECORE_PLUGIN",
		MagicCookieValue: "v1",
	}
	var buf bytes.Buffer

	// Execute the template and overwrite if the file already exists
	// tmpl.ExecuteTemplate()

	if err := tmpl.Execute(&buf, spacecoreData); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func WriteBytes(filepath string, data []byte) error {

	// if _, err := os.Stat(filepath); !errors.Is(err, fs.ErrNotExist) {
	// 	return fmt.Errorf("file (%s) already exists", filepath)
	// }

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}
