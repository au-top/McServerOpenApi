package main

import (
	"flag"
	"fmt"
	"fromMCServerGetMod/src/function"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

type modDownloadFileInfo struct {
	filename string
	length   int64
}
type modFileInfo struct {
	filePath string
	fileName string
}

func main() {

	var gamePath string = ""
	var configPath string = ""
	flag.StringVar(&configPath, "config", "", "配置文件路径")
	flag.StringVar(&gamePath, "gamePath", "", "游戏路径")
	gamePath, _ = filepath.Abs(gamePath)
	configPath, _ = filepath.Abs(configPath)
	flag.Parse()
	println("gamePath: " + gamePath + "\nconfig: " + configPath)
	var configName = ""
	var serviceHost string = ""
	{
		b, _ := ioutil.ReadFile(configPath)
		var readJsonStr = string(b)
		serviceHost = gjson.Get(readJsonStr, "serviceHost").String()
		configName = gjson.Get(readJsonStr, "configName").String()
	}

	var cachePath = "./cache"
	var downloadSaveCachePath = cachePath + "/" + configName + "/"
	var gameModsDirPath = gamePath + "/mods"

	ise, err := function.Exists(downloadSaveCachePath)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if !ise {
		os.MkdirAll(downloadSaveCachePath, os.ModePerm)
	}

	isdir, err := function.IsDir(downloadSaveCachePath)
	if !isdir {
		fmt.Println(downloadSaveCachePath + " is no dir ")
		return
	}
	cahceFileList, successed := downloadModFiles(serviceHost, downloadSaveCachePath)
	if successed {
		println("download mods successed")
	} else {
		println("download mods failed")
		return
	}
	_ = os.Rename(gameModsDirPath, gamePath+"/mods_backup_"+strconv.FormatInt(time.Now().UnixNano(), 10))
	_ = os.MkdirAll(gameModsDirPath, os.ModePerm)
	copyCahceToMods(cahceFileList, gameModsDirPath)

}

func copyCahceToMods(cahceFileList []modFileInfo, gameModsDirPath string) {
	for _, fileInfo := range cahceFileList {
		rFile, err := ioutil.ReadFile(fileInfo.filePath)
		if err != nil {
			fmt.Println(err)
		}
		err = ioutil.WriteFile(gameModsDirPath+"/"+fileInfo.fileName, rFile, 0644)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func downloadModFiles(serviceHost string, downloadSaveCachePath string) ([]modFileInfo, bool) {
	var modList = getModsList(serviceHost)
	var cacheModList = make([]modFileInfo, 0)
	var successed = false
	for i := 0; i < len(modList); i++ {
		modIndex := modList[i]
		var downloadLength int64 = 0
		reDownloadCount := 0
		downloadFileCahcePath := downloadSaveCachePath + modIndex.filename
		//read cache
		downloadFileIsEx, _ := function.Exists(downloadFileCahcePath)
		downloadModFix := serviceHost + "/modSource/"
		if downloadFileIsEx {
			finfo, _ := os.Stat(downloadFileCahcePath)
			downloadLength = finfo.Size()
		}
		if downloadFileIsEx && downloadLength == modIndex.length {
			fmt.Println("file " + modIndex.filename + " read cache ")
			cacheModList = append(cacheModList, modFileInfo{
				filePath: downloadFileCahcePath,
				fileName: modIndex.filename,
			})
		} else {
			for downloadLength != modIndex.length && reDownloadCount < 4 {
				downloadLength = downloadFile(downloadModFix+modIndex.filename, downloadFileCahcePath)
				if downloadLength == modIndex.length {
					fmt.Println(modIndex.filename + " download " + strconv.FormatInt(modIndex.length, 10) + "b" + " successed")
					cacheModList = append(cacheModList, modFileInfo{
						filePath: downloadFileCahcePath,
						fileName: modIndex.filename,
					})
					break
				} else {
					reDownloadCount++
					fmt.Println(modIndex.filename + " download failed ")
					fmt.Println(modIndex.filename + " is " + strconv.FormatInt(modIndex.length, 10) + "b no " + strconv.FormatInt(downloadLength, 10) + " b")
				}
			}
		}
	}

	println("download " + strconv.FormatInt(int64(len(cacheModList)), 10) + " modfiles count: " + strconv.FormatInt(int64(len(modList)), 10) + " successed")
	successed = len(cacheModList) == len(modList)
	return cacheModList, successed
}

func readDirCache(cachePath string) []os.FileInfo {
	readDir, err := ioutil.ReadDir(cachePath)
	if err != nil {
		fmt.Println(err)
		return make([]os.FileInfo, 0)
	}
	return readDir
}

func downloadFile(url string, saveFile string) int64 {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)

		return 0
	}
	defer res.Body.Close()
	out, err := os.Create(saveFile)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	defer out.Close()
	_, err = io.Copy(out, res.Body)
	if err != nil {
		fmt.Println(err)
		return 0
	} else {
		fileInfo, _ := os.Stat(saveFile)
		return fileInfo.Size()
	}
}

func getModsList(serviceHost string) []modDownloadFileInfo {
	res, err := http.Get(serviceHost + "/modInfo")
	if err != nil {
		fmt.Println("get error")
		return make([]modDownloadFileInfo, 0)
	}
	defer res.Body.Close()
	fmt.Println(res.StatusCode)

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("ReadAll error")
		return make([]modDownloadFileInfo, 0)
	}
	dataStr := string(data)
	if !gjson.Valid(dataStr) {
		fmt.Println("get mod list json error ")
		fmt.Println(dataStr)
		return make([]modDownloadFileInfo, 0)
	}
	var rawData = gjson.Parse(dataStr).Value().([]interface{})
	resArr := make([]modDownloadFileInfo, 0)
	for i := 0; i < len(rawData); i++ {
		mf := modDownloadFileInfo{filename: "", length: 0}
		readDataV := rawData[i].(map[string]interface{})
		mf.filename = readDataV["filename"].(string)
		mf.length = int64(readDataV["length"].(float64))
		resArr = append(resArr, mf)
	}
	return resArr
}
