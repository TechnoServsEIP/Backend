package models

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"gitlab.sysroot.ovh/technoservs/microservices/game-servers/utils"
	"golang.org/x/crypto/bcrypt"
)

//JWT claims struct
type Token struct {
	UserId uint
	jwt.StandardClaims
}

//a struct to rep user account
type Account struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
	Verified bool
}

func (account Account) generateJWT() (string, error) {
	fmt.Println("attribute user number ", account.ID)
	tk := &Token{
		UserId: account.ID,
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

	url := "http://localhost:9096/user/confirm?token=" + tokenString

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
