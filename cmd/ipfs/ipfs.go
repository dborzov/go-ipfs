package main

import (
	"fmt"
	"os"

	"github.com/gonuts/flag"
	"github.com/jbenet/commander"
	config "github.com/jbenet/go-ipfs/config"
	core "github.com/jbenet/go-ipfs/core"
	u "github.com/jbenet/go-ipfs/util"
)

// The IPFS command tree. It is an instance of `commander.Command`.
var CmdIpfs = &commander.Command{
	UsageLine: "ipfs [<flags>] <command> [<args>]",
	Short:     "global versioned p2p merkledag file system",
	Long: `ipfs - global versioned p2p merkledag file system

Basic commands:

    add <path>    Add an object to ipfs.
    cat <ref>     Show ipfs object data.
    ls <ref>      List links from an object.
    refs <ref>    List link hashes from an object.

Tool commands:

    config        Manage configuration.
    version       Show ipfs version information.
    commands      List all available commands.

Advanced Commands:

    mount         Mount an ipfs read-only mountpoint.

Use "ipfs help <command>" for more information about a command.
`,
	Run: ipfsCmd,
	Subcommands: []*commander.Command{
		cmdIpfsAdd,
		cmdIpfsCat,
		cmdIpfsLs,
		cmdIpfsRefs,
		cmdIpfsConfig,
		cmdIpfsVersion,
		cmdIpfsCommands,
		cmdIpfsMount,
	},
	Flag: *flag.NewFlagSet("ipfs", flag.ExitOnError),
}

func ipfsCmd(c *commander.Command, args []string) error {
	u.POut(c.Long)
	return nil
}

var cfg *config.Config

func main() {
	var err error
	cfg, err = config.Load("")
	if err != nil {
		fmt.Printf("Unable to load configuration file: %s \n", err)
		return
	}

	if len(os.Args) > 1 {
		if os.Args[1] != "config" && os.Args[1] != "update" {
			// when user attempts to update or tweaks the config variables
			// she should not be prevented to do that with the Update check
			// (so that if someone gets the version obsolete error message, the suggested fixes are actually working)
			cfg.Updates.Check = "ignore"
		}
	}

	if cfg.Updates.Check != "ignore" {
		// we don't check for updates whenever explicitly forbidden in config
		err := checkForUpdates()
		if err != nil {
			fmt.Println(err)
			if cfg.Updates.Check != "warn" {
				return
			}
		}
	}

	err = CmdIpfs.Dispatch(os.Args[1:])
	if err != nil {
		if len(err.Error()) > 0 {
			fmt.Fprintf(os.Stderr, "ipfs %s: %v\n", os.Args[1], err)
		}
		os.Exit(1)
	}
	return
}

func localNode() (*core.IpfsNode, error) {
	return core.NewIpfsNode(cfg)
}
