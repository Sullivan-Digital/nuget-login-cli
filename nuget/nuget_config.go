package nuget

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/beevik/etree"
)

const defaultNugetConfigFileName = "NuGet.config"

func isPointingToFile(hint string) bool {
	if hint == "" {
		return false
	}

	s, err := os.Stat(hint)
	if err == nil {
		return !s.IsDir()
	}

	if strings.HasSuffix(hint, "/") {
		return false
	}

	return filepath.Ext(hint) != ""
}

func GetNugetConfigPath(hint string) string {
	if hint == "" {
		return getDefaultNugetConfigPath()
	}

	if isPointingToFile(hint) {
		return hint
	}

	return filepath.Join(hint, defaultNugetConfigFileName)
}

func getDefaultNugetConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	defaultPath := filepath.Join(home, ".nuget", defaultNugetConfigFileName)
	paths := []string{
		defaultPath,
		filepath.Join(home, ".nuget", "nuget.config"),
		filepath.Join(home, ".nuget", "NuGet", "nuget.config"),
		filepath.Join(home, ".nuget", "NuGet", "NuGet.config"),
	}

	found := firstFilepathThatExists(paths...)
	if found != "" {
		return found
	}

	return defaultPath
}

func firstFilepathThatExists(paths ...string) string {
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

func GetNameForNugetSource(filePath string, urlOrName string) (string, error) {
	doc, err := readNugetConfig(filePath)
	if os.IsNotExist(err) {
		return defaultNameForNugetSource(urlOrName), nil
	}

	if err != nil {
		return "", err
	}

	sources := doc.FindElement("configuration/packageSources")

	for _, child := range sources.ChildElements() {
		if child.Tag == "add" && (child.SelectAttr("value").Value == urlOrName || child.SelectAttr("key").Value == urlOrName) {
			return child.SelectAttr("key").Value, nil
		}
	}

	return defaultNameForNugetSource(urlOrName), nil
}

func defaultNameForNugetSource(urlOrName string) string {
	u, err := url.Parse(urlOrName)
	if err != nil {
		return urlOrName
	}

	path := strings.TrimSuffix(u.Path, "/index.json")
	path = strings.ReplaceAll(path, "/", "-")
	return u.Host + path
}

func AddSourceToNugetConfig(filePath string, name string, url string) error {
	doc, err := readNugetConfig(filePath)
	if err != nil {
		return err
	}

	configuration := doc.FindElement("configuration")
	if configuration == nil {
		configuration = etree.NewElement("configuration")
		doc.AddChild(configuration)
	}

	sources := configuration.FindElement("packageSources")
	if sources == nil {
		sources = etree.NewElement("packageSources")
		configuration.AddChild(sources)
	}

	source := etree.NewElement("add")
	source.CreateAttr("key", name)
	source.CreateAttr("value", url)

	sources.AddChild(source)

	return writeNugetConfig(filePath, doc)
}

func AddPackageSourceCredentialsToNugetConfig(filePath string, name string, username string, password string) error {
	doc, err := readNugetConfig(filePath)
	if err != nil {
		return err
	}

	configuration := doc.FindElement("configuration")
	if configuration == nil {
		configuration = etree.NewElement("configuration")
		doc.AddChild(configuration)
	}

	credentials := configuration.FindElement("packageSourceCredentials")
	if credentials == nil {
		credentials = etree.NewElement("packageSourceCredentials")
		configuration.AddChild(credentials)
	}

	var packageSource *etree.Element
	for _, child := range credentials.ChildElements() {
		if child.Tag == name {
			packageSource = child
			break
		}
	}

	if packageSource == nil {
		packageSource = etree.NewElement(name)
		credentials.AddChild(packageSource)
	}

	packageSource.Child = []etree.Token{}
	usernameElement := etree.NewElement("add")
	usernameElement.CreateAttr("key", "Username")
	usernameElement.CreateAttr("value", username)
	packageSource.AddChild(usernameElement)

	passwordElement := etree.NewElement("add")
	passwordElement.CreateAttr("key", "ClearTextPassword")
	passwordElement.CreateAttr("value", password)
	packageSource.AddChild(passwordElement)

	return writeNugetConfig(filePath, doc)
}

func AddMappingToNugetConfig(filePath string, name string, pattern string) error {
	doc, err := readNugetConfig(filePath)
	if err != nil {
		return err
	}

	configuration := doc.FindElement("configuration")
	if configuration == nil {
		configuration = etree.NewElement("configuration")
		doc.AddChild(configuration)
	}

	mappings := configuration.FindElement("packageSourceMapping")
	if mappings == nil {
		mappings = etree.NewElement("packageSourceMapping")
		configuration.AddChild(mappings)
	}

	var packageSource *etree.Element
	for _, child := range mappings.ChildElements() {
		if child.Tag == "packageSource" && child.SelectAttr("key").Value == name {
			packageSource = child
			break
		}
	}

	if packageSource == nil {
		packageSource = etree.NewElement("packageSource")
		packageSource.CreateAttr("key", name)
		mappings.AddChild(packageSource)
	}

	mapping := etree.NewElement("package")
	mapping.CreateAttr("pattern", pattern)

	packageSource.AddChild(mapping)

	return writeNugetConfig(filePath, doc)
}

// ReadNugetConfig reads a NuGet config file and returns the XML document.
func readNugetConfig(filePath string) (*etree.Document, error) {
	if _, err := os.Stat(filePath); err != nil {
		return nil, err
	}

	doc := etree.NewDocument()
	if err := doc.ReadFromFile(filePath); err != nil {
		return nil, err
	}

	return doc, nil
}

func writeNugetConfig(filePath string, doc *etree.Document) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	doc.Indent(2)
	return doc.WriteToFile(filePath)
}
