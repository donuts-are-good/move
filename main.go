package main

import (
	"flag"
	"fmt"
	"hash/crc32"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Progress struct {
	BytesMoved int64
	TotalBytes int64
	FilesMoved int
	TotalFiles int
}

func main() {
	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Println("Error: Too few arguments.")
		fmt.Println("Usage: move /path/to/source /path/to/destination")
		os.Exit(1)
	}

	if flag.NArg() > 2 {
		fmt.Println("Error: Too many arguments.")
		fmt.Println("Usage: move /path/to/source /path/to/destination")
		os.Exit(1)
	}

	source := flag.Arg(0)
	destination := flag.Arg(1)

	fmt.Printf("moving: %s -> %s\n", source, destination)

	progressChan := make(chan Progress)
	go displayProgress(progressChan)

	startTime := time.Now()

	var totalFiles int
	var totalBytes int64
	filepath.Walk(source, func(srcPath string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			totalFiles++
			totalBytes += info.Size()
		}
		return nil
	})

	var progress Progress
	progress.TotalFiles = totalFiles
	progress.TotalBytes = totalBytes

	var wg sync.WaitGroup

	err := filepath.Walk(source, func(srcPath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(source, srcPath)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(destination, relPath)

		if info.IsDir() {
			return os.Mkdir(dstPath, info.Mode())
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := moveFile(srcPath, dstPath, info.Mode(), progressChan, &progress)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				os.Exit(1)
			}
		}()

		return nil
	})
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	wg.Wait()
	close(progressChan)

	duration := time.Since(startTime)
	fmt.Printf("\nMoving took %s... done!\n", duration)
}

func moveFile(src, dst string, mode fs.FileMode, progressChan chan<- Progress, progress *Progress) error {
	err := copyFile(src, dst, mode, progressChan, progress)
	if err != nil {
		return err
	}

	if !compareFiles(src, dst) {
		return fmt.Errorf("CRC check failed for: %s and %s", src, dst)
	}

	err = os.Remove(src)
	if err != nil {
		return err
	}

	return nil
}

func compareFiles(file1, file2 string) bool {
	f1, err := os.Open(file1)
	if err != nil {
		return false
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		return false
	}
	defer f2.Close()

	h1 := crc32.NewIEEE()
	h2 := crc32.NewIEEE()

	_, err = io.Copy(h1, f1)
	if err != nil {
		return false
	}

	_, err = io.Copy(h2, f2)
	if err != nil {
		return false
	}

	return h1.Sum32() == h2.Sum32()
}

func copyFile(src, dst string, mode fs.FileMode, progressChan chan<- Progress, progress *Progress) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, 32*1024)
	for {
		n, err := srcFile.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := dstFile.Write(buf[:n]); err != nil {
			return err
		}

		progress.BytesMoved += int64(n)
		progressChan <- *progress

		if err == io.EOF {
			break
		}
	}

	progress.FilesMoved++
	progressChan <- *progress

	return nil
}

func displayProgress(progressChan <-chan Progress) {
	for progress := range progressChan {
		fmt.Printf("\r%d files - %.1fGB [%s]",
			progress.FilesMoved,
			float64(progress.BytesMoved)/(1<<30),
			progressBar(progress.BytesMoved, progress.TotalBytes, 20),
		)
	}
}

func progressBar(current, total int64, width int) string {
	ratio := float64(current) / float64(total)
	completeChars := int(ratio * float64(width))
	progressBar := strings.Repeat("=", completeChars)
	if completeChars < width {
		progressBar += ">"
		remainingChars := width - completeChars - 1
		progressBar += strings.Repeat(" ", remainingChars)
	}
	return progressBar
}
