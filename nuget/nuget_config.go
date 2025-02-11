package nuget

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/beevik/etree"
)

const defaultNugetConfigFileName = "NuGet.config"

func isPointingToFile(hint string, verbose bool) bool {
	if verbose {
		fmt.Println("Checking if hint points to a file:", hint)
	}
	if hint == "" {
		return false
	}

	s, err := os.Stat(hint)
	if err == nil {
		if verbose {
			fmt.Println("Hint is a file:", !s.IsDir())
		}
		return !s.IsDir()
	}

	if strings.HasSuffix(hint, "/") {
		return false
	}

	return filepath.Ext(hint) != ""
}

func GetNugetConfigPath(hint string, verbose bool) string {
	if verbose {
		fmt.Println("Getting NuGet config path with hint:", hint)
	}
	if hint == "" {
		return getDefaultNugetConfigPath(verbose)
	}

	if isPointingToFile(hint, verbose) {
		return hint
	}

	return filepath.Join(hint, defaultNugetConfigFileName)
}

func getDefaultNugetConfigPath(verbose bool) string {
	if verbose {
		fmt.Println("Getting default NuGet config path")
	}
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

	found := firstFilepathThatExists(verbose, paths...)
	if found != "" {
		return found
	}

	return defaultPath
}

func firstFilepathThatExists(verbose bool, paths ...string) string {
	if verbose {
		fmt.Println("Checking for first existing filepath in paths:", paths)
	}
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			if verbose {
				fmt.Println("Found existing filepath:", path)
			}
			return path
		}
	}

	return ""
}

func GetNameForNugetSource(filePath string, urlOrName string, verbose bool) (string, error) {
	if verbose {
		fmt.Println("Getting name for NuGet source from file:", filePath, "with URL or name:", urlOrName)
	}
	doc, err := readNugetConfig(filePath, verbose)
	if os.IsNotExist(err) {
		return defaultNameForNugetSource(urlOrName, verbose), nil
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

	return defaultNameForNugetSource(urlOrName, verbose), nil
}

func defaultNameForNugetSource(urlOrName string, verbose bool) string {
	if verbose {
		fmt.Println("Generating default name for NuGet source:", urlOrName)
	}
	u, err := url.Parse(urlOrName)
	if err != nil {
		return urlOrName
	}

	path := strings.TrimSuffix(u.Path, "/index.json")
	path = strings.ReplaceAll(path, "/", "-")
	return u.Host + path
}

func AddSourceToNugetConfig(filePath string, name string, url string, verbose bool) error {
	if verbose {
		fmt.Println("Adding source to NuGet config:", filePath, "Name:", name, "URL:", url)
	}

	doc, err := readNugetConfig(filePath, verbose)
	if os.IsNotExist(err) {
		doc = etree.NewDocument()
	} else if err != nil {
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

	return writeNugetConfig(filePath, doc, verbose)
}

func AddPackageSourceCredentialsToNugetConfig(filePath string, name string, username string, password string, verbose bool) error {
	if verbose {
		fmt.Println("Adding package source credentials to NuGet config:", filePath, "Name:", name)
	}

	doc, err := readNugetConfig(filePath, verbose)
	if os.IsNotExist(err) {
		doc = etree.NewDocument()
	} else if err != nil {
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

	return writeNugetConfig(filePath, doc, verbose)
}

func AddMappingToNugetConfig(filePath string, name string, pattern string, verbose bool) error {
	if verbose {
		fmt.Println("Adding mapping to NuGet config:", filePath, "Name:", name, "Pattern:", pattern)
	}
	doc, err := readNugetConfig(filePath, verbose)
	if os.IsNotExist(err) {
		doc = etree.NewDocument()
	} else if err != nil {
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

	return writeNugetConfig(filePath, doc, verbose)
}

// ReadNugetConfig reads a NuGet config file and returns the XML document.
func readNugetConfig(filePath string, verbose bool) (*etree.Document, error) {
	if verbose {
		fmt.Println("Reading NuGet config from file:", filePath)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, err
	}

	doc := etree.NewDocument()
	if err := doc.ReadFromFile(filePath); err != nil {
		return nil, err
	}

	return doc, nil
}

func writeNugetConfig(filePath string, doc *etree.Document, verbose bool) error {
	if verbose {
		fmt.Println("Writing NuGet config to file:", filePath)
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	doc.Indent(2)
	return doc.WriteToFile(filePath)
}
