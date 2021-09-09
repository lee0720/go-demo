package main

import (
	"flag"
	"fmt"
	"mime"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"gitlab.com/lilh/go-demo/internal/config"
	"gitlab.com/lilh/go-demo/internal/utils"
	"gorm.io/gorm"
)

var configFileName = flag.String("cfn", "config", "name of configs file")
var configFilePath = flag.String("cfp", "./configs", "path of configs file")

var Totalcount = 10000 * 100

var BatchSize = 10000

type Recruitment struct {
	RecID        string
	LocationName string
}

func main() {

	flag.Parse()

	err := config.InitConfig(*configFileName, strings.Split(*configFilePath, ","))

	if err != nil {
		panic(err)
	}
	utils.InitRecruitmentDatapackDB()
	f := excelize.NewFile()

	index := f.NewSheet("Sheet2")
	f.SetCellValue("Sheet2", "D1", "未匹配上的")
	f.SetCellValue("Sheet2", "A1", "底层数据")
	f.SetCellValue("Sheet2", "B1", "结构化后的数据")

	locationDeal := func(re *regexp.Regexp, idx int, offset int, records []Recruitment) bool {
		f.SetCellValue("Sheet2", fmt.Sprintf("A%d", idx+offset+2), records[idx].LocationName)

		params := re.FindStringSubmatch(records[idx].LocationName)
		if len(params) != 0 {
			f.SetCellValue("Sheet2", fmt.Sprintf("B%d", idx+offset+2), params[0])
		} else {
			return false
		}
		return true
	}
	nextID := ""
	sema := utils.NewSemaphore(10)
	maps := make(map[string]struct{})
	for i := 0; i <= Totalcount; i += BatchSize {
		sema.Add(1)
		var records []Recruitment
		err := getRecruitmentLocation(nextID, BatchSize).Scan(&records).Error
		if err != nil {
			panic(err)
		}

		if len(records) != 0 {
			nextID = records[len(records)-1].RecID
		}
		go func(records []Recruitment, idx int) {

			defer sema.Done()

			for j := range records {

				f.SetCellValue("Sheet2", fmt.Sprintf("A%d", idx+j+2), records[j])

				flag := locationDeal(utils.PrefectureLevelCity, j, idx, records)
				if flag {
					continue
				}

				flag = locationDeal(utils.CountyLevelCity, j, idx, records)
				if !flag {
					maps[records[j].LocationName] = struct{}{}
				}

			}
		}(records, i)

	}

	sema.Wait()

	f.SetActiveSheet(index)
	// 根据指定路径保存文件
	if err := f.SaveAs("Book2.xlsx"); err != nil {
		fmt.Println(err)
	}

	for key, value := range maps {
		_ = value
		fmt.Println(key)
	}

}

func getRecruitmentLocation(nextID string, BatchSize int) *gorm.DB {

	if len(nextID) == 0 {

		return utils.RecruitmentDatapackDB.Raw(`SELECT rec_id, location_name
		FROM datapack_recruitment.recruitments order by rec_id limit ?
		`, BatchSize)

	} else {
		return utils.RecruitmentDatapackDB.Raw(`SELECT rec_id,location_name
		FROM datapack_recruitment.recruitments where rec_id > ? order by rec_id limit ?
		`, nextID, BatchSize)
	}

}

func getCityName() ([]string, error) {
	f, err := excelize.OpenFile("./areas.xlsx")
	if err != nil {
		panic(err)
	}

	// Get all the rows in the Sheet1.
	cols, err := f.GetCols("Sheet1")
	if err != nil {
		return nil, err

	}
	re, _ := regexp.Compile("市$")
	enterpriseNames := make([]string, 0, 100050)
	for idx, v := range cols[0] {
		if idx == 0 {
			continue
		}
		if re.MatchString(v) {

			enterpriseNames = append(enterpriseNames, v)
		}

	}

	return enterpriseNames, nil
}

func TypeByExtension(filename string) string {
	extension := mime.TypeByExtension(filepath.Ext(filename))
	if extension == "" {
		extension = "application/octet-stream"
	}
	return extension
}
