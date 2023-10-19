package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Alphonnse/file_server/internal/database"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// it holds the connection to a database
type ApiConfig struct {
	DB *database.Queries // its defined in db.go
}

func (apiCfg *ApiConfig) HandlerSigninTmpl(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/signin/signin.html")
}

func (apiCfg *ApiConfig) HandlerSignin(w http.ResponseWriter, r *http.Request) {

	// getting the body of response
	type signUpRecuestBody struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	body := signUpRecuestBody{}

	body.Name = r.FormValue("username")
	body.Email = r.FormValue("email")
	body.Password = r.FormValue("password")

	// decoder := json.NewDecoder(r.Body)
	// err := decoder.Decode(&body)
	// if err != nil {
	// 	RespondWithError(w, 400, fmt.Sprint("Error while parsing JSON:", err))
	// 	return
	// }

	// hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10) // last argument is the coast
	if err != nil {
		RespondWithError(w, 500, fmt.Sprint("Error while hashing the password:", err))
		return
	}

	// create the user
	_, err = apiCfg.DB.SignUp(r.Context(), database.SignUpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      body.Name,
		Email:     body.Email,
		Password:  string(hash),
	})
	if err != nil {
		RespondWithError(w, 500, fmt.Sprint("Error while creating the user:", err))
	}

	http.Redirect(w, r, "/login", 301)
}

// login
func (apiCfg *ApiConfig) HandlerLoginTmpl(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/login/login.html")
}

func (apiCfg *ApiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {

	// Get the email and pass pf required body
	type logInRequestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	body := logInRequestBody{}

	body.Email = r.FormValue("email")
	body.Password = r.FormValue("password")

	// decoder := json.NewDecoder(r.Body)
	// err := decoder.Decode(&body)
	// if err != nil {
	// 	RespondWithError(w, 400, fmt.Sprint("Error while parsing JSON:", err))
	// 	return
	// }

	// Look up requested user

	user, err := apiCfg.DB.LogIn(r.Context(), body.Email)
	if err != nil {
		RespondWithError(w, 400, "Invalid email or password")
		return
	}

	// Compare sent in pass with saved user pass hash

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		RespondWithError(w, 400, "Invalid email or password")
		return
	}

	// Generate a jwt token

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	//// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		ResondWithJSON(w, 500, fmt.Sprint("Failed to create token:", err))
	}

	// Send it back with cookie here
	cookie := http.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		MaxAge:   3600 * 24 * 30,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, fmt.Sprintf("/"+user.Name+"/disk/files/"), 301)
}
