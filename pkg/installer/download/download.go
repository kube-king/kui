package download

import (
	"fmt"
	"io"
	"kube-invention/pkg/installer/config"
	"kube-invention/pkg/installer/constant"
	"kube-invention/pkg/installer/global"
	"kube-invention/pkg/utils/common"
	"net/http"
	"os"
)

type DownloadData struct {
}

func (d *DownloadData) Exec(config *config.Config) (err error) {

	common.MkDirs(os.ModePerm, constant.DataPath)
	dataFilePath := fmt.Sprintf("%v/data-%v-%v.tar.gz", constant.RootPath, config.KubernetsOption.Version, config.Core.Arch)

	if ok, _ := common.PathExists(dataFilePath); !ok {
		url := fmt.Sprintf("%v/data-%v-%v.tar.gz", constant.DownloadAddress, config.KubernetsOption.Version, config.Core.Arch)
		global.Log.Info(fmt.Sprintf("Run Download Data File:%v ......", url))
		response, err := http.Get(url)
		if err != nil {
			return err
		}

		defer func() {
			if response != nil && response.Body != nil {
				response.Body.Close()
			}
		}()

		file, err := os.OpenFile(dataFilePath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(file, response.Body)
		if err != nil {
			return err
		}
		global.Log.Info("Download Success")
	}

	_, err = common.DeCompress(dataFilePath, constant.RootPath)
	if err != nil {
		return err
	}

	return err
}
