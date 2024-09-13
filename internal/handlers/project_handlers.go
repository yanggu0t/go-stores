package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yanggu0t/go-rdbms-practice/internal/models"
	"github.com/yanggu0t/go-rdbms-practice/internal/services"
	"github.com/yanggu0t/go-rdbms-practice/internal/utils"
)

type ProjectHandler struct {
	ProjectService *services.ProjectService
}

func NewProjectHandler(projectService *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{ProjectService: projectService}
}

func (h *ProjectHandler) GetAllProjects(c *gin.Context) {
	projects, err := h.ProjectService.GetAllProjects()
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, "error", "error_get_projects", nil)
		return
	}
	utils.Response(c, http.StatusOK, "success", "success_get_projects", projects)
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		utils.Response(c, http.StatusBadRequest, "error", "error_invalid_request_body", nil)
		return
	}

	// 生成 ProjectID
	projectID, err := models.GenerateID()
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, "error", "error_generate_project_id", nil)
		return
	}
	project.ProjectID = projectID

	if err := h.ProjectService.CreateProject(&project); err != nil {
		utils.Response(c, http.StatusInternalServerError, "error", "error_create_project", nil)
		return
	}
	utils.Response(c, http.StatusCreated, "success", "success_create_project", project)
}

func (h *ProjectHandler) GetProject(c *gin.Context) {
	id := c.Param("id")
	project, err := h.ProjectService.GetProjectByID(id)
	if err != nil {
		utils.Response(c, http.StatusNotFound, "error", "error_project_not_found", nil)
		return
	}
	utils.Response(c, http.StatusOK, "success", "success_get_project", project)
}

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id := c.Param("id")
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		utils.Response(c, http.StatusBadRequest, "error", "error_invalid_request_body", nil)
		return
	}

	if err := h.ProjectService.UpdateProject(id, &project); err != nil {
		utils.Response(c, http.StatusInternalServerError, "error", "error_update_project", nil)
		return
	}
	utils.Response(c, http.StatusOK, "success", "success_update_project", project)
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id := c.Param("id")
	if err := h.ProjectService.DeleteProject(id); err != nil {
		utils.Response(c, http.StatusInternalServerError, "error", "error_delete_project", nil)
		return
	}
	utils.Response(c, http.StatusOK, "success", "success_delete_project", nil)
}
