package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// CPU Struct
type CPU struct {
	cpu []string
}

func (c *CPU) parse(output []byte) (sum, work int) {
	outStr := string(output)

	// SPLIT CPU
	i := strings.Index(outStr, "\n")
	c.cpu = strings.SplitAfter(outStr[:i], " ")

	// SUM USAGE
	totalJiffies := 0
	for _, val := range c.cpu[2:] {
		i, _ := strconv.Atoi(strings.TrimRight(val, " "))
		totalJiffies += i
	}

	workJiffies := 0
	for _, val := range c.cpu[2:5] {
		i, _ := strconv.Atoi(strings.TrimRight(val, " "))
		workJiffies += i
	}

	return totalJiffies, workJiffies
}

func handleError(e error) {
	if e != nil {
		panic(e)
	}
}

func getCPUOutput() []byte {
	cmd := exec.Command("cat", "/proc/stat")
	out, err := cmd.Output()
	handleError(err)

	return out
}

func main() {
	// Assign Polybar Color from Args
	iconClr := "#FF" // Default
	strClr := "#FF"  // Default

	// Verify Given Arguments
	if len(os.Args) > 1 {
		iconClr = os.Args[1]
	}
	if len(os.Args) > 2 {
		strClr = os.Args[2]
	}

	// CONFIG USED
	interval := 1 * time.Second // Seconds

	// START RUNNING
	cpu := CPU{}
	var tJ1, tJ2, tW1, tW2 int

	// INITIAL VALUES
	tJ1, tW1 = cpu.parse(getCPUOutput())

	for {
		time.Sleep(interval)                 // Every Interval
		tJ2, tW2 = cpu.parse(getCPUOutput()) // Get Current Data

		// Calculate Usage
		dWork := tW2 - tW1
		dJiff := tJ2 - tJ1
		usage := (float32(dWork) / float32(dJiff)) * 100.00
		//		fmt.Printf("Usage: %.2f%%\n", usage)
		fmt.Printf("%%{F%s} %%{F%s}%.1f%%\n", iconClr, strClr, usage)

		// Store Previous Values
		tW1 = tW2
		tJ1 = tJ2
	}
}
