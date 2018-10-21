/* cfg: cli tool for operate ini file
cfg get -f {filename} {section name} {key name}
cfg set -f {filename} {section name} {key name} {value}


default: filename = default.cfg
*/

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kmollee/ini"
)

func usage() {
	fmt.Fprintf(os.Stdout, `%s -f {filename} get {section_name} {key}
%s -f {filename} set  {section name} {key} {value}
`, os.Args[0], os.Args[0])
	os.Exit(0)
}

func main() {
	filepath := flag.String("f", "default.cfg", "file path to store config")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 3 {
		usage()
	}

	action := flag.Arg(0)
	switch action {
	case "get":
		section, key := flag.Arg(1), flag.Arg(2)
		f, err := os.OpenFile(*filepath, os.O_RDWR, 0660)
		must(err)
		conf, err := ini.Parse(f)
		must(err)
		val, err := conf.SectionGetKey(section, key)
		must(err)
		fmt.Fprintln(os.Stdout, val)
	case "set":
		// set should provide value
		if flag.NArg() < 4 {
			usage()
		}
		section, key, val := flag.Arg(1), flag.Arg(2), flag.Arg(3)
		f, err := os.OpenFile(*filepath, os.O_CREATE|os.O_RDWR, 0660)
		must(err)
		conf, err := ini.Parse(f)
		must(err)
		conf.SectionSetKey(section, key, val)
		f.Seek(0, 0)
		f.Truncate(0)
		must(conf.Write(f))
	default:
		usage()
	}

}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stdout, err.Error())
		os.Exit(1)
	}
}
