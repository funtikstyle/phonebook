package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
)

var increment = 1

type Contacts struct {
	Name  string
	Phone string
}

var contacts = make(map[string]Contacts)

func main() {
	router := httprouter.New()

	router.GET("/contact/:id", getContact)
	router.GET("/contacts/list", getContactList)
	router.POST("/contact", addContact)
	router.PUT("/contact/:id", updateContact)
	router.DELETE("/contact/:id", deleteContact)

	fmt.Println("Server is listening...")
	log.Fatal(http.ListenAndServe(":8080", router))

}

func getContact(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// проверить что существует элемет мапы под этим id(ключом)
	_, ok := contacts[ps.ByName("id")]
	if ok {
		fmt.Fprintf(w, "%s - %s - %s\n", ps.ByName("id"), contacts[ps.ByName("id")].Name, contacts[ps.ByName("id")].Phone)
		return
	}
	fmt.Fprintf(w, "Данный контакт отсутствует!")
}

func getContactList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// в случае путой мапы выдать сообщени что список пуст

	if len(contacts) != 0 {
		for i := range contacts {
			fmt.Fprintf(w, "%s - %s - %s\n", i, contacts[i].Name, contacts[i].Phone)
		}
		return
	}
	fmt.Fprintf(w, "Список контактов пуст!")
}

func updateContact(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// изменнение только существующих элементов!!! ключ должен быть в map[]
	// _ , ok := contacts[i]
	_, ok := contacts[ps.ByName("id")]
	if !ok {

	}
	contact := Contacts{}
	err := json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		w.Write([]byte(err.Error()))
		log.Println(err.Error())
		return
	}

	contacts[ps.ByName("id")] = contact
	fmt.Fprintf(w, "Пользователь изменен %s - %s - %s\n", ps.ByName("id"), contact.Name, contact.Phone)
}

func addContact(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// проверить что id еще не занят, если занят попробовать еще 10 раз
	log.Println("addContact")
	contact := Contacts{}
	err := json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		w.Write([]byte(err.Error()))
		log.Println(err.Error())
		return
	}
	_, ok := contacts[strconv.Itoa(increment)]
	if ok {
		for i := 0; i < 10; i++ {
			increment++
			_, ok := contacts[strconv.Itoa(increment)]
			if !ok {
				break
			}
		}
		if ok {
			fmt.Fprintf(w, "Не удалось добавить контакт , мы сделали все что могли!!!")
		}
	}
	//добавляем контакт в map[]
	contacts[strconv.Itoa(increment)] = contact

	fmt.Fprintf(w, "Добавлена запись,\n ИД: %v\n Имя: %s\n Телефон: %s!\n", increment, contact.Name, contact.Phone)
	increment++
}

func deleteContact(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// проверить что такой id существует
	_, ok := contacts[ps.ByName("id")]
	if ok {
		delete(contacts, ps.ByName("id"))
		fmt.Fprintf(w, "Контакт удален!")
		return
	}
	fmt.Fprintf(w, "Данный контакт не найден!!!")
}
