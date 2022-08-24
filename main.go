package main

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"image"
	"image/jpeg"
	"image/png"
	"regexp"

	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"log"
)

// masuk ke model

func DBConn() (db *sql.DB, err error) {
	//dbDriver := "mysql"
	//dbUser := "devel"
	////dbPass := ""
	//dbName := "entry_task"

	//	db, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	//db, err = sql.Open(dbDriver, dbUser+"@/"+dbName)
	db, err = sql.Open("mysql", "devel:devel@tcp(127.0.0.1:3306)/entry_task")

	return
}

type User struct {
	Id       int    `json:"id"`
	UserName string `json:"userName"`
	NickName string `json:"name"`
	Picture  string `json:"picture"`
}

// mengembalikan pesan kesalahan/sukses seperti insert/delete dll ke clien
type Pesan struct {
	NamaPesan string
	KodePesan int
}

type Gambar struct {
	datagambar string `json:"DataGambar"`
}

func main() {
	fmt.Println("before router")
	router := httprouter.New()

	// GET
	router.GET("/getprofile", GetProfileFunc)
	router.GET("/getprofilepict", GetProfilePictFunc)

	//POST
	router.POST("/registeruser", RegisterUserFunc)
	router.GET("/login", LoginFunc)
	router.POST("/edituser", EditUserFunc)

	//PUT
	router.PUT("/uploadprofilepict", UploadProfilePictFunc)

	//	base64toJpg(getJPEGbase64("flower.jpg"))
	//base64toJpg(getJPEGbase64("a.png"))

	http.ListenAndServe(":8080", router)
}

