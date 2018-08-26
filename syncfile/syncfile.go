package syncfile

import (
	"fmt"
	"path"
	"sync"
	"os"
)

var UsingFiles usingFiles

func init() {
	UsingFiles = usingFiles{
		fileInfos: make(map[string]*usingFile),
	}
}

type usingFile struct {
	filePath string
	filePtr  *os.File
	using    int
}

func (u usingFile)GetFilePtr() *os.File{
	return u.filePtr
}
type usingFiles struct {
	fileInfos map[string]*usingFile
	lock      sync.Mutex
}

func (f *usingFiles) ROpen(filepath string) (fileInfo *usingFile, err error) {
	if info, err := os.Stat(filepath); err != nil || info.IsDir() {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%s is not exist", filepath)
		}
		return nil, fmt.Errorf("%s is dir or err: %v", filepath, err)
	} else {
		//exist
		filename := path.Base(filepath)
		f.lock.Lock()
		defer f.lock.Unlock()
		fileInfo, ok := f.fileInfos[filename]
		if ok {
			filePtr := fileInfo.filePtr
			if _, err := filePtr.Stat(); err == nil{
				fileInfo.using += 1
				fmt.Printf("filepath using is %d\n", fileInfo.using)
				return fileInfo, nil
			}
		}
		filePtr, err := os.OpenFile(filepath, os.O_RDONLY, 0640)
		if err != nil {
			return nil, fmt.Errorf("open file[%s] is error: %v", filepath, err)
		}
		fileInfo = &usingFile{
			filePath: filepath,
			filePtr: filePtr,
			using:    1,
		}
		fmt.Printf("filepath[%s] register: %+v\n",filepath, fileInfo)
		f.fileInfos[filename] = fileInfo
		return fileInfo, nil
	}
}

func (f *usingFiles) Close(filepath string){
	filename := path.Base(filepath)
	f.lock.Lock()
	defer f.lock.Unlock()
	fileInfo, ok := f.fileInfos[filename]
	if ok {
		filePtr := fileInfo.filePtr
		if _, err := filePtr.Stat(); err == nil{
			if fileInfo.using <= 1 {
				filePtr.Close()
				fmt.Printf("filepath[%s] unregister: %+v\n",filepath, fileInfo)
				delete(f.fileInfos, filename)
			} else {
				fileInfo.using -= 1
				return
			}
		}
	}
}
