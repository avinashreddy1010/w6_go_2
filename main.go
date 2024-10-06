package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Define the structure for storing student details
type Student struct {
	ID      int    `json:"id"`      // Unique ID for each student
	Name    string `json:"name"`    // Student's name
	Program string `json:"program"` // Program they are enrolled in
	College string `json:"college"` // College name
}

// Global slice to store multiple student records
var students []Student

// Function to add a new student to the list
func createStudent(w http.ResponseWriter, r *http.Request) {
	// Allow only POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "This method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse incoming JSON request into a Student struct
	var newStudent Student
	json.NewDecoder(r.Body).Decode(&newStudent)

	// Assign an ID to the new student and add to the list
	newStudent.ID = len(students) + 1
	students = append(students, newStudent)

	// Set content type and return the created student as a response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newStudent)
}

// Function to get the list of all students
func getStudents(w http.ResponseWriter, r *http.Request) {
	// Allow only GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "This method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set content type and return all students as a response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

// Function to get a specific student by their ID
func getStudent(w http.ResponseWriter, r *http.Request) {
	// Extract the student ID from the URL
	idStr := strings.TrimPrefix(r.URL.Path, "/students/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	// Search for the student by ID and return if found
	for _, student := range students {
		if student.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(student)
			return
		}
	}

	// Return an error if the student is not found
	http.Error(w, "Student not found", http.StatusNotFound)
}

// Function to update student details
func updateStudent(w http.ResponseWriter, r *http.Request) {
	// Extract the student ID from the URL
	idStr := strings.TrimPrefix(r.URL.Path, "/students/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	// Allow only PUT requests
	if r.Method != http.MethodPut {
		http.Error(w, "This method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Find the student by ID and update their details
	for i, student := range students {
		if student.ID == id {
			students = append(students[:i], students[i+1:]...) // Remove the old entry
			var updatedStudent Student
			json.NewDecoder(r.Body).Decode(&updatedStudent)
			updatedStudent.ID = id                      // Retain the original ID
			students = append(students, updatedStudent) // Add the updated entry
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedStudent)
			return
		}
	}

	// Return an error if the student is not found
	http.Error(w, "Student not found", http.StatusNotFound)
}

// Function to delete a student by their ID
func deleteStudent(w http.ResponseWriter, r *http.Request) {
	// Extract the student ID from the URL
	idStr := strings.TrimPrefix(r.URL.Path, "/students/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	// Allow only DELETE requests
	if r.Method != http.MethodDelete {
		http.Error(w, "This method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Find the student by ID and remove them from the list
	for i, student := range students {
		if student.ID == id {
			students = append(students[:i], students[i+1:]...) // Remove the student
			w.WriteHeader(http.StatusNoContent)                // Respond with 204 No Content
			return
		}
	}

	// Return an error if the student is not found
	http.Error(w, "Student not found", http.StatusNotFound)
}

// Main function to handle all incoming requests
func handleRequests(w http.ResponseWriter, r *http.Request) {
	// Clean up the path to remove leading or trailing slashes
	path := strings.Trim(r.URL.Path, "/")

	// Handle different endpoints and methods
	switch {
	case path == "students" && r.Method == http.MethodGet:
		getStudents(w, r) // Get all students
	case path == "students" && r.Method == http.MethodPost:
		createStudent(w, r) // Add a new student
	case strings.HasPrefix(path, "students/"):
		if r.Method == http.MethodGet {
			getStudent(w, r) // Get a specific student by ID
		} else if r.Method == http.MethodPut {
			updateStudent(w, r) // Update a student by ID
		} else if r.Method == http.MethodDelete {
			deleteStudent(w, r) // Delete a student by ID
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

// Main function to start the server
func main() {
	http.HandleFunc("/", handleRequests)         // Handle all incoming requests
	log.Fatal(http.ListenAndServe(":8080", nil)) // Start the server on port 8080
}
