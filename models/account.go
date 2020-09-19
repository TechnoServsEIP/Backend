package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/TechnoServsEIP/Backend/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

//JWT claims struct
type Token struct {
	UserId uint
	Role   string
	jwt.StandardClaims
}

//a struct to rep user account
type Account struct {
	gorm.Model
	Role     string
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
	Verified bool
	Activate bool
}

type GithubData struct {
	Login string `json:"login"`
	Id    int    `json:"id"`
}

func (account Account) generateJWT() (string, error) {
	fmt.Println("attribute user number ", account.ID)
	tk := &Token{
		UserId: account.ID,
		Role:   account.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1 * 1).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
	tokenString, err := token.SignedString([]byte(os.Getenv("token_password")))
	return tokenString, err
}

//Validate incoming user details...
func (account *Account) Validate() (map[string]interface{}, bool) {
	if !strings.Contains(account.Email, "@") {
		return utils.Message(false, "Email address is required"), false
	}

	if len(account.Password) == 0 {
		return utils.Message(false, "Password is required"), false
	}

	if len(account.Password) < 6 {
		return utils.Message(false, "Password len need to be 6 characters minimum"), false
	}

	//Email must be unique
	temp := &Account{}

	//check for errors and duplicate emails
	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return utils.Message(false, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return utils.Message(false, "Email address already in use by another user."), false
	}

	return utils.Message(false, "Requirement passed"), true
}

func (account *Account) Create() map[string]interface{} {
	if resp, ok := account.Validate(); !ok {
		fmt.Println("account validate resp:", resp)
		return resp
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)
	account.Role = "user"
	account.Activate = true

	//if account.ID <= 0 {
	//	return utils.Message(false, "Failed to create account, connection error.")
	//}
	//Create new JWT token for the newly registered account
	GetDB().Create(account)
	tokenString, err := account.generateJWT()
	if err != nil {
		return utils.Message(false, "Failed to create account")
	}
	account.Token = tokenString

	url := "https://technoservs.ichbinkour.eu/confirm?token=" + tokenString

	err = utils.SendConfirmationEmail(url, account.Email)
	if err != nil {
		return utils.Message(false, "Failed to send email")
	}

	account.Password = "" //delete password

	response := utils.Message(true, "Account has been created")
	response["account"] = account
	return response
}

func Login(email, password string) map[string]interface{} {
	account := &Account{}
	err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error
	if !account.Activate {
		return utils.Message(false, "account deacivate")
	}
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.Message(false, "Email address not found")
		}
		return utils.Message(false, "Connection error. Please retry")
	}

	if !account.Verified {
		return utils.Message(false, "Email address not verified")
	}
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return utils.Message(false, "Invalid login credentials. Please try again")
	}
	account.Password = ""

	//Create JWT token
	tokenString, err := account.generateJWT()
	if err != nil {
		return utils.Message(false, "Failed to create account")
	}
	account.Token = tokenString //Store the token in the response

	resp := utils.Message(true, "Logged In")
	resp["account"] = account
	return resp
}

func GetUserFromId(Id int) *Account {
	acc := &Account{}

	err := GetDB().Table("accounts").Where("id = ?", Id).First(acc).Error
	if err != nil {
		fmt.Println("error fetching user ", err)
		return nil
	}
	fmt.Println(acc)
	if acc.Email == "" { //User not found!
		return nil
	}

	acc.Password = ""
	return acc
}

func Update(Id int, fieldsToUpdate map[string]interface{}) *Account {
	acc := &Account{}

	GetDB().Table("accounts").Where("id = ?", Id).Update(fieldsToUpdate).First(acc)

	if acc.Email == "" { //User not found!
		return nil
	}

	acc.Password = ""
	return acc
}

func GetUsers() interface{} {
	users := []Account{}
	res := GetDB().Find(&users)
	return res
}

func DeactivateUser(id int) interface{} {
	user := GetUserFromId(id)
	if user == nil {
		return nil
	}
	user.Activate = false
	res := GetDB().Save(&user)
	return res
}

func ActivateUser(id int) interface{} {
	user := GetUserFromId(id)
	if user == nil {
		return nil
	}
	user.Activate = true
	res := GetDB().Save(&user)
	return res
}

func VerifyUser(id int) interface{} {
	user := GetUserFromId(id)
	if user == nil {
		return nil
	}
	user.Verified = true
	res := GetDB().Save(&user)
	return res
}

func RemoveVerification(id int) interface{} {
	user := GetUserFromId(id)
	if user == nil {
		return nil
	}
	user.Verified = false
	res := GetDB().Save(&user)
	return res
}

func ChangePassword(password string, id uint) error {
	account := GetUserFromId(int(id))
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	account.Password = string(hashedPassword)
	GetDB().Save(&account)
	return nil
}

// OAuth2

func DecryptToken(tokenString string) (jwt.Claims, bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Token{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})
	if err != nil {
		fmt.Println("Malformed authentication token ", err)
		return token.Claims, token.Valid, err
	}

	fmt.Println(token.Claims.(*Token).UserId)
	return token.Claims, token.Valid, nil
}

func AuthenticateOAuth2User(account *Account) map[string]interface{} {
	passwordTmp := account.Password
	if _, newUser := account.Validate(); !newUser {
		return Login(account.Email, account.Password)
	} else {
		response := account.Create()

		jsonString, _ := json.Marshal(response["account"])
		r := bytes.NewReader(jsonString)
		json.NewDecoder(r).Decode(account)

		claims, _, _ := DecryptToken(account.Token)

		user := GetUserFromId(int(claims.(*Token).UserId))
		user.Verified = true
		Update(int(user.ID), map[string]interface{}{
			"verified": true,
		})

		return Login(account.Email, passwordTmp)
	}
}

func LoginGithub(email string, password string) map[string]interface{} {
	account := &Account{}

	account.Email = email
	account.Password = password

	return AuthenticateOAuth2User(account)
}
