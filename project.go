package naos

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mholt/archiver"
)

// InventoryFileName specifies the inventory file name.
const InventoryFileName = "naos.json"

// HiddenDirectory specifies the hidden directory name.
const HiddenDirectory = ".naos"

// A Project is a project available on disk.
type Project struct {
	Location  string
	Inventory *Inventory
}

// CreateProject will initialize a project in the specified directory.
func CreateProject(path string) (*Project, error) {
	// ensure p directory
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, err
	}

	// create p
	p := &Project{
		Location:  path,
		Inventory: NewInventory(),
	}

	// save inventory
	err = p.SaveInventory()
	if err != nil {
		return nil, err
	}

	return p, nil
}

// FindProject will look for project in the specified path.
func FindProject(path string) (*Project, error) {
	// attempt to read inventory
	inv, err := ReadInventory(filepath.Join(path, InventoryFileName))
	if err != nil {
		return nil, err
	}

	// TODO: Walk up the tree and try parent directories.

	// prepare project
	project := &Project{
		Location:  path,
		Inventory: inv,
	}

	return project, nil
}

// SaveInventory will save the associated inventory to disk.
func (p *Project) SaveInventory() error {
	// save inventory
	err := p.Inventory.Save(filepath.Join(p.Location, InventoryFileName))
	if err != nil {
		return err
	}

	return nil
}

// HiddenDirectory returns the hidden used to store the toolchain, development
// framework and other necessary files.
func (p *Project) HiddenDirectory() string {
	return filepath.Join(p.Location, HiddenDirectory)
}

// SetupToolchain will setup the compilation toolchain. An existing toolchain
// will be removed if force is set to true. If out is not nil, it will be used
// to log information about the process.
func (p *Project) SetupToolchain(force bool, out io.Writer) error {
	// get toolchain url
	var url string
	switch runtime.GOOS {
	case "darwin":
		url = "https://dl.espressif.com/dl/xtensa-esp32-elf-osx-1.22.0-61-gab8375a-5.2.0.tar.gz"
	case "linux":
		url = "https://dl.espressif.com/dl/xtensa-esp32-elf-linux64-1.22.0-61-gab8375a-5.2.0.tar.gz"
	default:
		return errors.New("unsupported os")
	}

	// prepare toolchain directory
	toolchainDir := filepath.Join(p.HiddenDirectory(), "xtensa-esp32-elf")

	// check if already exists
	ok, err := exists(toolchainDir)
	if err != nil {
		return err
	}

	// return immediately if already exists and not forced
	if ok && !force {
		log(out, fmt.Sprintln("Skipping toolchain as it already exists"))
		return nil
	}

	// remove existing directory if existing
	if ok {
		log(out, fmt.Sprintln("Removing existing toolchain (force=true)"))
		err = os.RemoveAll(toolchainDir)
		if err != nil {
			return err
		}
	}

	// get a temporary file
	tmp, err := ioutil.TempFile("", "naos")
	if err != nil {
		return err
	}

	// make sure temporary file gets closed
	defer tmp.Close()

	// download toolchain
	log(out, fmt.Sprintf("Downloading toolchain from '%s'...\n", url))
	err = download(tmp.Name(), url)
	if err != nil {
		return err
	}

	// unpack toolchain
	log(out, fmt.Sprintf("Unpacking toolchain to '%s'...\n", p.HiddenDirectory()))
	err = archiver.TarGz.Open(tmp.Name(), p.HiddenDirectory())
	if err != nil {
		return err
	}

	// close temporary file
	tmp.Close()

	// remove temporary file
	err = os.Remove(tmp.Name())
	if err != nil {
		return err
	}

	return nil
}

// ToolchainLocation returns the location of the toolchain if it exists or an
// error if it does not exist.
func (p *Project) ToolchainLocation() (string, error) {
	// calculate directory
	dir := filepath.Join(p.HiddenDirectory(), "xtensa-esp32-elf")

	// check if toolchain directory exists
	ok, err := exists(dir)
	if err != nil {
		return "", err
	} else if !ok {
		return "", errors.New("toolchain not found")
	}

	return dir, nil
}

// SetupDevelopmentFramework will setup the development framework. An existing
// development framework will be removed if force is set to true. If out is not
// nil, it will be used to log information about the process.
func (p *Project) SetupDevelopmentFramework(force bool, out io.Writer) error {
	// prepare toolchain directory
	frameworkDir := filepath.Join(p.HiddenDirectory(), "esp-idf")

	// check if already exists
	ok, err := exists(frameworkDir)
	if err != nil {
		return err
	}

	// return immediately if already exists and not forced
	if ok && !force {
		log(out, fmt.Sprintln("Skipping development framework as it already exists"))
		return nil
	}

	// remove existing directory if existing
	if ok {
		log(out, fmt.Sprintln("Removing existing development framework (force=true)"))
		err = os.RemoveAll(frameworkDir)
		if err != nil {
			return err
		}
	}

	// clone development framework
	err = clone("https://github.com/espressif/esp-idf.git", frameworkDir, ESPIDFVersion, out)
	if err != nil {
		return err
	}

	return nil
}

