package syncfile

import (
	"sync"
	"testing"
	"fmt"
)

func readFile(filePath string) {

}

func Test_syncfile(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0;i < 10; i++ {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			fileInfo, err := UsingFiles.ROpen(filePath)
			if err == nil {
				defer UsingFiles.Close(filePath)
				data := make([]byte, 100)
				lenth, err := fileInfo.GetFilePtr().ReadAt(data, 0)
				if err != nil {
					fmt.Printf("ReadAt is err: %v\n", err)
					return
				}
				data = data[:lenth]
				fmt.Printf("data is %s\n", string(data))
			}
			fmt.Printf("ROpen is err: %v\n", err)
		}("./syncfile.go")
	}
	wg.Wait()
}
