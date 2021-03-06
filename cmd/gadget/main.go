package main

import (
	"flag"
	//~ "fmt"
	"os"
	"path/filepath"
	"strings"
	"errors"
	log "github.com/sirupsen/logrus"
)

var (
	Version   = "unknown"
	GitCommit = "unknown"
)

type GadgetCommandFunc func([]string, *GadgetContext) error

type GadgetCommand struct {
	Name        string
	Function    GadgetCommandFunc
	NeedsConfig bool
}

var Commands = []GadgetCommand {
	{ Name: "init",    Function: GadgetInit,    NeedsConfig: false },
	{ Name: "add",     Function: GadgetAdd,     NeedsConfig: true  },
	{ Name: "build",   Function: GadgetBuild,   NeedsConfig: true  },
	{ Name: "deploy",  Function: GadgetDeploy,  NeedsConfig: true  },
	{ Name: "start",   Function: GadgetStart,   NeedsConfig: true  },
	{ Name: "stop",    Function: GadgetStop,    NeedsConfig: true  },
	{ Name: "status",  Function: GadgetStatus,  NeedsConfig: true  },
	{ Name: "delete",  Function: GadgetDelete,  NeedsConfig: true  },
	{ Name: "shell",   Function: GadgetShell,   NeedsConfig: false },
	{ Name: "logs",    Function: GadgetLogs,    NeedsConfig: true  },
	{ Name: "run",     Function: GadgetRun,     NeedsConfig: false },
	{ Name: "version", Function: GadgetVersion, NeedsConfig: false },
	{ Name: "help",    Function: GadgetHelp,    NeedsConfig: false },
}

func GadgetVersion(args []string, g *GadgetContext) error {
	log.Infoln(filepath.Base(os.Args[0]))
	log.Infof("  version: %s\n", Version)
	log.Infof("  commit: %s\n", GitCommit)
	return nil
}

func GadgetHelp(args []string, g *GadgetContext) error {
	flag.Usage()
	return nil
}

func FindCommand(name string) (*GadgetCommand, error) {
	for _,cmd := range Commands {
		if cmd.Name == name {
			return &cmd,nil
		}
	}
	return nil, errors.New("Failed to find command")
}

func main() {
	// Hey, Listen! 
	// Everything that outputs needs to come after g.Verbose check!
	flag.Usage = func() {
		log.Info ("")
		log.Infof("USAGE: %s [options] COMMAND", filepath.Base(os.Args[0]))
		log.Info ("")
		log.Info ("Commands:")
		log.Info ("  init        Initialize gadget project")
		log.Info ("  add         Initialize gadget project")
		log.Info ("  build       Build gadget config file")
		log.Info ("  deploy      Build gadget config file")
		log.Info ("  start       Build gadget config file")
		log.Info ("  stop        Build gadget config file")
		log.Info ("  status      Build gadget config file")
		log.Info ("  delete      Build gadget config file")
		log.Info ("  shell       Connect to remote device running GadgetOS")
		log.Info ("  logs        Build gadget config file")
		log.Info ("  version     Print version information")
		log.Info ("  help        Print this message")
		log.Info ("")
		log.Infof("Run '%s COMMAND --help' for more information on the command", filepath.Base(os.Args[0]))
		log.Info ("")
		log.Infof("Options:")
		log.Info ("  -C string                             ")
		log.Info ("    	Run in directory (default \".\")  ")
		log.Info ("  -v	Verbose execution                 ")
		log.Info ("")
	}

	g := GadgetContext{}
	
	flag.BoolVar(&g.Verbose, "v", false, "Verbose execution")
	flag.StringVar(&g.WorkingDirectory, "C", ".", "Run in directory")
	flag.Parse()

	if g.Verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	
	// Hey, Listen! 
	// Everything that outputs needs to come after g.Verbose check!
	

	err := RequiredSsh()
	if err != nil {
		log.Error("Failed to verify ssh requirements")
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		log.Error("No Command Specified")
		os.Exit(1)
	}
		
	// file command
	cmd,err := FindCommand(args[0])
	if err != nil {
		flag.Usage()
		log.WithFields(log.Fields{
			"command": strings.Join(args[0:], " "),
		}).Error("Command is not valid")
		os.Exit(1)
	}

	// if command needs to use the config file, load it
	if cmd.NeedsConfig {
		err = g.LoadConfig()
		if err != nil {
			log.Error("Failed to load config")
			log.Warn("Be sure to run gadget in the same directory as 'gadget.yml'")
			log.Warn("Or specify a directory e.g. 'gadget -C ../projects/gpio/ [command]'")
			os.Exit(1)
		}
	}

	err = cmd.Function(args[1:], &g)
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
