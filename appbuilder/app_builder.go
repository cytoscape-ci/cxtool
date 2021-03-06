/*
	Basic settings for cxtool app.

	(An utility to create new App settings)
 */
package appbuilder

import (
	"os"
	"github.com/codegangsta/cli"
	"github.com/cyService/cxtool/converter"
	"github.com/cyService/cxtool/converter/fromcyjs"
	"bufio"
	"io"
	"log"
	"github.com/cyService/cxtool/converter/filter"
)


func BuildApp() *cli.App {
	app := cli.NewApp()

	basicSetup(app)
	setFlags(app)
	setAction(app)

	return app
}


func basicSetup(app *cli.App) {
	app.Name = Name
	app.Usage = Usage
	app.Version = Version
}


func setFlags(app *cli.App) {

	app.Flags = []cli.Flag {
		cli.StringFlag {
			Name:  "in-format, i",
			Value: "cx",
			Usage: "Source file format.  Default input file format is CX.",
		},

		cli.StringFlag {
			Name:  "out-format, o",
			Value: "cyjs",
			Usage: "Output file format. Default is Cytoscape.js.",
		},
	}
}

func setAction(app *cli.App) {
	app.Action = func(c *cli.Context) {
		commandLineArgs := c.Args()

		inFileFormat := c.String("in-format")
		outFileFormat := c.String("out-format")

		con := getConverter(inFileFormat, outFileFormat)

		// Two cases: Run from file or piped text stream
		if len(commandLineArgs) == 0 {

			fi, err := os.Stdin.Stat()
			if err != nil {
				panic(err)
			}

			if fi.Mode() & os.ModeNamedPipe == 0 {
				// Show help menu if there is no input
				cli.ShowAppHelp(c)
			} else {
				// No param.  Use Pipe
				r := bufio.NewReader(os.Stdin)
				w := io.Writer(os.Stdout)
				con.Convert(r, w)
			}
		} else {
			source := commandLineArgs[0]

			f, err := os.Open(source)
			if err != nil {
				log.Fatal("Error:", err)
				return
			}

			// Close input file at the end of this
			defer f.Close()

			w := io.Writer(os.Stdout)
			con.Convert(bufio.NewReader(f), w)
		}
	}
}

func getConverter(inFormat string, outFormat string) converter.Converter {
	switch inFormat {
	case cx:
		if outFormat == cytoscapejs {
			return converter.Cx2Cyjs{}
		}

		if outFormat == sif {
			return converter.Cx2Sif{}
		}

		if outFormat == nodeids {
			return filter.Cx2IdList{}
		}
	case sif:
		return converter.Sif2Cx{Delimiter:' '}
	case cytoscapejs:
		return fromcyjs.Cyjs2Cx{}
	default:
		return converter.Cx2Cyjs{}
	}
	return nil
}