package naos

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/kr/pty"
	"github.com/shiftr-io/naos/espidf"
	"github.com/shiftr-io/naos/mqtt"
	"github.com/shiftr-io/naos/utils"
	"github.com/shiftr-io/naos/xtensa"
)

// TODO: Implement dependency updating.

// TODO: Read versions from naos-esp?

const idfVersion = "9b955f4c9f1b32652ea165d3e4cdaad01bba170e"
const mqttVersion = "cc87172126aa7aacc3b982f7be7489950429b733"
const comVersion = "41c8c6a839eb20c46f19258b8fdcc5caa3aba01c"

// A Project is a project available on disk.
type Project struct {
	Location  string
	Inventory *Inventory
}

// CreateProject will initialize a project in the specified directory. If out is
// not nil, it will be used to log information about the process.
func CreateProject(path string, out io.Writer) (*Project, error) {
	// ensure project directory
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, err
	}

	// create project
	p := &Project{
		Location:  path,
		Inventory: NewInventory(),
	}

	// save inventory
	err = p.SaveInventory()
	if err != nil {
		return nil, err
	}

	// print info
	utils.Log(out, "Created new empty inventory.")
	utils.Log(out, "Please update the settings to suit your needs.")

	// ensure source directory
	utils.Log(out, "Ensuring source directory.")
	err = os.MkdirAll(filepath.Join(path, "src"), 0755)
	if err != nil {
		return nil, err
	}

	// prepare main source path and check if it already exists
	mainSourcePath := filepath.Join(path, "src", "main.c")
	ok, err := utils.Exists(mainSourcePath)
	if err != nil {
		return nil, err
	}

	// create main source file if it not already exists
	if !ok {
		utils.Log(out, "Creating default 'main.c' source file.")
		err = ioutil.WriteFile(mainSourcePath, []byte(mainSourceFile), 0644)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

// OpenProject will open the project in the specified path.
func OpenProject(path string) (*Project, error) {
	// attempt to read inventory
	inv, err := ReadInventory(filepath.Join(path, "naos.json"))
	if err != nil {
		return nil, err
	}

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
	err := p.Inventory.Save(filepath.Join(p.Location, "naos.json"))
	if err != nil {
		return err
	}

	return nil
}

// InternalDirectory returns the internal directory used to store the toolchain,
// development framework and other necessary files.
func (p *Project) InternalDirectory() string {
	return filepath.Join(p.Location, "naos")
}

// TreeDirectory returns the internal directory that contains the whole build tree.
func (p *Project) TreeDirectory() string {
	return filepath.Join(p.InternalDirectory(), "tree")
}

// Setup will setup necessary dependencies. Any existing dependencies will be
// removed if force is set to true. If out is not nil, it will be used to log
// information about the process.
func (p *Project) Setup(force bool, cmake bool, out io.Writer) error {
	// install xtensa toolchain
	err := xtensa.Install(p.InternalDirectory(), force, out)
	if err != nil {
		return err
	}

	// install esp32 development framework
	err = espidf.Install(p.InternalDirectory(), idfVersion, force, out)
	if err != nil {
		return err
	}

	// check if already exists
	ok, err := utils.Exists(p.TreeDirectory())
	if err != nil {
		return err
	}

	// return immediately if already exists and no forced
	if ok && !force {
		utils.Log(out, "Skipping build tree as it already exists.")
		return nil
	}

	// remove existing directory if existing
	if ok {
		utils.Log(out, "Removing existing build tree (forced).")
		err = os.RemoveAll(p.TreeDirectory())
		if err != nil {
			return err
		}
	}

	// print log
	utils.Log(out, "Creating build tree directory structure...")

	// adding directory
	err = os.MkdirAll(p.TreeDirectory(), 0777)
	if err != nil {
		return err
	}

	// adding sdk config file
	err = ioutil.WriteFile(filepath.Join(p.TreeDirectory(), "sdkconfig"), []byte(sdkconfigFile), 0644)
	if err != nil {
		return err
	}

	// adding makefile
	err = ioutil.WriteFile(filepath.Join(p.TreeDirectory(), "Makefile"), []byte(makeFile), 0644)
	if err != nil {
		return err
	}

	// adding main directory
	err = os.MkdirAll(filepath.Join(p.TreeDirectory(), "main"), 0777)
	if err != nil {
		return err
	}

	// adding component.mk
	err = ioutil.WriteFile(filepath.Join(p.TreeDirectory(), "main", "component.mk"), []byte(componentFile), 0644)
	if err != nil {
		return err
	}

	// linking src
	err = os.Symlink(filepath.Join(p.Location, "src"), filepath.Join(p.TreeDirectory(), "main", "src"))
	if err != nil {
		return err
	}

	// adding components directory
	err = os.MkdirAll(filepath.Join(p.TreeDirectory(), "components"), 0777)
	if err != nil {
		return err
	}

	// clone component
	utils.Log(out, "Installing MQTT component...")
	err = utils.Clone("https://github.com/256dpi/esp-mqtt.git", filepath.Join(p.TreeDirectory(), "components", "esp-mqtt"), mqttVersion, out)
	if err != nil {
		return err
	}

	// clone component
	utils.Log(out, "Installing NAOS component...")
	err = utils.Clone("https://github.com/shiftr-io/naos-esp.git", filepath.Join(p.TreeDirectory(), "components", "naos-esp"), comVersion, out)
	if err != nil {
		return err
	}

	if cmake {
		// write internal cmake file
		utils.Log(out, "Creating internal CMake file.")
		err := ioutil.WriteFile(filepath.Join(p.InternalDirectory(), "CMakeLists.txt"), []byte(internalCMakeListsFile), 0644)
		if err != nil {
			return err
		}

		// get project path
		projectPath := filepath.Join(p.Location, "CMakeLists.txt")
		ok, err := utils.Exists(projectPath)
		if err != nil {
			return err
		}

		if !ok || force {
			utils.Log(out, "Creating project CMake file.")
			err = ioutil.WriteFile(projectPath, []byte(projectCMakeListsFile), 0644)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Build will build the project.
func (p *Project) Build(clean, appOnly bool, out io.Writer) error {
	// clean project if requested
	if clean {
		utils.Log(out, "Cleaning project...")
		err := p.exec(out, nil, "make", "clean")
		if err != nil {
			return err
		}
	}

	// build project (app only)
	if appOnly {
		utils.Log(out, "Building project (app only)...")
		err := p.exec(out, nil, "make", "app")
		if err != nil {
			return err
		}

		return nil
	}

	// build project
	utils.Log(out, "Building project...")
	err := p.exec(out, nil, "make", "all")
	if err != nil {
		return err
	}

	return nil
}

// Flash will flash the project to the attached device.
func (p *Project) Flash(device string, erase bool, appOnly bool, out io.Writer) error {
	// calculate paths
	espTool := filepath.Join(espidf.Directory(p.InternalDirectory()), "components", "esptool_py", "esptool", "esptool.py")
	bootLoaderBinary := filepath.Join(p.TreeDirectory(), "build", "bootloader", "bootloader.bin")
	projectBinary := filepath.Join(p.TreeDirectory(), "build", "naos-project.bin")
	partitionsBinary := filepath.Join(p.TreeDirectory(), "build", "partitions_two_ota.bin")

	// prepare erase flash command
	eraseFlash := []string{
		espTool,
		"--chip", "esp32",
		"--port", device,
		"--baud", "921600",
		"--before", "default_reset",
		"--after", "hard_reset",
		"erase_flash",
	}

	// prepare flash all command
	flashAll := []string{
		espTool,
		"--chip", "esp32",
		"--port", device,
		"--baud", "921600",
		"--before", "default_reset",
		"--after", "hard_reset",
		"write_flash",
		"-z",
		"--flash_mode", "dio",
		"--flash_freq", "40m",
		"--flash_size", "detect",
		"0x1000", bootLoaderBinary,
		"0x10000", projectBinary,
		"0x8000", partitionsBinary,
	}

	// prepare flash app command
	flashApp := []string{
		espTool,
		"--chip", "esp32",
		"--port", device,
		"--baud", "921600",
		"--before", "default_reset",
		"--after", "hard_reset",
		"write_flash",
		"-z",
		"--flash_mode", "dio",
		"--flash_freq", "40m",
		"--flash_size", "detect",
		"0x10000", projectBinary,
	}

	// erase attached device if requested
	if erase {
		utils.Log(out, "Erasing flash...")
		err := p.exec(out, nil, "python", eraseFlash...)
		if err != nil {
			return err
		}
	}

	// flash attached device (app only)
	if appOnly {
		utils.Log(out, "Flashing device (app only)...")
		err := p.exec(out, nil, "python", flashApp...)
		if err != nil {
			return err
		}

		return nil
	}

	// flash attached device
	utils.Log(out, "Flashing device...")
	err := p.exec(out, nil, "python", flashAll...)
	if err != nil {
		return err
	}

	return nil
}

// Attach will attach to the attached device.
func (p *Project) Attach(device string, simple bool, out io.Writer, in io.Reader) error {
	// prepare command
	var cmd *exec.Cmd

	// set simple or advanced command
	if simple {
		// construct command
		cmd = exec.Command("miniterm.py", "--rts", "0", "--dtr", "0", "--raw", device, "115200")
	} else {
		// get path of monitor tool
		tool := filepath.Join(espidf.Directory(p.InternalDirectory()), "tools", "idf_monitor.py")

		// get elf path
		elf := filepath.Join(p.TreeDirectory(), "build", "naos-project.elf")

		// construct command
		cmd = exec.Command("python", tool, "--baud", "115200", "--port", device, elf)
	}

	// set working directory
	cmd.Dir = p.TreeDirectory()

	// inherit current environment
	cmd.Env = os.Environ()

	// go through all env variables
	for i, str := range cmd.Env {
		if strings.HasPrefix(str, "PATH=") {
			// prepend toolchain bin directory
			cmd.Env[i] = "PATH=" + xtensa.BinDirectory(p.InternalDirectory()) + ":" + os.Getenv("PATH")
		} else if strings.HasPrefix(str, "PWD=") {
			// override shell working directory
			cmd.Env[i] = "PWD=" + p.TreeDirectory()
		}
	}

	// start process and get tty
	utils.Log(out, "Attaching to device (press Ctrl+C to exit)...")
	tty, err := pty.Start(cmd)
	if err != nil {
		return err
	}

	// make sure tty gets closed
	defer tty.Close()

	// prepare channel
	quit := make(chan struct{})

	// read data until EOF
	go func() {
		io.Copy(out, tty)
		close(quit)
	}()

	// write data until EOF
	go func() {
		io.Copy(tty, in)
		close(quit)
	}()

	// wait for quit
	<-quit

	return nil
}

// Format will format all source files in the project if 'clang-format' is
// available.
func (p *Project) Format(out io.Writer) error {
	// get source and header files
	sourceFiles, headerFiles, err := p.sourceAndHeaderFiles()
	if err != nil {
		return err
	}

	// prepare arguments
	arguments := []string{"-style", "{BasedOnStyle: Google, ColumnLimit: 120}", "-i"}
	arguments = append(arguments, sourceFiles...)
	arguments = append(arguments, headerFiles...)

	// format source files
	err = p.exec(out, nil, "clang-format", arguments...)
	if err != nil {
		return err
	}

	return nil
}

// Update will update the devices that match the supplied glob pattern with the
// previously built image. The specified callback is called for every change in
// state or progress.
func (p *Project) Update(pattern string, timeout time.Duration, callback func(*Device, *mqtt.UpdateStatus)) error {
	// get image path
	image := filepath.Join(p.InternalDirectory(), "tree", "build", "naos-project.bin")

	// read image
	bytes, err := ioutil.ReadFile(image)
	if err != nil {
		return err
	}

	// run update
	err = p.Inventory.Update(pattern, bytes, timeout, callback)
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) exec(out io.Writer, in io.Reader, name string, arg ...string) error {
	// construct command
	cmd := exec.Command(name, arg...)

	// set working directory
	cmd.Dir = p.TreeDirectory()

	// connect output and inputs
	cmd.Stdout = out
	cmd.Stderr = out
	cmd.Stdin = in

	// inherit current environment
	cmd.Env = os.Environ()

	// go through all env variables
	for i, str := range cmd.Env {
		if strings.HasPrefix(str, "PATH=") {
			// prepend toolchain bin directory
			cmd.Env[i] = "PATH=" + xtensa.BinDirectory(p.InternalDirectory()) + ":" + os.Getenv("PATH")
		} else if strings.HasPrefix(str, "PWD=") {
			// override shell working directory
			cmd.Env[i] = "PWD=" + p.TreeDirectory()
		}
	}

	// run command
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) sourceAndHeaderFiles() ([]string, []string, error) {
	// prepare list
	sourceFiles := make([]string, 0)
	headerFiles := make([]string, 0)

	// scan directory
	err := filepath.Walk(filepath.Join(p.Location, "src"), func(path string, f os.FileInfo, err error) error {
		// directly return errors
		if err != nil {
			return err
		}

		// add files with matching extension
		if filepath.Ext(path) == ".c" {
			sourceFiles = append(sourceFiles, path)
		} else if filepath.Ext(path) == ".h" {
			headerFiles = append(headerFiles, path)
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return sourceFiles, headerFiles, nil
}
