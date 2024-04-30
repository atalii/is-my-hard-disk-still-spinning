package conf

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/sblinch/kdl-go"
	"github.com/sblinch/kdl-go/document"
)

type Service struct {
	Name string
}

type Link struct {
	Name string
	Url  string
}

type Meta struct {
	Slogans []string
}

var conf struct {
	links    []Link
	services []Service
	meta     *Meta
}

func ReadConf(path string) error {
	reader, err := os.Open(path)
	if err != nil {
		return err
	}

	if doc, err := kdl.Parse(reader); err != nil {
		return err
	} else {
		err = load(doc)
		return err
	}
}

func Links() []Link {
	return conf.links
}

func Services() []Service {
	return conf.services
}

func Slogan() string {
	n := len(conf.meta.Slogans)

	if n == 0 {
		return "meow?? why is anyone reading this"
	} else {
		return conf.meta.Slogans[rand.Intn(n)]
	}
}

func load(doc *document.Document) error {
	for _, node := range doc.Nodes {
		var err error

		name := node.Name.Value.(string)
		switch name {
		case "links":
			err = loadLinks(node.Children)
		case "services":
			err = loadServices(node.Children)
		case "meta":
			err = loadMeta(node.Children)
		}

		if err != nil {
			log.Printf("while loading config: %v", err)
			return err
		}
	}

	return nil
}

func loadLinks(links []*document.Node) error {
	if conf.links != nil {
		return errors.New("links declared twice")
	}

	for _, linkNode := range links {
		urlNode, _ := linkNode.Properties.Get("url")

		name := linkNode.Name.Value.(string)
		url := urlNode.Value.(string)
		conf.links = append(conf.links, Link{
			Name: name,
			Url:  url,
		})
	}

	return nil
}

func loadServices(services []*document.Node) error {
	if conf.services != nil {
		return errors.New("services declared twice")
	}

	for _, svNode := range services {
		name := svNode.Name.Value.(string)
		conf.services = append(conf.services, Service{
			Name: name,
		})
	}

	return nil
}

func loadMeta(nodes []*document.Node) error {
	if conf.meta != nil {
		return errors.New("meta declared twice")
	}

	conf.meta = &Meta{}

	for _, node := range nodes {
		name := node.Name.Value.(string)
		switch name {
		case "slogan":
			if err := loadSlogan(node.Arguments); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unexpected node: %s", name)
		}
	}

	return nil
}

func loadSlogan(slogans []*document.Value) error {
	if len(slogans) != 1 {
		return errors.New("slogan nodes should only have one argument")
	}

	slogan := slogans[0].Value.(string)
	conf.meta.Slogans = append(conf.meta.Slogans, slogan)
	return nil
}
