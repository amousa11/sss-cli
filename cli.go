package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/amousa11/sss"
	"github.com/amousa11/sss/utils"

	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "Generation of Shares and Recovery of Secrets - Shamir Secret Sharing"
	app.Usage = "Generate Shares and Recover Secrets using a threshold of those shares"
	app.Commands = []cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g", "gen"},
			Usage:   "generate [minNumShares] [numShares] - generates numShares shares of a secret which can only be recovered by numMinShares shares",
			Action: func(c *cli.Context) error {
				generateSharesCLI(c)
				return nil
			},
		},
		{
			Name:    "recover",
			Aliases: []string{"c"},
			Usage:   "recover [pathToSharesFolder] - recovers a secret given a path to a folder with the generated shares",
			Action: func(c *cli.Context) error {
				recoverSecretsCLI(c)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func generateSharesCLI(c *cli.Context) {
	m, err := strconv.Atoi(c.Args().Get(0))
	if err != nil {
		log.Fatal(err.Error())
	}

	n, err := strconv.Atoi(c.Args().Get(1))
	if err != nil {
		log.Fatal(err.Error())
	}

	prime, _ := big.NewInt(1).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)

	fmt.Println("FIELD MODULUS =", prime.Text(16))
	fmt.Println("Generating ", n, "shares with threshold", m, " for recovery:")
	secret, points, e := sss.GenerateShares(m, n, prime)
	if e != nil {
		log.Fatal(e.Error())
	}
	fmt.Println("Secret : ", secret.Text(16))
	fmt.Println("Shares : ")
	for i := 0; i < len(points); i++ {
		fmt.Println(points[i].X.Text(16), "\t", points[i].Y.Text(16))
		dir := "./shares/"
		createDirIfNotExist(dir)
		dir += points[i].X.Text(10)
		f, err := os.Create(dir)
		check(err, c)

		fmt.Fprintf(f, "%x %x %x", points[i].X, points[i].Y, prime)

		f.Close()
	}
}

func recoverSecretsCLI(c *cli.Context) {
	var modulus *big.Int
	shareFiles := getShareFiles(c.Args().Get(0))
	points := make([]*utils.Point, len(shareFiles))
	for i := 0; i < len(shareFiles); i++ {
		dat, err := ioutil.ReadFile(shareFiles[i])
		if err != nil {
			log.Fatal(err)
		}
		// Share file format: X Y modulus
		// Reads the byte array into a string and spilits the string by spaces into an array
		share := strings.Split(string(dat[:]), " ")

		xString := share[0]
		yString := share[1]
		mod := share[2]

		xValue, okX := new(big.Int).SetString(xString, 16)
		yValue, okY := new(big.Int).SetString(yString, 16)

		if !okX {
			log.Fatal(fmt.Errorf("Error parsing xValue from share: %s", string(dat)))
		}

		if !okY {
			log.Fatal(fmt.Errorf("Error parsing yValue from share: %s", string(dat)))
		}

		if i == 0 {
			firstModulus, okMod := new(big.Int).SetString(mod, 16)
			if !okMod {
				log.Fatal(fmt.Errorf("Error parsing modulus from share: %s", string(dat)))
			}
			modulus = firstModulus
		}

		point := new(utils.Point)
		point.X = xValue.Mod(xValue, modulus)
		point.Y = yValue.Mod(yValue, modulus)

		points[i] = point
	}
	recoveredSecret, e := sss.RecoverSecret(points, modulus)

	if e != nil {
		log.Fatal(e.Error())
	}

	fmt.Printf("Secret successfully recovered from %d shares: %x\n", len(points), recoveredSecret)
}

func getShareFiles(path string) (out []string) {
	pwd, _ := os.Getwd()
	files, err := ioutil.ReadDir(pwd + "/shares/")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		out = append(out, pwd+"/shares/"+f.Name())
	}
	return out
}

func check(e error, c *cli.Context) {
	if e != nil {
		log.Fatal(e)
	}
}

// credit: https://siongui.github.io/2017/03/28/go-create-directory-if-not-exist/
func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}
