package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/cogolabs/fleet2kubes/pkg/kubes"
	"github.com/cogolabs/fleet2kubes/pkg/unit"
)

var (
	port = flag.Int("port", 80, "")
	vlan = flag.String("vlan", "16", "")
)

func init() {
	flag.Parse()
}

func do(filename string, output io.Writer) error {
	name := strings.Split(filename, ".service")[0]

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

	svc := kubes.NewService(name, ip, *port)
	dpl := kubes.NewDeployment(name, u.RunImage, u.RunCommand, *port)

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
		fmt.Printf("usage: %s <legacy.service>", os.Args[0])
		os.Exit(1)
		return
	}

	err := do(flag.Arg(0), os.Stdout)
	if err != nil {
		log.Println(err)
	}
}
