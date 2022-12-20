package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type TaskAPI interface {
	GetTask(w http.ResponseWriter, r *http.Request)
	CreateNewTask(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
	UpdateTaskCategory(w http.ResponseWriter, r *http.Request)
}

type taskAPI struct {
	taskService service.TaskService
}

func NewTaskAPI(taskService service.TaskService) *taskAPI {
	return &taskAPI{taskService}
}

func (t *taskAPI) GetTask(w http.ResponseWriter, r *http.Request) {
	Uid, err := strconv.Atoi(r.Context().Value("id").(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	Tid, _ := strconv.Atoi(r.URL.Query().Get("task_id"))
	if Tid == 0 {
		tasks, err := t.taskService.GetTasks(r.Context(), Uid)
		fmt.Println("test")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			json.NewEncoder(w).Encode(entity.NewErrorResponse("internal server error"))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tasks)
		return
	}

	task, err := t.taskService.GetTaskByID(r.Context(), Tid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("internal server error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (t *taskAPI) CreateNewTask(w http.ResponseWriter, r *http.Request) {
	Uid := r.Context().Value("id").(string)
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}
	this := entity.Task{}
	this.Title = task.Title
	this.Description = task.Description
	this.CategoryID = task.CategoryID
	this.UserID, _ = strconv.Atoi(Uid)
	if task.Title == "" || task.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}
	tasks, err := t.taskService.StoreTask(r.Context(), &this)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("internal server error"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": Uid,
		"task_id": tasks.ID,
		"message": "success create new task",
	})
}

func (t *taskAPI) DeleteTask(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value("id")
	Tid, _ := strconv.Atoi(r.URL.Query().Get("task_id"))
	err := t.taskService.DeleteTask(r.Context(), Tid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("internal server error"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": uid,
		"task_id": Tid,
		"message": "success delete task",
	})
}

func (t *taskAPI) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest
	Uid := r.Context().Value("id").(string)
	Tid, _ := strconv.Atoi(r.URL.Query().Get("task_id"))
	this := entity.Task{}
	this.Title = task.Title
	this.Description = task.Description
	this.UserID, _ = strconv.Atoi(Uid)
	this.ID = Tid
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}
	tasks, err := t.taskService.UpdateTask(r.Context(), &this)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("internal server error"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": tasks.UserID,
		"task_id": Tid,
		"message": "success update task",
	})
}

func (t *taskAPI) UpdateTaskCategory(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskCategoryRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	Uid := r.Context().Value("id")

	idLogin, err := strconv.Atoi(Uid.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	var updateTask = entity.Task{
		ID:         task.ID,
		CategoryID: task.CategoryID,
		UserID:     int(idLogin),
	}

	_, err = t.taskService.UpdateTask(r.Context(), &updateTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": Uid,
		"task_id": task.ID,
		"message": "success update task category",
	})
}
