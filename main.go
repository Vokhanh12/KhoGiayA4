package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const NUMTHREAD = 10

type KhoGiay struct {
	width  int
	length int
	x      int
	v      int
}

func main() {

	fmt.Println("Hello world")

	filePath := "DanhHCN_31012024.txt"

	// Open the file in
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// handle create file out
	fileName := "outputDanhHCN.txt"

	fileOut, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer fileOut.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	lineNumber := 0

	// Iterate through each line in the file

	queue := make(chan KhoGiay, 100)

	go func() {

		for scanner.Scan() {

			if lineNumber == 0 {
				lineNumber++
				continue
			}

			line := scanner.Text()

			arrs := strings.Fields(line)

			width, err := strconv.Atoi(arrs[0])
			if err != nil {
				fmt.Println("Loi doi width string sang int")
				return
			}

			length, err := strconv.Atoi(arrs[1])
			if err != nil {
				fmt.Println("Loi doi length string sang int")
				return
			}

			queue <- KhoGiay{width: width, length: length}

		}

		close(queue)

	}()

	for i := 0; i < NUMTHREAD; i++ {

		go func(threadNum int) {

			for v := range queue {

				width := float64(v.width) / 100

				length := float64(v.length) / 100

				chx := make(chan float64, 100)

				go func() {

					for i := 0.0; i <= float64(width)/2; i += 0.001 {
						// x range 0 <= x <= w/2
						chx <- i
					}

				}()

				// goutine chia 5

				Vmax := 0.0
				Xmax := 0.0

				for i := 0; i < threadNum; i++ {

					go func(threadNum int) {

						for ix := range chx {

							crV := ix * (float64(width) - 2*ix) * (float64(length) - 2*ix)

							if Vmax < crV {
								Vmax = crV
								Xmax = ix
							}

							time.Sleep(time.Millisecond)

						}

						close(chx)

					}(i)

				}

				time.Sleep(time.Second * 3)

				outputString := fmt.Sprintf("%f %f %f %f numThread:%d \n", width, length, Xmax, Vmax, threadNum)
				_, err = fileOut.WriteString(outputString)
				if err != nil {
					fmt.Println("Error writing to file:", err)
					return
				}

				fmt.Println()
				time.Sleep(time.Second)

			}
		}(i)

	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	time.Sleep(time.Second * 30)

}
