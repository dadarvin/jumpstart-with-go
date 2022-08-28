package main

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang/gddo/httputil/header"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

// DBConn make connection to database
func DBConn() (db *sql.DB, err error) {
	//	db, err = sql.Open("mysql", "devel:devel@tcp(127.0.0.1:3306)/entry_task")
	db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/entry_task")

	return
}

// User struct for user data
type User struct {
	Id       int    `json:"id" form:"id"`
	UserName string `json:"username" form:"username"`
	NickName string `json:"nickname" form:"nickname"`
	Password string `json:"password" form:"password"`
}

// UserPicture struct for containing basecode64 from client to convert into base64 format
type UserPicture struct {
	Id             int    `json:"id"`
	UserName       string `json:"nickname"`
	ProfilePicture string `json:"profilepicture"`
}

// JsonResponse struct for JSON response
type JsonResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// JSONErrorResponse struct for JSON error response
type JsonErrorResponse struct {
	Error *ApiResponse `json:"error"`
}

// ApiResponse response for Api request
type ApiResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func main() {
	router := httprouter.New()

	// GET
	router.GET("/get-profile/:id", GetProfileFunc)
	router.GET("/get-profile-pict/:id", GetProfilePictFunc)

	//POST
	router.POST("/register", RegisterUserFunc)
	router.POST("/login", LoginFunc)
	router.POST("/uploadprofilepict", UploadProfilePictFunc)

	//PUT
	router.PUT("/change-nickname", EditUserFunc)

	fmt.Println("Running Server in :8080......")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func checkPostHeader(r *http.Request) string {
	msg := ""
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg = "Content-Type header is not application/json"
			return msg
		}
	}
	return msg
}

func httpResponse(w http.ResponseWriter, responseCode int, message string, m interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(responseCode)

	if err := json.NewEncoder(w).Encode(&JsonResponse{Status: responseCode, Message: message, Data: m}); err != nil {
		errorResponse(w, http.StatusInternalServerError, "Internal Server Error")
	}
}

func errorResponse(w http.ResponseWriter, errorCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(errorCode)
	json.NewEncoder(w).Encode(&JsonErrorResponse{Error: &ApiResponse{Status: errorCode, Message: errorMsg}})
}

