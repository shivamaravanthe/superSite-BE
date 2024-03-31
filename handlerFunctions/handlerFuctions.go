package handlerFunctions

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	cryptopwd "github.com/shivamaravanthe/superSite-BE/crypto"
	"github.com/shivamaravanthe/superSite-BE/dto"
	"github.com/shivamaravanthe/superSite-BE/jwt"
	"github.com/shivamaravanthe/superSite-BE/model"
	"github.com/shivamaravanthe/superSite-BE/store"
	"github.com/shivamaravanthe/superSite-BE/utils"
	"golang.org/x/crypto/bcrypt"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	ok := "BKtZEdc/6mpkiCid1VITHHqv67Q0OcKJkZgssZ6eN0J+KgajY7N6qRdQOYR9oq61JA9Zsoyscm5nua+s2V8R2EIHlKFjMk9arhSDWDVHYCW9ZjHQsGsTjl/tRysgEDWC85JOOtM="
	// bytes, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	fmt.Printf("Errored Reading Post Body: %v\n", err)
	// 	utils.CreateErrorResponse(w, 400, "Errored Reading Post Body", err)
	// 	return
	// }
	// // fmt.Printf("Body: %+v", r.Body)
	// r.Body.Close()
	// fmt.Printf("Body: %+v\n", string(bytes))
	utils.CreateResponse(w, "Pong", &ok)
}

func SingUp(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Errored Reading Post Body: %v\n", err)
		utils.CreateErrorResponse(w, 400, "Errored Reading Post Body", err)
		return
	}

	user := model.Users{}
	if err = json.Unmarshal(bytes, &user); err != nil {
		fmt.Printf("Errored Unmarshalling Post Body: %v\n", err.Error())
		utils.CreateErrorResponse(w, 400, "Errored Unmarshalling Post Body", err)
		return
	}

	if user.Password == "" {
		fmt.Println("Password is Missing")
		utils.CreateErrorResponse(w, 403, "Password is Missing", nil)
		return
	}

	if user.Email == "" {
		fmt.Println("Email is Missing")
		utils.CreateErrorResponse(w, 403, "Email is Missing", nil)
		return
	}
	pwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		fmt.Printf("Errored hashing password %v\n", err.Error())
		utils.CreateErrorResponse(w, 400, "Errored hashing password", err)
		return
	}

	if err := store.DB.Create(&model.Users{
		Email:     user.Email,
		Password:  string(pwd),
		CreatedAt: time.Now(),
		CreatedBy: "Admin",
		UpdatedAt: time.Now(),
		UpdatedBy: "Admin",
	}).Error; err != nil {
		fmt.Printf("Errored Adding to Users Table %v\n", err.Error())
		utils.CreateErrorResponse(w, 400, "Errored Adding to Users Table", err)
		return
	}

	utils.CreateResponse(w, "OK", &utils.Empty{})
}

func Login(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Errored Reading Post Body: %v\n", err)
		utils.CreateErrorResponse(w, 400, "Errored Reading Post Body", err)
		return
	}

	user := dto.Users{}
	if err = json.Unmarshal(bytes, &user); err != nil {
		fmt.Printf("Errored Unmarshalling Post Body: %v\n", err.Error())
		utils.CreateErrorResponse(w, 400, "Errored Unmarshalling Post Body", err)
		return
	}

	if user.Password == "" {
		fmt.Println("Password is Missing")
		utils.CreateErrorResponse(w, 400, "Password is Missing", nil)
		return
	}
	if user.Email == "" {
		fmt.Println("Email is Missing")
		utils.CreateErrorResponse(w, 400, "Email is Missing", nil)
		return
	}

	userDto := model.Users{}
	if err := store.DB.First(&userDto, "Email = ?", user.Email).Error; err != nil {
		fmt.Printf("Errored Selecting user data %v\n", err.Error())
		utils.CreateErrorResponse(w, 401, "User Not Found", err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userDto.Password), []byte(user.Password)); err != nil {
		fmt.Printf("Password not Valid %v\n", err.Error())
		utils.CreateErrorResponse(w, 401, "Password not Valid ", err)
		return
	}

	tokenString, err := jwt.CreateToken(userDto.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Error creating token %v", err)
		return
	}
	w.Header().Add("token", tokenString)
	w.Header().Add("Access-Control-Expose-Headers", "token")
	expire := time.Now().Add(6 * time.Minute)
	cookie := http.Cookie{Name: "Auth", Value: tokenString, Path: "/", Expires: expire, MaxAge: 90000}
	http.SetCookie(w, &cookie)
	utils.CreateResponse(w, "Ok", &utils.Empty{})
}

