package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/cogolabs/fleet2kubes/pkg/kubes"
	"github.com/cogolabs/fleet2kubes/pkg/unit"
)

var (
	port     = flag.Int("port", 80, "expose this port")
	replicas = flag.Int("n", 2, "replicas")
	vlan     = flag.String("vlan", "16", "import address from this interface")
)

func init() {
	flag.Parse()
}

func do(filename string, output io.Writer) error {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	u, err := unit.NewUnit(string(raw))
	if err != nil {
		return err
	}

	ip := u.Network["bond0."+*vlan]
	if ip == "" {
		ip = u.Network["br"+*vlan]
	}
	if ip == "" {
		return fmt.Errorf("unknown IP: %+v", u.Network)
	}
	ip = strings.Split(ip, "/")[0]

	name := filepath.Base(strings.Split(filename, ".service")[0])
	svc := kubes.NewService(name, ip, *port)
	dpl := kubes.NewDeployment(name, u.RunImage, u.RunCommand, *replicas, *port)

	err = yaml.NewEncoder(output).Encode(svc)
	if err != nil {
		return err
	}
	fmt.Fprint(output, "---\n")
	err = yaml.NewEncoder(output).Encode(dpl)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if len(flag.Args()) != 1 {
		fmt.Printf("Usage: %s [flags] <legacy.service>", os.Args[0])
		os.Exit(1)
		return
	}

	err := do(flag.Arg(0), os.Stdout)
	if err != nil {
		log.Println(err)
	}
}