// RegisterUserFunc register new user
// @Router /register [POST]
// @Param username body string
// @Param nickname body string
func RegisterUserFunc(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var newUser User

	checkPost := checkPostHeader(r)
	if checkPost != "" {
		errorResponse(w, http.StatusUnsupportedMediaType, checkPost)
		return
	}

	bodyVal, err := io.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = json.Unmarshal(bodyVal, &newUser)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// hash password with bcrypt
	hashedPass, err := HashPassword(newUser.Password)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// connecting to database
	db, err := DBConn()
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	defer db.Close()

	// insert ke database (username uniq key, id sequance )
	_, err = db.Exec("insert into user(username, nickname, password)  values (?, ?, ?) ", newUser.UserName, newUser.NickName, hashedPass)
	if err != nil {
		//http.Error(w, "failed inserting data : "+err.Error(), http.StatusInternalServerError)
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	httpResponse(w, http.StatusOK, "Register Sukses", nil)
}

// LoginFunc endpoint for user login
// @Routes /login
// @Param id body int
// @Param password body string
// jika user ada, return json (id,namauser dan nickname) sebaliknya null
func LoginFunc(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var newUser User

	checkPost := checkPostHeader(r)
	if checkPost != "" {
		//		http.Error(w, checkPost, http.StatusUnsupportedMediaType)
		errorResponse(w, http.StatusUnsupportedMediaType, " Harus menggunakan POST")
		return
	}

	bodyVal, err := io.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = json.Unmarshal(bodyVal, &newUser)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// connect ke database
	db, err := DBConn()
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	defer db.Close()

	// var result []User

	rows, err := db.Query("SELECT id, username, nickname, password FROM user WHERE username=? LIMIT 1", newUser.UserName)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return

	}
	defer rows.Close()
	//	fmt.Println(rows)

	//	var each = User{}
	//	each := User{}
	var each User
	for rows.Next() {
		var err = rows.Scan(&each.Id, &each.UserName, &each.NickName, &each.Password)
		if err != nil {
			errorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// result = append(result, each)
	}

	if each.UserName == "" {
		httpResponse(w, http.StatusOK, "Username tidak ditemukan", nil)
		return
	}

	match := CheckPasswordHash(newUser.Password, each.Password)
	if !match {
		httpResponse(w, http.StatusOK, "Password Salah", nil)
		return
	}

	// UserData := UserData{
	// 	Id: each.Id,
	// 	UserName: each.UserName,
	// 	NickName: each.NickName,
	// }

	// Return data user tanpa field password
	userData := map[string]interface{}{
		"id":       each.Id,
		"username": each.UserName,
		"nickname": each.NickName,
	}
	httpResponse(w, http.StatusOK, "Login Sukses", userData)
	/*
		// for _, each := range result {
		// 	fmt.Println(each.NickName)
		// }
		var jsonData, errj = json.Marshal(result)
		if errj != nil {
			panic(err.Error())
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
		return
	*/
}

// EditUserFunc Edit user nickname
// @Router /edit-user
// -------------------------------------------------------------------------------------------------
func EditUserFunc(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var newUser User

	bodyVal, err := io.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = json.Unmarshal(bodyVal, &newUser)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	db, err := DBConn()
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	defer db.Close()

	//newUser.UserName = r.URL.Query().Get("username")
	//newUser.UserName = r.FormValue("username")
	//newUser.NickName = r.FormValue("nickname")
	//	newUser.UserName = r.PostFormValue("username")
	//	newUser.NickName = r.PostFormValue("nickname")
	// r.ParseForm   r.Formu sername]

	_, err = db.Exec("UPDATE user SET nickname=? WHERE id=?", newUser.NickName, newUser.Id)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	httpResponse(w, http.StatusOK, "Edit User Sukses", newUser)

	//pesan := Pesan{"Sukses Update", 200}
	//jsonData, err := json.Marshal(pesan)
	//if err != nil {
	//	panic(err.Error())
	//}
	//w.Write(jsonData)
	//return

	//w.Header().Set("Content-Type", "application/json")
}

// GetProfileFunc Get Profile info
// @Routes /user-profile/:id
// @Params id queryParam int
// jika username tdk ditemukan, return null
// -----------------------------------------------------------------------------------------------
func GetProfileFunc(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	var (
		newUser User
		err     error
	)
	/*
		checkPost := checkPostHeader(r)
		if checkPost != "" {
			http.Error(w, checkPost, http.StatusUnsupportedMediaType)
		}

		bodyVal, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("error", err)
		}

		err = json.Unmarshal(bodyVal, &newUser)
		if err != nil {
			fmt.Println("error unmarshaling")
		}
	*/

	//		newUser.UserName = r.FormValue("username")
	// newUser.UserName = r.URL.Query().Get("id")
	idString := param.ByName("id")
	newUser.Id, err = strconv.Atoi(idString)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	db, err := DBConn()
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	defer db.Close()

	//????		rows, err := db.Query("SELECT * FROM Employee WHERE id=? LIMIT 1", "a")
	rows, err := db.Query("SELECT id, username, nickname FROM user WHERE id=? LIMIT 1", newUser.Id)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	defer rows.Close()

	var result User

	for rows.Next() {
		// var each = User{}

		var err = rows.Scan(&result.Id, &result.UserName, &result.NickName)
		if err != nil {
			errorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// result = append(result, each)
	}
	httpResponse(w, http.StatusOK, "Sukses", result)

	/*
		var jsonData, errj = json.Marshal(result)
		if errj != nil {
			fmt.Println(err.Error())
			return
		}

		w.Write(jsonData)
		return
	*/
}

// ----------------------------------
// Picture Logic
// ----------------------------------

// UploadProfilePictFunc pict stlh di decode di write/dikirim ke user dgn json
// file di simpan di server dgn nama file sesuai nama user
// Converts  base64 data to   png
func UploadProfilePictFunc(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	/*
		db, err := DBConn()

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer db.Close()
	*/
	var (
		newUser UserPicture
		err     error
	)

	checkPost := checkPostHeader(r)
	if checkPost != "" {
		http.Error(w, checkPost, http.StatusUnsupportedMediaType)
		return
	}

	bodyVal, err := io.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = json.Unmarshal(bodyVal, &newUser)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	//sementara tdk perlu cek ke database apakah username tsb ada atau tidak

	//dikirim ke function base64Png untuk diconver kt file png dan disimpan di disk file dgn Id user sebagai filename.
	base64toPng(strconv.Itoa(newUser.Id), newUser.ProfilePicture)

	//return

}

// GetProfilePictFunc, pict stlh di decode di write/dikirim ke user dgn json
// @Routes /get-profile-pict/:id
// @Params id parameter id user int
// jika username tdk ditemukan, return null
func GetProfilePictFunc(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

	var (
		newUser UserPicture
		err     error
	)
	// var newGambar [] Gambar
	//	var newGambar []string

	// newUser.UserName = r.URL.Query().Get("username")
	idString := param.ByName("id")
	newUser.Id, err = strconv.Atoi(idString)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	//fgetbase64= fungsi utk menkonversi disk file ke basecode 64 dan dikirim ke clien
	//		newGambar.datagambar = fgetbase64(newUser.UserName)
	//	newGambar[0] = fgetbase64(newUser.UserName)
	//newUser.ProfilePicture = fgetbase64(newUser.UserName)

	httpResponse(w, http.StatusOK, "Sukses Get Profile Picture", fgetbase64(strconv.Itoa(newUser.Id)))
	/*
		w.Header().Set("Content-Type", "application/json")

		var jsonData, errj = json.Marshal(newGambar.datagambar)
		if errj != nil {
			fmt.Println(err.Error())
			return
		}

		w.Write(jsonData)

		return
	*/

}

// fungsi
func base64toPng(fIdUser string, fPicture string) {
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

	//	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	//fPicture adalah base64cie yg dikirim dari clien utk diubah jadi png
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(fPicture))

	m, formatString, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	bounds := m.Bounds()
	fmt.Println(bounds, formatString)

	//Encode from image format to writer
	//fUser=nama user yg dijadinakan file name (nama user unique)
	pngFilename := "Pict_" + fIdUser + ".png"

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

// fgetbase64 Gets base64 string of an existing JPEG file
// fungsi utk mengamfile file berdasarkan nama user, utk diconversi kebase64cide dan dikirim ke clien
func fgetbase64(fileName string) string {

	var filename = "Pict_" + fileName + ".png"
	imgFile, err := os.Open(filename)

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
