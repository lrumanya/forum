package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/inscription", inscription)
	http.ListenAndServe(":9000", nil)
	log.Println("Serveur démarré sur http://localhost:9000")
}

func index(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index")

}

var db *sql.DB

func init() {
	// Initialise la connexion à la base de données (à adapter selon votre base de données)
	db, err := sql.Open("sqlite3", "digibook.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Vérifie la connexion à la base de données
	err = db.Ping()
	if err != nil {
		panic(err)
	}
}

func inscription(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "inscription")
	if r.Method == http.MethodPost {
		// Récupère les données du formulaire
		username := r.FormValue("pseudo")
		email := r.FormValue("email")
		datenaissance := r.FormValue("datenaissance")
		photo := r.FormValue("photo")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirmpassword")
		accepte_termes := r.FormValue("accepte_termes")

		// Vérifie que le mot de passe et la confirmation du mot de passe correspondent
		if password != confirmPassword {
			http.Error(w, "Le mot de passe et la confirmation du mot de passe ne correspondent pas", http.StatusBadRequest)
			return
		}

		// Utilise bcrypt pour hacher le mot de passe
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Erreur lors du hachage du mot de passe", http.StatusInternalServerError)
			return
		}

		// Stocke les informations de l'utilisateur dans la base de données
		err = insertUser(username, email, datenaissance, photo, string(hashedPassword), accepte_termes)
		if err != nil {
			http.Error(w, "Erreur lors de l'inscription de l'utilisateur", http.StatusInternalServerError)
			return
		}

		// Redirige vers une page de succès
		http.Redirect(w, r, "/connexion", http.StatusSeeOther)
		return
	}
}

func insertUser(username, email, datenaissance, photo, hashedPassword, accepte_termes string) error {
	// Insère l'utilisateur dans la base de données (à adapter selon votre schéma de base de données)
	_, err := db.Exec("INSERT INTO user (username, email, datenaissance, photo, Password, accepte_termes) VALUES (?, ?, ?, ?, ?, ?)", username, email, datenaissance, photo, hashedPassword, accepte_termes)
	return err
}

func renderTemplate(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles("./template/" + tmpl + ".html")
	if err != nil {
		http.Error(w, "Erreur lors du lancement de la page accueil", http.StatusInternalServerError)
		return
	}

	t.Execute(w, nil)
}
