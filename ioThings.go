package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"golang.org/x/text/encoding/unicode"
)

// get flag info and argument
// NOTE: arg MUST occur AFTER flags when calling program
// icemaker -export=tmp/outfilename.tex -sigDigits=4 infilename.prb
func commandFlags(version string) (inFile fileInfo, outFile fileInfo, randomStr, sigDigits, logOut string) {
	var inFileStr string

	outFilePtr := flag.String("export", "", "outFile - REQUIRED FLAG\nFile extension should be .tex")
	sigDigitsPtr := flag.String("sigDigits", "4", "number of significant digits for output\n")
	// determines sig digits number for prob2tex
	randomPtr := flag.String("random", "false", "Choices are false, true, min, max, minMax, or positive integer")
	// determines whether parameters are default or random chosen from a set
	versionPtr := flag.Bool("version", false, "Print out version")

	flag.Parse()
	if *versionPtr {
		fmt.Println("problem2tex: ", version)
	}
	exitCode := 0
	inFileStr = flag.Arg(0)
	if inFileStr == "" {
		exitCode = 1
		fmt.Println("No input file name given\nRun with -help to see inputs required")
		os.Exit(exitCode)
	}
	if *outFilePtr == "" {
		exitCode = 1
		fmt.Println("No outFile given\nRun with -help to see inputs required")
		os.Exit(exitCode)
	}

	inFile = getFileInfo(inFileStr)
	outFile = getFileInfo(*outFilePtr)
	randomStr = *randomPtr
	_, logOut = checkRandom(randomStr)
	if logOut != "" {
		return
	}
	sigDigits, logOut = checkSigDigits(*sigDigitsPtr, logOut)
	// sigDigits = strIncrement(sigDigits, -1) // needed so that TOTAL significant digits is sigDigits
	if outFile.ext == "" {
		outFile.ext = ".log"
		outFile.full = filepath.Join(outFile.path, outFile.name+outFile.ext)
		logOut = logOut + "Output file needs a file extension of either .tex or .svg\n"
	}
	switch outFile.ext {
	case ".tex":
		if inFile.ext != ".prb" {
			logOut = logOut + "Input should be .prb file\n"
		}
	default:
		logOut = logOut + "Output file needs a file extension of .tex\n"
		outFile.ext = ".log"
		outFile.full = filepath.Join(outFile.path, outFile.name+outFile.ext)
	}
	return
}

func checkRandom(randomStr string) (int, string) {
	var random int
	var logOut string
	var err error
	switch randomStr {
	case "false", "0":
		random = 0
	case "true", "-1":
		random = -1
	case "min":
		random = -2
	case "max":
		random = -3
	case "minMax":
		random = -4
	default: //check that string is a positive integer
		random, err = strconv.Atoi(randomStr)
		if err != nil {
			random = 0
			logOut = logOut + "random should be either \"false\", \"true\", \"min\", \"max\", \"minMax\", or a positive integer"
		} else {
			if random < 1 {
				random = 0
				logOut = logOut + "random should be a positive integer"
			}
		}
	}
	return random, logOut
}

func checkSigDigits(sigDigits, logOut string) (string, string) {
	i, err := strconv.Atoi(sigDigits)
	if err != nil {
		logOut = logOut + "sigDigits should be a positive integer\n"
		sigDigits = "4"
	} else {
		if i < 1 {
			logOut = logOut + "sigDigits should be a positive integer\n"
			sigDigits = "4"
		}
	}
	return sigDigits, logOut
}

func getFileInfo(inString string) (file fileInfo) {
	var base string
	var re0 = regexp.MustCompile(`(?m)^(?P<res1>\w*)`)
	//var result []string

	file.path = filepath.Dir(inString)
	file.ext = filepath.Ext(inString)
	file.full = inString

	base = filepath.Base(inString)
	if re0.MatchString(base) {
		file.name = re0.FindStringSubmatch(base)[1]
	}
	return
}

func fileWriteString(inString, fileNameandPath string) {
	// write inString to file "fileNameandPath" (does NOT append, it overwrites)
	outbytes := []byte(inString)
	err := ioutil.WriteFile(fileNameandPath, outbytes, 0644)
	if err != nil { // if error, then create an ERROR.log file and write to it the error
		outbytes := []byte("Cannot write " + fileNameandPath + "\n")
		_ = ioutil.WriteFile("ERROR.log", outbytes, 0644) // ERROR log file created
		os.Exit(1)
	}
}

func fileAppendString(inString, fileNameandPath string) {
	// append inString to file "fileNameandPath" (will create it if it does not exist)
	f, err := os.OpenFile(fileNameandPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte(inString + "\n")); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func fileReadString(fileNameandPath string) (string, string) {
	var fileString, logOut string
	inbytes, err := ioutil.ReadFile(fileNameandPath) //
	if err != nil {
		//	fmt.Print(err)
		logOut = fmt.Sprint(err)
	}
	fileString = string(inbytes)
	return fileString, logOut
}

// func checkIfNewer(newFile, oldFile string) bool {
// 	var newer bool
// 	var commandTest string
// 	// var svgFile, pdfTexFile string
// 	// svgFile = inFile.path + inFile.name + inFile.extension
// 	// pdfTexFile = outPath + inFile.name + ".pdf_tex"
// 	commandTest = "test " + newFile + " -nt " + oldFile + " ; echo $?"
// 	out, errout, err := shellout(commandTest)
// 	if err != nil {
// 		log.Printf("error: %v\n", err)
// 	}
// 	_ = errout
// 	// fmt.Println("--- stdout ---")
// 	// fmt.Println(out)
// 	// fmt.Println("--- stderr ---")
// 	// fmt.Println(errout)
// 	var re = regexp.MustCompile(`(?m)0`)
// 	if re.MatchString(out) {
// 		newer = true
// 	} else {
// 		newer = false
// 	}
// 	return newer
// }

// func shellout(command string) (string, string, error) {
// 	var stdout bytes.Buffer
// 	var stderr bytes.Buffer
// 	cmd := exec.Command("bash", "-c", command)
// 	cmd.Stdout = &stdout
// 	cmd.Stderr = &stderr
// 	err := cmd.Run()
// 	return stdout.String(), stderr.String(), err
// }

// Checks if file is utf16 encoded and if so, it converts it to utf8 for better regex matching
func convertIfUtf16(inString string) (string, bool) {
	// requires import "golang.org/x/text/encoding/unicode"
	var inBytes []byte
	var codeUtf16 bool
	inBytes = []byte(inString)
	if len(inBytes) > 7 {
		if inBytes[1] == 0 && inBytes[3] == 0 && inBytes[5] == 0 && inBytes[7] == 0 { // VERY likely utf16 encoded so need to change to utf8
			codeUtf16 = true
			decoder := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()
			inString, _ = decoder.String(inString)
		}
	}
	return inString, codeUtf16
}
