package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/manifoldco/promptui"
)

const prev = "-parent-"

var (
	Build     = "devel"
	V         = flag.Bool("version", false, "Show Version.")
	H         = flag.Bool("help", false, "Show help.")
	log       = GetLogger()
	templates = &promptui.SelectTemplates{
		Help:     "Use: \x1b[91m↓(j) ↑(k) → ← \x1b[0m Search: \x1b[91m/\x1b[0m ",
		Label:    "✨ {{ . | green}}",
		Active:   "\U0001F437  {{ .Name | cyan }} {{if .Host}}{{.Host | cyan}}{{end}}",
		Inactive: "  {{.Name | faint}} ",
	}
)

func main() {
	flag.Parse()
	if !flag.Parsed() {
		flag.Usage()
		return
	}

	if *H {
		flag.Usage()
		return
	}
	if *V {
		fmt.Println("go version: ", runtime.Version())
		fmt.Println("sk: v1.0.0")
		return
	}

	for {
		err := LoadConfig()
		if err != nil {
			log.Error("load config error", err)
			os.Exit(1)
		}
		node := choose(nil, GetConfig())
		if node == nil {
			return
		}
		client := NewClient(node)
		client.Login()

	}
}

func choose(parent, trees []*Node) *Node {
	prompt := promptui.Select{
		Label:     "Select Host:",
		Items:     trees,
		Templates: templates,
		Searcher: func(input string, index int) bool {
			node := trees[index]
			content := fmt.Sprintf("%s %s %s", node.Name, node.User, node.Host)
			if strings.Contains(input, " ") {
				for _, key := range strings.Split(input, " ") {
					key = strings.TrimSpace(key)
					if key != "" {
						if !strings.Contains(content, key) {
							return false
						}
					}
				}
				return true
			}
			if strings.Contains(content, input) {
				return true
			}
			return false
		},
	}
	index, _, err := prompt.Run()
	if err != nil {
		return nil
	}
	node := trees[index]
	if len(node.Children) > 0 {
		first := node.Children[0]
		if first.Name != prev {
			first = &Node{Name: prev}
			node.Children = append(node.Children[:0], append([]*Node{first}, node.Children...)...)
		}
		return choose(trees, node.Children)
	}
	if node.Name == prev {
		return choose(nil, parent)
	}
	return node
}
