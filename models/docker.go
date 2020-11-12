package models

import (
	"github.com/TechnoServsEIP/Backend/utils"
	"github.com/docker/docker/api/types"
	"github.com/jinzhu/gorm"
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
	IdDocker     string               `json:"id_docker"`
	Game         string               `json:"game"`
	ServerName   string               `json:"server_name"`
	ServerStatus string               `json:"server_status"`
	UserId       uint                 `json:"user_id"` //The user that this id belongs to
	LimitPlayers int64                `json:"limit_players"`
	Settings     *types.ContainerJSON `json:"settings"`
}

type DockerList struct {
	UserId string `json:"user_id"`
}

type GameServer struct {
	UserId      string `json:"user_id"`
	ContainerId string `json:"container_id"`
	ServerName  string `json:"server_name"`
}

type DockerLimitPlayers struct {
	UserId       string `json:"user_id"`
	ContainerId  string `json:"container_id"`
	LimitPlayers int64  `json:"limit_players"`
}

type DockerLimitPlayersUserServers struct {
	UserId       string `json:"user_id"`
	LimitPlayers int64  `json:"limit_players"`
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

func OneUserServer(id uint, docker_id string) *DockerStore {
	docker := &DockerStore{}

	err := GetDB().Table("docker_stores").Where("id_docker = ?", docker_id).First(docker).Error
	if err != nil {
		return nil
	}
	return docker
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

func (docker *DockerStore) UpdateServerStatus(status string) error {

	err := GetDB().Table("docker_stores").Where("id_docker = ?", docker.IdDocker).First(docker).Error

	if err != nil {
		return err
	}

	docker.ServerStatus = status

	GetDB().Save(docker)

	return nil
}

func (docker *DockerStore) UpdateGameServer(gameSrv *GameServer) error {

	err := GetDB().Table("docker_stores").Where("id_docker = ?", docker.IdDocker).First(docker).Error

	if err != nil {
		return err
	}

	docker.ServerName = gameSrv.ServerName

	GetDB().Save(docker)

	return nil
}

func ListAllDockers() *[]DockerStore {
	dockers := &[]DockerStore{}

	err := GetDB().Table("docker_stores").Find(dockers).Error
	if err != nil {
		return nil
	}
	return dockers
}

/*
 * Change the limit number players on the server
 */
func (docker *DockerStore) ChangeLimitPlayer(limit int64) error {
	err := GetDB().Table("docker_stores").Where("id_docker = ?", docker.IdDocker).First(docker).Error

	if err != nil {
		return err
	}

	docker.LimitPlayers = limit

	GetDB().Save(docker)

	return nil
}

/*
 * Get the limit number players on the server
 */
func (docker *DockerStore) GetLimitPlayer() (int64, error) {
	err := GetDB().Table("docker_stores").Where("id_docker = ?", docker.IdDocker).First(docker).Error

	if err != nil {
		return 20, err
	}

	if docker.LimitPlayers == 0 {
		docker.LimitPlayers = 20
		GetDB().Save(docker)
	}

	return docker.LimitPlayers, nil
}
