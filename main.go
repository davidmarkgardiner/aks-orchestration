package main

import (
    "embed"
    "encoding/json"
    "fmt"
    "html/template"
    "log"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "sync"
)

//go:embed templates
var templates embed.FS

// Script represents a single script to be executed
type Script struct {
    Name        string `json:"name"`
    Description string `json:"description"`
}

// ExecutionStatus represents the current state of script execution
type ExecutionStatus struct {
    CurrentScript string   `json:"current_script"`
    Completed    []string `json:"completed"`
    Errors       []string `json:"errors"`
    InProgress   bool     `json:"in_progress"`
    mu           sync.Mutex
}

var (
    status = &ExecutionStatus{
        Completed:  make([]string, 0),
        Errors:     make([]string, 0),
        InProgress: false,
    }

    // Define your script sequence here
    scriptSequence = []Script{
        {Name: "1_aks_prerequisites.sh", Description: "Installing AKS Prerequisites"},
        {Name: "2_create_aks_cluster.sh", Description: "Creating AKS Cluster"},
        {Name: "3_setup_networking.sh", Description: "Configuring Network"},
        {Name: "4_install_monitoring.sh", Description: "Setting up Monitoring"},
    }
)

func (s *ExecutionStatus) reset() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.CurrentScript = ""
    s.Completed = make([]string, 0)
    s.Errors = make([]string, 0)
    s.InProgress = false
}

func (s *ExecutionStatus) updateCurrentScript(script string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.CurrentScript = script
}

func (s *ExecutionStatus) addCompleted(script string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.Completed = append(s.Completed, script)
}

func (s *ExecutionStatus) addError(err string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.Errors = append(s.Errors, err)
}

func executeScripts() {
    if status.InProgress {
        return
    }

    status.mu.Lock()
    status.InProgress = true
    status.mu.Unlock()

    // Reset status
    status.reset()
    status.InProgress = true

    // Create scripts directory if it doesn't exist
    if err := os.MkdirAll("scripts", 0755); err != nil {
        status.addError(fmt.Sprintf("Failed to create scripts directory: %v", err))
        status.InProgress = false
        return
    }

    for _, script := range scriptSequence {
        status.updateCurrentScript(script.Description)
        scriptPath := filepath.Join("scripts", script.Name)

        // Check if script exists
        if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
            status.addError(fmt.Sprintf("Script not found: %s", script.Name))
            continue
        }

        // Make script executable
        if err := os.Chmod(scriptPath, 0755); err != nil {
            status.addError(fmt.Sprintf("Failed to make script executable: %s", script.Name))
            continue
        }

        // Execute script
        cmd := exec.Command(scriptPath)
        if output, err := cmd.CombinedOutput(); err != nil {
            status.addError(fmt.Sprintf("Error executing %s: %v\nOutput: %s", script.Name, err, output))
            break
        }

        status.addCompleted(script.Description)
    }

    status.mu.Lock()
    status.InProgress = false
    status.CurrentScript = ""
    status.mu.Unlock()
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFS(templates, "templates/index.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, nil)
}

func handleStart(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    go executeScripts()

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "started"})
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    status.mu.Lock()
    defer status.mu.Unlock()
    json.NewEncoder(w).Encode(status)
}

func main() {
    // Set up routes
    http.HandleFunc("/", handleIndex)
    http.HandleFunc("/start", handleStart)
    http.HandleFunc("/status", handleStatus)

    // Start server
    log.Println("Server starting on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
} 