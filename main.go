package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Student struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Program string `json:"program"`
}

var students []Student

// Create a New Student
func createStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var newStudent Student
	json.NewDecoder(r.Body).Decode(&newStudent)
	newStudent.ID = len(students) + 1
	students = append(students, newStudent)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newStudent)
}

// Read All Students
func getStudents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

// Read a Single Student by ID
func getStudent(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/students/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	for _, student := range students {
		if student.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(student)
			return
		}
	}
	http.Error(w, "Student not found", http.StatusNotFound)
}

// Update a Student
func updateStudent(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/students/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	for i, student := range students {
		if student.ID == id {
			students = append(students[:i], students[i+1:]...)
			var updatedStudent Student
			json.NewDecoder(r.Body).Decode(&updatedStudent)
			updatedStudent.ID = id
			students = append(students, updatedStudent)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedStudent)
			return
		}
	}
	http.Error(w, "Student not found", http.StatusNotFound)
}

// Delete a Student
func deleteStudent(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/students/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	for i, student := range students {
		if student.ID == id {
			students = append(students[:i], students[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Student not found", http.StatusNotFound)
}

// Route Handling
func handleRequests(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")

	switch {
	case path == "students" && r.Method == http.MethodGet:
		getStudents(w, r)
	case path == "students" && r.Method == http.MethodPost:
		createStudent(w, r)
	case strings.HasPrefix(path, "students/"):
		if r.Method == http.MethodGet {
			getStudent(w, r)
		} else if r.Method == http.MethodPut {
			updateStudent(w, r)
		} else if r.Method == http.MethodDelete {
			deleteStudent(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

// Main Function
func main() {
	http.HandleFunc("/", handleRequests)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
