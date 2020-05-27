package models

import (
	"github.com/docker/docker/api/types"
	"github.com/jinzhu/gorm"
	"gitlab.sysroot.ovh/technoservs/microservices/game-servers/utils"
)

type Docker struct {
	UserId     string `json:"user_id"`
	Game       string `json:"game"`
	ServerName string `json:"server_name"`
}

type DockerDelete struct {
	UserId      string `json:"user_id"`
	ContainerId string `json:"container_id"`
}

type DockerStore struct {
	gorm.Model
	IdDocker     string              `json:"id_docker"`
	Game         string              `json:"game"`
	ServerName   string              `json:"server_name"`
	ServerStatus string              `json:"server_status"`
	UserId       uint                `json:"user_id"` //The user that this id belongs to
	Settings     types.ContainerJSON `json:"settings"`
}

type DockerList struct {
	UserId string `json:"user_id"`
}

func (docker *DockerStore) Validate() (map[string]interface{}, bool) {

	if docker.IdDocker == "" {
		return utils.Message(false, "Docker container id can't be null"), false
	}
	//All the required parameters are present
	return utils.Message(true, "success"), true
}

func (docker *DockerStore) Create() map[string]interface{} {

	if resp, ok := docker.Validate(); !ok {
		return resp
	}

	GetDB().Create(docker)

	resp := utils.Message(true, "success")
	resp["docker"] = docker
	return resp
}

func (docker *DockerStore) Update() map[string]interface{} {

	if resp, ok := docker.Validate(); !ok {
		return resp
	}

	GetDB().Update(docker)

	resp := utils.Message(true, "success")
	resp["docker"] = docker
	return resp
}

func ListDockerByUserId(id uint) *Docker {
	docker := &Docker{}
	err := GetDB().Table("docker").Where("id = ?", id).First(docker).Error
	if err != nil {
		return nil
	}
	return docker
}

func UserServers(id uint) *[]DockerStore {
	dockers := &[]DockerStore{}

	err := GetDB().Table("docker_stores").Where("user_id = ?", id).Find(dockers).Error
	if err != nil {
		return nil
	}
	return dockers
}

func RemoveContainer(user_id uint, docker_id string) map[string]interface{} {
	docker := &DockerStore{}

	err := GetDB().Table("docker_stores").Where("id_docker = ?", docker_id).Delete(docker).Error
	if err != nil {
		return map[string]interface{}{
			"error": err.Error,
		}
	}

	return map[string]interface{}{}
}

func (docker *DockerStore) UpdateServerStatus(status string) map[string]interface{} {

	err := GetDB().First(docker).Error

	if err != nil {
		return map[string]interface{}{
			"error": err.Error,
		}
	}

	docker.ServerStatus = status

	GetDB().Save(docker)

	return nil
}
