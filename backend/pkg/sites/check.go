package sites

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/htchan/BookSpider/internal/utils"
)

func (site *Site) CheckEnd() {
	site.OpenDatabase()
	tx, err := site.database.Begin()
	utils.CheckError(err)
	matchingChapterCriteria := []string{"后记", "後記", "新书", "新書", "结局", "結局", "感言",
		"尾声", "尾聲", "终章", "終章", "外传", "外傳", "完本", "结束", "結束", "完結",
		"完结", "终结", "終結", "番外", "结尾", "結尾", "全书完", "全書完", "全本完"}
	sqlStmt := "update books set end=true, download=false where ("
	for _, criteria := range matchingChapterCriteria {
		sqlStmt += "chapter like '%" + criteria + "%' or "
	}
	sqlStmt += "date < '" + strconv.Itoa(time.Now().Year()-1) + "') and (end <> true or end is null)"
	result, err := tx.Exec(sqlStmt)
	utils.CheckError(err)
	rowAffect, err := result.RowsAffected()
	utils.CheckError(err)
	utils.CheckError(tx.Commit())
	log.Println(site.SiteName, "Row affected: ", rowAffect)
	site.CloseDatabase()
}

func (site *Site) checkDuplicateBook() {
	// check duplicate record in books table
	rows, tx := site.query("select num, version from books group by num, version having count(*) > 1")
	log.Print(site.SiteName, "duplicate books id : [")
	var id, version, count int
	for ; rows.Next(); count++ {
		if count > 0 {
			log.Println(", ")
		}
		rows.Scan(&id, &version)
		log.Print("(" + strconv.Itoa(id) + ", " + strconv.Itoa(version) + ")")
	}
	log.Println("]")
	log.Println(site.SiteName, "duplicate books count : "+strconv.Itoa(count))
	closeQuery(rows, tx)
}

func (site *Site) checkDuplicateError() {
	// check duplicate record in error table
	rows, tx := site.query("select num from error group by num having count(*) > 1")
	log.Print(site.SiteName, "duplicate error id : [")
	var id, count int
	for ; rows.Next(); count++ {
		if count > 0 {
			log.Print(", ")
		}
		rows.Scan(&id)
		log.Print(strconv.Itoa(id))
	}
	log.Println("]")
	log.Println(site.SiteName, "duplicate error count : "+strconv.Itoa(count))
	closeQuery(rows, tx)
}

func (site *Site) checkDuplicateCrossTable() {
	// check duplicate record crossing books and error table
	rows, tx := site.query("select distinct num from books where num in (select num from error)")
	log.Print(site.SiteName, "duplicate cross id : [")
	var id, count int
	for ; rows.Next(); count++ {
		if count > 0 {
			log.Print(", ")
		}
		rows.Scan(&id)
		log.Print(strconv.Itoa(id))
	}
	log.Println("]")
	log.Println(site.SiteName, "duplicate cross count : "+strconv.Itoa(count))
	closeQuery(rows, tx)
}

func (site *Site) checkMissingId() {
	missingBookIds := site.missingIds()
	jsonByte, err := json.Marshal(missingBookIds)
	utils.CheckError(err)
	log.Println(site.SiteName, "missing id : ", string(jsonByte))
	log.Println(site.SiteName, "missing count : "+strconv.Itoa(len(missingBookIds)))
}

func (site *Site) Check() {
	// init variable
	site.OpenDatabase()
	site.checkDuplicateBook()
	site.checkDuplicateError()
	site.checkDuplicateCrossTable()

	// check missing record
	site.checkMissingId()
	site.CloseDatabase()
}
