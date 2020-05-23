package models

import (
	//"github.com/jinzhu/gorm"
	//"gitlab.sysroot.ovh/technoservs/microservices/game-servers/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Offer struct {
	UUID       string      `json:"uuid" bson:"uuid,omitempty"`
	Name       string      `json:"name" bson:"name"`
	OfferTypes []OfferType `json:"offer_types" bson:"offer_types"`
}

type OfferType struct {
	UUID   string  `json:"uuid" bson:"uuid,omitempty"`
	Name   string  `json:"name" bson:"name"`
	Type   string  `json:"type" bson:"type"`
	Models []Model `json:"models" bson:"models"`
}

type Model struct {
	UUID string `json:"uuid" bson:"uuid,omitempty"`
	Name string `json:"name" bson:"name"`

	PlayersSlots int64     `json:"players_slots" bson:"players_slots"`
	PluginsLimit int64     `json:"plugins_limit" bson:"plugins_limit"`
	RAMMemory    int64     `json:"ram_memory" bson:"ram_memory"`
	DiskSpace    DiskSpace `json:"disk_space" bson:"disk_space"`
	Port         Port      `json:"port" bson:"port"`

	CustomDomainAddress  bool               `json:"custom_domain_address" bson:"custom_domain_address"`
	DedicatedIPAddress   DedicatedIPAddress `json:"dedicated_ip_address" bson:"dedicated_ip_address"`
	ModdedServersAllowed bool               `json:"modded_servers_allowed" bson:"modded_servers_allowed"`
	AutomatedBackups     bool               `json:"automated_backups" bson:"automated_backups"`
	PlannedTasks         bool               `json:"planned_tasks" bson:"planned_tasks"`

	PrioritizedSupport bool    `json:"prioritized_support" bson:"prioritized_support"`
	SLA                float32 `json:"sla" bson:"sla"`
	MonthlyPrice       float32 `json:"monthly_price" bson:"monthly_price"`
	HourlyPrice        float32 `json:"hourly_price" bson:"hourly_price"`
}

type DiskSpace struct {
	CurrentDiskSpace int64 `json:"current_disk_space" bson:"current_disk_space"`
	Min              int64 `json:"min" bson:"min"`
	Max              int64 `json:"max" bson:"max"`

	MonthlyPrice float32 `json:"monthly_price" bson:"monthly_price"`
	HourlyPrice  float32 `json:"hourly_price" bson:"hourly_price"`
}

type DedicatedIPAddress struct {
	DedicatedIp bool `json:"dedicated_ip" bson:"dedicated_ip"`

	MonthlyPrice float32 `json:"monthly_price" bson:"monthly_price"`
	HourlyPrice  float32 `json:"hourly_price" bson:"hourly_price"`
}

type Port struct {
	CurrentNbPort int64 `json:"current_nb_port" bson:"current_nb_port"`
	Min           int64 `json:"min" bson:"min"`
	Max           int64 `json:"max" bson:"max"`

	MonthlyPrice float32 `json:"monthly_price" bson:"monthly_price"`
	HourlyPrice  float32 `json:"hourly_price" bson:"hourly_price"`
}


func (offer *Offer) Create() map[string]interface{} {
	offerdb := GetMongoDB().Collection("offer")

	result, err := offerdb.InsertOne(ctx, offer)
	if err != nil {
		// return errors.Wrap(err, "message")
		return map[string]interface{}{
			"result": err,
		}
	}

	uuid := result.InsertedID.(primitive.ObjectID)

	filter := bson.M{"_id": uuid}

	data := Offer{
		UUID: uuid.Hex(),
		Name:       offer.Name,
		OfferTypes: offer.OfferTypes,
	}

	// update uuid of the new offer
	update := offerdb.FindOneAndUpdate(ctx, filter, bson.M{"$set": data}, options.FindOneAndUpdate().SetReturnDocument(1))

	decoded := Offer{}
	err = update.Decode(&decoded)

	if err != nil {
		return map[string]interface{}{
			"result": err,
		}
	}

	return map[string]interface{} {
		"result": offer,
	}
}

func (offer *Offer) Update(uuid string) map[string]interface{} {
	offerdb := GetMongoDB().Collection("offer")

	filter := bson.M{"uuid": uuid}

	// Result is the BSON encoded result
	// To return the updated document instead of original we have to add options.
	result := offerdb.FindOneAndUpdate(ctx, filter, bson.M{"$set": offer}, options.FindOneAndUpdate().SetReturnDocument(1))

	// Decode result and write it to 'decoded'
	decoded := Offer{}
	err := result.Decode(&decoded)

	if err != nil {
		return map[string]interface{}{
			"result": err,
		}
	}

	return map[string]interface{} {
		"result": offer,
	}
}

func GetOfferList() map[string]interface{} {
	iter, err := GetMongoDB().Collection("offer").Find(ctx, bson.D{})
	if err != nil {
		return map[string]interface{}{
			"result": err,
		}
	}

	items := make([]Offer, 0)
	for iter.Next(ctx) {
		item := Offer{}
		if err := iter.Decode(&item); err != nil {
			return map[string]interface{}{
				"result": err,
			}
		}
		items = append(items, item)
	}

	return map[string]interface{}{
		"result": items,
	}
}

func GetOffer(uuid string) map[string]interface{} {
	offerdb := GetMongoDB().Collection("offer")

	result := offerdb.FindOne(ctx, bson.M{"uuid": uuid})

	data := Offer{}

	//decode and write to data
	if err := result.Decode(&data); err != nil {
		return map[string]interface{}{
			"result": err,
		}
	}

	return map[string]interface{}{
		"uuid":       data.UUID,
		"name":       data.Name,
		"offer_types": data.OfferTypes,
	}
}

func Delete(uuid string) map[string]interface{} {
	offerdb := GetMongoDB().Collection("offer")

	_, err := offerdb.DeleteOne(ctx, bson.M{"uuid": uuid})

	if err != nil {
		return map[string]interface{}{
			"result": err,
		}
	}

	return map[string]interface{}{}
}