// DevelopmentFrameworkLocation returns the location of the development
// framework if it exists or an error if it does not exist.
func (p *Project) DevelopmentFrameworkLocation() (string, error) {
	// calculate directory
	dir := filepath.Join(p.HiddenDirectory(), "esp-idf")

	// check if toolchain directory exists
	ok, err := exists(dir)
	if err != nil {
		return "", err
	} else if !ok {
		return "", errors.New("development framework not found")
	}

	return dir, nil
}

// SetupBuildTree will setup the build tree. An existing build tree will be
// removed if force is set to true. If out is not nil, it will be used to log
// information about the process.
func (p *Project) SetupBuildTree(force bool, out io.Writer) error {
	// prepare build tree directory
	buildTreeDir := filepath.Join(p.HiddenDirectory(), "tree")

	// check if already exists
	ok, err := exists(buildTreeDir)
	if err != nil {
		return err
	}

	// return immediately if already exists and no forced
	if ok && !force {
		log(out, fmt.Sprintln("Skipping build tree as it already exists"))
		return nil
	}

	// remove existing directory if existing
	if ok {
		log(out, fmt.Sprintln("Removing existing build tree (force=true)"))
		err = os.RemoveAll(buildTreeDir)
		if err != nil {
			return err
		}
	}

	// adding directory
	log(out, fmt.Sprintf("Adding build tree directory to '%s'\n", buildTreeDir))
	err = os.MkdirAll(buildTreeDir, 0777)
	if err != nil {
		return err
	}

	// adding sdk config file
	log(out, fmt.Sprintf("Adding 'sdkconfig' to '%s'\n", buildTreeDir))
	err = ioutil.WriteFile(filepath.Join(buildTreeDir, "sdkconfig"), []byte(SdkconfigContent), 0644)
	if err != nil {
		return err
	}

	// adding makefile
	log(out, fmt.Sprintf("Adding 'Makefile' to '%s'\n", buildTreeDir))
	err = ioutil.WriteFile(filepath.Join(buildTreeDir, "Makefile"), []byte(MakefileContent), 0644)
	if err != nil {
		return err
	}

	// adding main directory
	log(out, fmt.Sprintf("Adding 'main' directory to '%s'\n", buildTreeDir))
	err = os.MkdirAll(filepath.Join(buildTreeDir, "main"), 0777)
	if err != nil {
		return err
	}

	// adding components directory
	log(out, fmt.Sprintf("Adding 'components' directory to '%s'\n", buildTreeDir))
	err = os.MkdirAll(filepath.Join(buildTreeDir, "components"), 0777)
	if err != nil {
		return err
	}

	// construct esp-mqtt component dir
	espMQTTDir := filepath.Join(buildTreeDir, "components", "esp-mqtt")

	// clone component
	err = clone("https://github.com/256dpi/esp-mqtt.git", espMQTTDir, ESPMQTTVersion, out)
	if err != nil {
		return err
	}

	return nil
}

// BuildTreeLocation returns the location of the build tree if it exists or an
// error if it does not exist.
func (p *Project) BuildTreeLocation() (string, error) {
	// calculate directory
	dir := filepath.Join(p.HiddenDirectory(), "tree")

	// check if build tree directory exists
	ok, err := exists(dir)
	if err != nil {
		return "", err
	} else if !ok {
		return "", nil
	}

	return dir, nil
}

// Build will build the project.
func (p *Project) Build(out io.Writer) error {
	// get toolchain location
	toolchain, err := p.ToolchainLocation()
	if err != nil {
		return err
	}

	// get build tree location
	buildTree, err := p.BuildTreeLocation()
	if err != nil {
		return err
	}

	// construct build command
	cmd := exec.Command("make")

	// set working directory
	cmd.Dir = buildTree

	// inherit current environment
	cmd.Env = os.Environ()

	// go through all env variables
	for i, str := range cmd.Env {
		if strings.HasPrefix(str, "PATH=") {
			// prepend toolchain bin directory
			cmd.Env[i] = "PATH=" + filepath.Join(toolchain, "bin") + ":" + os.Getenv("PATH")
		} else if strings.HasPrefix(str, "PWD=") {
			// override shell working directory
			cmd.Env[i] = "PWD=" + buildTree
		}
	}

	// connect output if provided
	if out != nil {
		cmd.Stdout = out
		cmd.Stderr = out
	}

	// build project
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
