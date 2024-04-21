package conf

import (
	"errors"
	"log"
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

var conf struct {
	links    []Link
	services []Service
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

func load(doc *document.Document) error {
	for _, node := range doc.Nodes {
		var err error

		name := node.Name.Value.(string)
		switch name {
		case "links":
			err = loadLinks(node.Children)
		case "services":
			err = loadServices(node.Children)
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
