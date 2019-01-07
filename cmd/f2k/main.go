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
	name     = flag.String("name", "", "name to use for the label")
	port     = flag.Int("port", 80, "expose this port")
	replicas = flag.Int("replicas", 1, "replicas")

	concurrencyPolicy = flag.String("concurrencyPolicy", "Forbid", "Allow, Replace, or Forbid")
	restartPolicy     = flag.String("restartPolicy", "OnFailure", "Always, OnFailure, or Never")

	vlan = flag.String("vlan", "16", "import address from this interface")
)

func init() {
	flag.Parse()
}

func doDeployService(name string, u *unit.Unit, output io.Writer) error {
	ip := u.Network["bond0."+*vlan]
	if ip == "" {
		ip = u.Network["br"+*vlan]
	}

	ip = strings.Split(ip, "/")[0]
	if ip != "" {
		fmt.Fprint(output, "---\n")
		if err := yaml.NewEncoder(output).Encode(kubes.NewService(name, ip, *port)); err != nil {
			return err
		}
	}

	fmt.Fprint(output, "---\n")
	err := yaml.NewEncoder(output).Encode(
		kubes.NewDeployment(name, u.RunImage, u.RunCommand, *replicas, *port, u.Env),
	)
	return err
}

func doCronJob(filename, name string, u *unit.Unit, output io.Writer) error {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	timer, err := unit.NewUnit(string(raw))
	if err != nil {
		return err
	}

	schedule := timer.Contents["Timer"]["OnCalendar"][0]
	annotations := kubes.Annotations{}
	if unit, ok := timer.Contents["Unit"]; ok {
		if description, ok := unit["Description"]; ok {
			annotations["description"] = description[0]
		}
		if documentation, ok := unit["Documentation"]; ok {
			annotations["documentation"] = documentation[0]
		}
	}

	fmt.Fprintf(output, "---\n")
	err = yaml.NewEncoder(output).Encode(
		kubes.NewCronJob(name, schedule, *concurrencyPolicy, *restartPolicy,
			u.RunImage, u.RunCommand, u.Env, annotations),
	)
	return err
}

func do(filename string, output io.Writer) error {
	basename := filename
	ext := filepath.Ext(filename)
	if ext == "" {
		filename += ".service"
	} else {
		basename = basename[:len(basename)-len(ext)]
	}

	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	u, err := unit.NewUnit(string(raw))
	if err != nil {
		return err
	}

	if *name == "" {
		*name = filepath.Base(strings.Split(filename, ".service")[0])
		*name = strings.Replace(*name, ".", "-", -1)
	}

	timerFname := filename[:len(filename)-len(filepath.Ext(filename))] + ".timer"
	if _, err := os.Stat(timerFname); os.IsNotExist(err) {
		return doDeployService(*name, u, output)
	}

	return doCronJob(timerFname, *name, u, output)
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