func StorePassword(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Errored Reading Post Body: %v\n", err)
		utils.CreateErrorResponse(w, 400, "Errored Reading Post Body", err)
		return
	}

	pwd := dto.PasswordStore{}
	if err = json.Unmarshal(bytes, &pwd); err != nil {
		fmt.Printf("Errored Unmarshalling Post Body: %v\n", err.Error())
		utils.CreateErrorResponse(w, 400, "Errored Unmarshalling Post Body", err)
		return
	}

	if pwd.Password == "" {
		fmt.Println("Password is Missing")
		utils.CreateErrorResponse(w, 400, "Password is Missing", nil)
		return
	}

	if pwd.UserName == "" {
		fmt.Println("UserName is Missing")
		utils.CreateErrorResponse(w, 400, "UserName is Missing", nil)
		return
	}

	if pwd.Description == "" {
		fmt.Println("Description is Missing")
		utils.CreateErrorResponse(w, 400, "Description is Missing", nil)
		return
	}

	if pwd.Master == "" {
		fmt.Println("Master is Missing")
		utils.CreateErrorResponse(w, 400, "Key is Missing", nil)
		return
	}

	encUserName, err := cryptopwd.Encrypt([]byte(pwd.UserName), pwd.Master)
	if err != nil {
		fmt.Printf("Errored Encrypt UserName: %v\n", err)
		utils.CreateErrorResponse(w, 400, "Errored Encrypt UserName:", err)
		return
	}

	encLink, err := cryptopwd.Encrypt([]byte(pwd.Link), pwd.Master)
	if err != nil {
		fmt.Printf("Errored Encrypt Link: %v\n", err)
		utils.CreateErrorResponse(w, 400, "Errored Encrypt Link:", err)
		return
	}

	encPwd, err := cryptopwd.Encrypt([]byte(pwd.Password), pwd.Master)
	if err != nil {
		fmt.Printf("Errored Encrypt Link: %v\n", err)
		utils.CreateErrorResponse(w, 400, "Errored Encrypt Link:", err)
		return
	}

	mdl := model.PasswordStore{
		Link:        string(encLink),
		UserName:    string(encUserName),
		Password:    string(encPwd),
		Description: pwd.Description,
		CreatedAt:   time.Now(),
		CreatedBy:   "Admin",
		UpdatedAt:   time.Now(),
		UpdatedBy:   "Admin",
	}

	if err := store.DB.Create(&mdl).Error; err != nil {
		fmt.Printf("Errored Adding to PasswordStore Table %v\n", err.Error())
		utils.CreateErrorResponse(w, 400, "Errored Adding to PasswordStore Table", err)
		return
	}

	stars := strings.Repeat("*", 8)

	resp := dto.PasswordStore{
		ID:          int(mdl.ID),
		Link:        stars,
		UserName:    stars,
		Password:    stars,
		Description: mdl.Description,
		Master:      "",
	}

	utils.CreateResponse(w, "Ok", &resp)
}

