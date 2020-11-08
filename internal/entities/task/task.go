package task

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	jsonFile = "todo/todo.json"
)

type Task struct {
	ID       int    `json:"id"`
	Task     string `json:"task"`
	RemindTime string `json:"remindtime"`
}

type Handler struct {
	JSONPath string
	tasks    []*Task
}

func NewHandler() (*Handler, error) {
	t := new(Handler)

	if err := t.exploreJSONPath(); err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadFile(t.JSONPath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bytes, &t.tasks); err != nil {
		return nil, err
	}

	return t, nil
}

func (h *Handler) GetTask(id int) *Task {
	if id >= len(h.tasks) {
		return nil
	}
	return h.tasks[id-1]
}

func (h *Handler) GetTasks() []*Task {
	return h.tasks
}

func (h *Handler) AppendTask(t *Task) {
	h.tasks = append(h.tasks, t)
	h.align()
}

func (h *Handler) exploreJSONPath() error {
	dataHome := os.Getenv("XDG_DATA_HOME")
	var jsonPath string
	var homeDir, _ = os.UserHomeDir()
	if dataHome != "" {
		jsonPath = filepath.Join(dataHome, jsonFile)
	} else {
		jsonPath = filepath.Join(homeDir, ".local/share/", jsonFile)
	}

	if err := createJSONFile(jsonPath); err != nil {
		return err
	}

	h.JSONPath = jsonPath
	return nil
}

func createJSONFile(path string) error {
	if _, err := os.Stat(filepath.Dir(path)); err != nil {
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}
	}

	if _, err := os.Stat(path); err != nil {
		if err := writeInitialSample(path); err != nil {
			return err
		}
	}

	return nil
}

func writeInitialSample(path string) error {
	tasks := &[]*Task{
		{
			ID:       1,
			Task:     "deleting or modifying this task is your first TODO",
			RemindTime: "2099/01/01 00:00",
		},
	}

	bytes, err := json.Marshal(tasks)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, bytes, 0644); err != nil {
		return err
	}

	return nil
}

func (h *Handler) Write() error {
	bytes, err := json.Marshal(&h.tasks)
	if err != nil {
		return nil
	}

	if err := ioutil.WriteFile(h.JSONPath, bytes, 0644); err != nil {
		return err
	}

	return nil
}

func (h *Handler) Remove(id int) {
	if id > len(h.tasks) {
		return
	}

	h.tasks = append(h.tasks[:id], h.tasks[id+1:]...)
	h.align()
}

func (h *Handler) align() {
	for i, t := range h.tasks {
		t.ID = i + 1
	}
}