func base64toPng(fUser string, fPicture string) {
	const data = `
/9j/4AAQSkZJRgABAQIAHAAcAAD/2wBDABALDA4MChAODQ4SERATGCgaGBYWGDEjJR0oOjM9PDkzODdA
SFxOQERXRTc4UG1RV19iZ2hnPk1xeXBkeFxlZ2P/2wBDARESEhgVGC8aGi9jQjhCY2NjY2NjY2NjY2Nj
Y2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2P/wAARCABnAJYDASIAAhEBAxEB/8QA
HwAAAQUBAQEBAQEAAAAAAAAAAAECAwQFBgcICQoL/8QAtRAAAgEDAwIEAwUFBAQAAAF9AQIDAAQRBRIh
MUEGE1FhByJxFDKBkaEII0KxwRVS0fAkM2JyggkKFhcYGRolJicoKSo0NTY3ODk6Q0RFRkdISUpTVFVW
V1hZWmNkZWZnaGlqc3R1dnd4eXqDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXG
x8jJytLT1NXW19jZ2uHi4+Tl5ufo6erx8vP09fb3+Pn6/8QAHwEAAwEBAQEBAQEBAQAAAAAAAAECAwQF
BgcICQoL/8QAtREAAgECBAQDBAcFBAQAAQJ3AAECAxEEBSExBhJBUQdhcRMiMoEIFEKRobHBCSMzUvAV
YnLRChYkNOEl8RcYGRomJygpKjU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6goOE
hYaHiImKkpOUlZaXmJmaoqOkpaanqKmqsrO0tba3uLm6wsPExcbHyMnK0tPU1dbX2Nna4uPk5ebn6Onq
8vP09fb3+Pn6/9oADAMBAAIRAxEAPwDlwKMD0pwzSiuK57QzGDxS7D6in8Y5ximnAPUfSlcq4m3ilUYp
2OKXHvRcVxnTtS7c07HNFK4DQPakC4PNOA+tOx70XAjK/So5gBGP94fzqfvUVx/qxx/EP51UXqRP4WSE
cmgjilP3jSEZqS0IO/NGDnpUiocDg/McDjvV6HTPOdVWYgsM5KcfzzQ2JySM2jp6VYu7SWzmMUwG4cgj
kMPUVBjjtTGtRu0Zopw+lFFxhinrGzuqqMsxAA9yaXFSRv5cqSEcIwYj6GpuZ30O30fSLKzhUpbpNMv3
5XGTn29BV28jt7pPLuIVljPBBFVreYx+VbqAjycgt3x14zRcNOxGyVFHQkIc/wA61exyKLbuzjdZ046d
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
`

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	//reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(fPicture))

	m, formatString, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	bounds := m.Bounds()
	fmt.Println(bounds, formatString)

	//Encode from image format to writer
	pngFilename := fUser + ".png"

	f, err := os.OpenFile(pngFilename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = png.Encode(f, m)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Png file", pngFilename, "created")

}

//??????????????????????????????????????????????????????????????????????????????????????????????????????????????????????

// Given a base64 string of a JPEG, encodes it into an JPEG image test.jpg
func base64toJpg(data string) {

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	m, formatString, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	bounds := m.Bounds()
	fmt.Println("base64toJpg", bounds, formatString)

	//Encode from image format to writer
	pngFilename := "test.jpg"
	f, err := os.OpenFile(pngFilename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = jpeg.Encode(f, m, &jpeg.Options{Quality: 75})
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Jpg file", pngFilename, "created")

}

// Gets base64 string of an existing JPEG file
// func getJPEGbase64(fileName string) string {
func fgetbase64(fileName string) string {

	var xx = fileName + ".png"
	imgFile, err := os.Open(xx)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer imgFile.Close()

	// create a new buffer base on file size
	fInfo, _ := imgFile.Stat()
	var size = fInfo.Size()
	buf := make([]byte, size)

	// read file content into buffer
	fReader := bufio.NewReader(imgFile)
	fReader.Read(buf)

	imgBase64Str := base64.StdEncoding.EncodeToString(buf)
	//fmt.Println("Base64 string is:", imgBase64Str)

	return imgBase64Str
}

//???????????????????????????????????????????????????????????????????????????????????????????????????????????????????

// ------------------------------------------------------------------------------
// function to register user
// end point http://localhost:8080/registeruser?username=xxxx&nickname=yyyyyy
// username (uniq key)
// ------------------------------------------------------------------------------
func RegisterUserFunc(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method == "POST" {
		db, err := DBConn()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer db.Close()

		var newUser User

		newUser.UserName = r.FormValue("username")
		newUser.NickName = r.FormValue("nickname")

		/*
			resBool, errStr := checkFormValue(w, r, xallStudents)
			if resBool == false {
				return false //???? perlu kode err/message ttg err , mis nama kosong dll???
			}
		*/

		// insert ke database (username uniq key, id sequance )
		_, err = db.Exec("insert into user(username,nickname)  values (?, ?) ", newUser.UserName, newUser.NickName)
		if err != nil {
			pesan := Pesan{"Gagal insert ", 0}
			jsonData, err := json.Marshal(pesan)
			if err != nil {
				//fmt.Println(err) // return
				//panic(err.Error())
				//http.Error(w, err.Error(), http.StatusInternalServerError) // return
				log.Fatal(err)
			}
			w.Write(jsonData)
			return
		}

		pesan := Pesan{"Sukses", 200}
		jsonData, err := json.Marshal(pesan)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(jsonData)
		return
	}

	http.Error(w, "harus menggunakan POST...........", http.StatusBadRequest)
}

// --------------------------------------------------------------------------------------------
// user login
// end point http://localhost:8080/login?username=xxxx
// jika user ada, return json (id,namauser dan nickname) sebaliknya null
// -----------------------------------------------------------------------------------------------
func LoginFunc(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	db, err := DBConn()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var newUser User
	var result []User

	//		newUser.UserName = r.FormValue("username")
	newUser.UserName = r.URL.Query().Get("username")

	rows, err := db.Query("SELECT id,username,name  FROM user WHERE  username=? LIMIT 1", newUser.UserName)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var each = User{}

		var err = rows.Scan(&each.Id, &each.UserName, &each.NickName)
		if err != nil {
			panic(err.Error())
			// fmt.Println(err.Error())
			// return
		}

		result = append(result, each)

	}
	// for _, each := range result {
	// 	fmt.Println(each.NickName)
	// }

	// jika data ditemukan, return data user, else  null
	// atau mau...	http.Error(w, "User not found", http.StatusNotFound) ??????

	var jsonData, errj = json.Marshal(result)
	if errj != nil {
		panic(err.Error())
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		//return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
	return
}

// -------------------------------------------------------------------------------------------------
// --Edit/update  nick name
// end point http://localhost:8080/edituser?username=xxxx&nicname=yyyyyy
// -------------------------------------------------------------------------------------------------
func EditUserFunc(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	//if r.Method == "POST" {
	db, err := DBConn()

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var newUser User

	//newUser.UserName = r.URL.Query().Get("username")
	newUser.UserName = r.FormValue("username")
	newUser.NickName = r.FormValue("nickname")

	//	newUser.UserName = r.PostFormValue("username")
	//	newUser.NickName = r.PostFormValue("nickname")
	// r.ParseForm   r.Formu sername]

	_, err = db.Exec("UPDATE user SET  nickname=? WHERE  username=?", newUser.NickName, newUser.UserName)
	if err != nil {
		panic(err.Error())
	}

	pesan := Pesan{"Sukses Update", 200}
	jsonData, err := json.Marshal(pesan)
	if err != nil {
		panic(err.Error())
	}
	w.Write(jsonData)
	//return

	//}
	//	http.Error(w, "harus menggunakan POST...........", http.StatusBadRequest)

}

// --------------------------------------------------------------------------------------------
// Get User Profile Func
// end point http://localhost:8080/userprofile?username=xxxx
// jika username tdk ditemukan, return null
// -----------------------------------------------------------------------------------------------
func GetProfileFunc(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {

		db, err := DBConn()

		if err != nil {
			//fmt.Println(err.Error())
			//return
			panic(err.Error())
		}
		defer db.Close()

		var newUser User

		//		newUser.UserName = r.FormValue("username")
		newUser.UserName = r.URL.Query().Get("username")

		//????		rows, err := db.Query("SELECT * FROM Employee WHERE id=? LIMIT 1", "a")
		rows, err := db.Query("SELECT id,username,nickname  FROM user WHERE  username=? LIMIT 1", newUser.UserName)
		if err != nil {
			panic(err.Error())
		}
		defer rows.Close()

		var result []User

		for rows.Next() {
			var each = User{}

			var err = rows.Scan(&each.Id, &each.UserName, &each.NickName)
			if err != nil {
				panic(err.Error())
			}

			result = append(result, each)
		}
		var jsonData, errj = json.Marshal(result)
		if errj != nil {
			fmt.Println(err.Error())
			return
		}

		w.Write(jsonData)
		return

	}

	http.Error(w, "", http.StatusBadRequest)
}

// --------------------------------------------------------------------------------------------
// UploadProfilePictFunc, pict stlh di decode di write/dikirim ke user dgn json
// file di simpan di server dgn nama file sesuai nama user
// Converts  base64 data to   png
// -----------------------------------------------------------------------------------------------
func UploadProfilePictFunc(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	/*
		db, err := DBConn()

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer db.Close()
	*/
	var newUser User

	newUser.UserName = r.FormValue("username")
	newUser.Picture = r.FormValue("picture")

	//sementara tdk perlu cek ke database apakah username tsb ada atau tidak

	base64toPng(newUser.UserName, newUser.Picture)

	return

	http.Error(w, "Harus Post", http.StatusBadRequest)
}

// --------------------------------------------------------------------------------------------
// GetProfilePictFunc, pict stlh di decode di write/dikirim ke user dgn json
//
// -----------------------------------------------------------------------------------------------
func GetProfilePictFunc(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {

		db, err := DBConn()

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer db.Close()

		var newUser User
		var newGambar Gambar

		newUser.UserName = r.FormValue("username")

		newGambar.datagambar = fgetbase64(newUser.UserName)

		var jsonData, errj = json.Marshal(newGambar.datagambar)
		if errj != nil {
			fmt.Println(err.Error())
			return
		}

		w.Write(jsonData)

		return

	}

	http.Error(w, "", http.StatusBadRequest)
}

// ---------------------------------------------------------------------------------------------------------
// function to check correct user adding input (regular expression and non-empty field input)
// ----------------------------------------------------------------------------------------------------------

func checkFormValue(w http.ResponseWriter, r *http.Request, forms ...string) (res bool, errStr string) {

	for _, form := range forms {
		m, _ := regexp.MatchString("^[a-zA-Z]+$", r.FormValue(form))
		if r.FormValue(form) == "" {
			return false, "All forms must be completed"
		}
		if m == false {
			return false, "Use only english letters if firstname,lastname forms"
		}

	}
	return true, ""
}
