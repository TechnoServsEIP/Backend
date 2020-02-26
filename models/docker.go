package models

import (
	"github.com/jinzhu/gorm"
	"oauth2server/utils"
)

type Docker struct {
	Game string `json:"game"`
}

type DockerDelete struct {
	UserId string `json:"user_id"`
	ContainerId string`json:"container_id"`
}

type DockerStore struct {
	gorm.Model
	Game string `json:"game"`
	Id string `json:"id"`
	UserId uint `json:"user_id"` //The user that this id belongs to
}

func (docker *DockerStore) Validate() (map[string] interface{}, bool) {

	if docker.Id == "" {
		return utils.Message(false, "Docker container id can't be null"), false
	}
	//All the required parameters are present
	return utils.Message(true, "success"), true
}

func (docker *DockerStore) Create() map[string] interface{} {

	if resp, ok := docker.Validate(); !ok {
		return resp
	}

	GetDB().Create(docker)

	resp := utils.Message(true, "success")
	resp["docker"] = docker
	return resp
}

func (docker *DockerStore) Update() map[string] interface{} {

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