package models

import (  
    "fmt"
    "os"
    "bufio"
    "encoding/json"
    "io/ioutil"
    "unicode"
    "strings"
    "reflect"

    "github.com/TechnoServsEIP/Backend/utils"
)

type properties struct {
    Choices         []string
    Description     string
    Min             int
    Max             int
    Name            string
    DataType        string
    DefaultValue    string
    Value           string
}

type ServerProperties struct {
	SpawnProtection          properties `json:"SpawnProtection"`
    GeneratorSettings        properties `json:"GeneratorSettings"`
    ForceGamemode            properties `json:"ForceGamemode"`
    AllowNether              properties `json:"AllowNether"`
    EnforceWhitelist         properties `json:"EnforceWhitelist"`
    Gamemode                 properties `json:"Gamemode"`
    PlayerIdleTimeout        properties `json:"PlayerIdleTimeout"`
    Difficulty               properties `json:"Difficulty"`
    SpawnMonsters            properties `json:"SpawnMonsters"`
    OpPermissionLevel        properties `json:"OpPermissionLevel"`
    Pvp                      properties `json:"Pvp"`
    SnooperEnabled           properties `json:"SnooperEnabled"`
    LevelType                properties `json:"LevelType"`
    Hardcore                 properties `json:"Hardcore"`
    EnableStatus             properties `json:"EnableStatus"`
    EnableCommandBlock       properties `json:"EnableCommandBlock"`
    MaxPlayers               properties `json:"MaxPlayers"`
    MaxWorldSize             properties `json:"MaxWorldSize"`
    FunctionPermissionLevel  properties `json:"FunctionPermissionLevel"`
    SpawnNpcs                properties `json:"SpawnNpcs"`
    AllowFlight              properties `json:"AllowFlight"`
    LevelName                properties `json:"LevelName"`
    ViewDistance             properties `json:"ViewDistance"`
    ResourcePack             properties `json:"ResourcePack"`
    SpawnAnimals             properties `json:"SpawnAnimals"`
    WhitLlist                properties `json:"WhitLlist"`
    GenerateStructures       properties `json:"GenerateStructures"`
    OnlineMode               properties `json:"OnlineMode"`
    LevelSeed                properties `json:"LevelSeed"`
    Motd                     properties `json:"Motd"`
}

func kebabToCamelCase(str string) string {
    tmp := ""
    upper := true

    for _, char := range str {
        
        if (char != '-' && !upper) {
            tmp += string(char)
        } else {
            upper = true
        }

        if (upper && char != '-') {
            tmp += string(unicode.ToUpper(char))
            upper = false
        }
    }

    return tmp
}

func ServerPropertiesByServerId(containerId string) *ServerProperties {
    /*
     * Load JSON width default description of server.properties
     */
    file, _ := ioutil.ReadFile("./models/serverProperties.json")

    data := ServerProperties{}

    _ = json.Unmarshal([]byte(file), &data)

    /*
     * Read server.properties file
     */
    tmpFolder := "./models/" + containerId

    if err := utils.CreateTmpFolder(tmpFolder); err != nil {
        fmt.Println(err)
		return nil
    }

    serverFile := containerId + ":/data/server.properties"

	if err := utils.DockerCopy(serverFile, tmpFolder); err != nil {
        fmt.Println(err)
		return nil
    }

    fileServerProp, err := os.Open(tmpFolder + "/server.properties")
    if err != nil {
		fmt.Println(err)
		return nil
    }
    defer fileServerProp.Close()

	/*
     * Set property value
     */
    scanner := bufio.NewScanner(fileServerProp)
    for scanner.Scan() {
        property := strings.Split(scanner.Text(), "=")

        if (len(property) > 1) {
            if (reflect.ValueOf(&data).Elem().FieldByName(kebabToCamelCase(property[0])) != reflect.Value{}) {
                reflect.ValueOf(&data).Elem().FieldByName(kebabToCamelCase(property[0])).FieldByName("Value").SetString(property[1])
            }
        }
    }

    if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil
	}

    /*
     * Remove server.property
     */
    if err := utils.DeleteTmpFolder(tmpFolder); err != nil {
        fmt.Println(err)
		return nil
    }

	return &data
}