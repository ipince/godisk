package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	Mil = 1000
	B   = 1
	KB  = Mil * B
	MB  = Mil * KB
	GB  = Mil * MB
)

type TreeNode struct {
	part       string // for debugging only
	selfTotal  int64  // of elements at this level, not including children
	totalSet   bool
	grandTotal int64
	children   map[string]*TreeNode
}

func NewNode(part string) *TreeNode {
	return &TreeNode{
		part:     part,
		children: make(map[string]*TreeNode),
	}
}

func (t *TreeNode) Total() int64 {
	if t.totalSet {
		return t.grandTotal
	}
	total := t.selfTotal
	for _, c := range t.children {
		total += c.Total()
	}
	t.grandTotal = total
	t.totalSet = true
	return total
}

func (t *TreeNode) String() string {
	return fmt.Sprintf("%d %s", t.selfTotal, t.part)
}

func (t *TreeNode) Print() string {
	return t.print("", true, 0)
}

func (t *TreeNode) print(prefix string, last bool, depth int) string {
	pad := "|--"
	childPad := "|  "
	if last {
		pad = "`--"
		childPad = "   "
	}
	str := fmt.Sprintf("%s%s %s %s (self %s)\n", prefix, pad, HumanReadable(t.Total()), t.part, HumanReadable(t.selfTotal))
	if depth > 1 {
		return str
	}

	childrenSlice := make([]*TreeNode, 0, len(t.children))
	for _, c := range t.children {
		childrenSlice = append(childrenSlice, c)
	}
	sort.Slice(childrenSlice, func(i, j int) bool {
		return childrenSlice[i].Total() > childrenSlice[j].Total()
	})
	for i, c := range childrenSlice {
		str += c.print(prefix+childPad, i == len(childrenSlice)-1, depth+1)
	}

	return str
}

func DirSize(top string) (*TreeNode, error) {
	// map from dir to selfTotal...
	// but then i have to walk it in order, how?
	root := NewNode(top)
	err := filepath.Walk(top, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Find the right node.
		parts := strings.Split(strings.TrimPrefix(path, top), string(os.PathSeparator))
		node := root
		//fmt.Println(parts)
		if len(parts) > 0 { // root
			for _, p := range parts[:len(parts)-1] {
				if n, ok := node.children[p]; ok { // walk down
					node = n
				} else { // not found, so create a new child
					node.children[p] = NewNode(p)
					node = node.children[p]
				}
			}
		}

		// find TreeNode for this dir, if not exist, then add?
		if !info.IsDir() {
			//dir := filepath.Dir(path)
			//fmt.Printf("adding %s to %s tally for file %s\n", HumanReadable(info.Size()), dir, info.Name())
			node.selfTotal += info.Size()
		}

		return nil
	})

	return root, err
}

func HumanReadable(size int64) string {
	if size > GB {
		return fmt.Sprintf("%6.2f GB", float64(size)/float64(GB))
	} else if size > MB {
		return fmt.Sprintf("%6.2f MB", float64(size)/float64(MB))
	} else if size > KB {
		return fmt.Sprintf("%6.2f KB", float64(size)/float64(KB))
	}
	return fmt.Sprintf("%6.2f B", float64(size))
}

func main() {
	dir := os.Args[1]

	sizes, err := DirSize(dir)
	if err != nil {
		log.Printf("failed to calculate size of directory %s: %v\n", dir, err)
	}
	fmt.Println(sizes.Print())
}