func GetPassword(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Errored Reading Post Body: %v\n", err)
		utils.CreateErrorResponse(w, 400, "Errored Reading Post Body", err)
		return
	}

	pwd := dto.PasswordStore{}
	if err = json.Unmarshal(bytes, &pwd); err != nil {
		fmt.Printf("Errored Unmarshalling Post Body: %v\n", err.Error())
		utils.CreateErrorResponse(w, 400, "Errored Unmarshalling Post Body", err)
		return
	}

	if pwd.ID < 1 {
		fmt.Println("ID is Missing")
		utils.CreateErrorResponse(w, 400, "ID is Missing", nil)
		return
	}

	pwdDto := model.PasswordStore{}
	if err := store.DB.First(&pwdDto, "ID = ?", pwd.ID).Error; err != nil {
		fmt.Printf("Errored Selecting user data %v\n", err.Error())
		utils.CreateErrorResponse(w, 400, "Selecting user data", err)
		return
	}

	password, err := cryptopwd.Decrypt([]byte(pwdDto.Password), pwd.Master)
	if err != nil {
		fmt.Printf("Errored Decrypt Password: %v\n", err)
		utils.CreateErrorResponse(w, 400, "Invalid Key", err)
		return
	}

	userName, err := cryptopwd.Decrypt([]byte(pwdDto.UserName), pwd.Master)
	if err != nil {
		fmt.Printf("Errored Decrypt userName: %v\n", err)
		utils.CreateErrorResponse(w, 400, "Invalid Key", err)
		return
	}

	link, err := cryptopwd.Decrypt([]byte(pwdDto.Link), pwd.Master)
	if err != nil {
		fmt.Printf("Errored Decrypt link: %v\n", err)
		utils.CreateErrorResponse(w, 400, "Invalid Key", err)
		return
	}

	dto := dto.PasswordStore{
		ID:          pwd.ID,
		Link:        string(link),
		UserName:    string(userName),
		Password:    string(password),
		Description: pwdDto.Description,
	}

	utils.CreateResponse(w, "Ok", &dto)
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Errored Reading Post Body: %v\n", err)
		utils.CreateErrorResponse(w, 400, "Errored Reading Post Body", err)
		return
	}

	pwd := dto.PasswordStore{}
	if err = json.Unmarshal(bytes, &pwd); err != nil {
		fmt.Printf("Errored Unmarshalling Post Body: %v\n", err.Error())
		utils.CreateErrorResponse(w, 400, "Errored Unmarshalling Post Body", err)
		return
	}

	if pwd.ID <= 0 {
		fmt.Println("ID is Missing")
		utils.CreateErrorResponse(w, 400, "ID is Missing", nil)
		return
	}
	if pwd.Password == "" {
		fmt.Println("Password is Missing")
		utils.CreateErrorResponse(w, 400, "Password is Missing", nil)
		return
	}

	if pwd.UserName == "" {
		fmt.Println("UserName is Missing")
		utils.CreateErrorResponse(w, 400, "UserName is Missing", nil)
		return
	}

	if pwd.Description == "" {
		fmt.Println("Description is Missing")
		utils.CreateErrorResponse(w, 400, "Description is Missing", nil)
		return
	}

	if pwd.Master == "" {
		fmt.Println("Master is Missing")
		utils.CreateErrorResponse(w, 400, "Key is Missing", nil)
		return
	}

	pwdDto := model.PasswordStore{}
	if err := store.DB.First(&pwdDto, "ID = ?", pwd.ID).Error; err != nil {
		fmt.Printf("Errored Selecting user data %v\n", err.Error())
		utils.CreateErrorResponse(w, 400, "Selecting user data", err)
		return
	}
	if pwdDto.Password == "" {
		fmt.Println("Wrong ID")
		utils.CreateErrorResponse(w, 400, "Wrong ID", nil)
		return
	}
	_, err = cryptopwd.Decrypt([]byte(pwdDto.Password), pwd.Master)
	if err != nil {
		fmt.Printf("Errored Decrypt Password: %v\n", err)
		utils.CreateErrorResponse(w, 400, "Invalid Key", err)
		return
	}

	encUserName, err := cryptopwd.Encrypt([]byte(pwd.UserName), pwd.Master)
	if err != nil {
		fmt.Printf("Errored Encrypt UserName: %v\n", err)
		utils.CreateErrorResponse(w, 400, "Errored Encrypt UserName:", err)
		return
	}

	encLink, err := cryptopwd.Encrypt([]byte(pwd.Link), pwd.Master)
	if err != nil {
		fmt.Printf("Errored Encrypt Link: %v\n", err)
		utils.CreateErrorResponse(w, 400, "Errored Encrypt Link:", err)
		return
	}

	encPwd, err := cryptopwd.Encrypt([]byte(pwd.Password), pwd.Master)
	if err != nil {
		fmt.Printf("Errored Encrypt Link: %v\n", err)
		utils.CreateErrorResponse(w, 400, "Errored Encrypt Link:", err)
		return
	}

	mdl := model.PasswordStore{
		ID:          uint(pwd.ID),
		Link:        string(encLink),
		UserName:    string(encUserName),
		Password:    string(encPwd),
		Description: pwd.Description,
		CreatedAt:   time.Now(),
		CreatedBy:   "Admin",
		UpdatedAt:   time.Now(),
		UpdatedBy:   "Admin",
	}

	if err := store.DB.Save(&mdl).Error; err != nil {
		fmt.Printf("Errored Updating to PasswordStore Table %v\n", err.Error())
		utils.CreateErrorResponse(w, 400, "Errored Updating to PasswordStore Table", err)
		return
	}
	utils.CreateResponse(w, "Ok", &mdl)
}

func ListPassword(w http.ResponseWriter, r *http.Request) {

	pwdDto := []model.PasswordStore{}
	if err := store.DB.Select("ID,Description").Table("password_stores").Scan(&pwdDto).Error; err != nil {
		fmt.Printf("Errored Selecting user data %v\n", err.Error())
		utils.CreateErrorResponse(w, 400, "Selecting user data", err)
		return
	}

	resp := []dto.PasswordStore{}
	stars := strings.Repeat("*", 8)
	for _, each := range pwdDto {
		resp = append(resp, dto.PasswordStore{
			ID:          int(each.ID),
			Link:        stars,
			UserName:    stars,
			Password:    stars,
			Description: each.Description,
			Master:      "",
		})
	}

	utils.CreateResponse(w, "Ok", &resp)
}
