package models

import (  
    "fmt"
    "os"
    "bufio"
    "encoding/json"
    "io/ioutil"
    "unicode"
    "strconv"
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

type UpdateServerProperties struct {
    UserId                  string  `json:"user_id"`
    ContainerId             string  `json:"container_id"`

	SpawnProtection          int    `json:"SpawnProtection"`
    GeneratorSettings        string `json:"GeneratorSettings"`
    ForceGamemode            bool   `json:"ForceGamemode"`
    AllowNether              bool   `json:"AllowNether"`
    EnforceWhitelist         bool   `json:"EnforceWhitelist"`
    Gamemode                 string `json:"Gamemode"`
    PlayerIdleTimeout        int    `json:"PlayerIdleTimeout"`
    Difficulty               string `json:"Difficulty"`
    SpawnMonsters            bool   `json:"SpawnMonsters"`
    OpPermissionLevel        int    `json:"OpPermissionLevel"`
    Pvp                      bool   `json:"Pvp"`
    SnooperEnabled           bool   `json:"SnooperEnabled"`
    LevelType                string `json:"LevelType"`
    Hardcore                 bool   `json:"Hardcore"`
    EnableStatus             bool   `json:"EnableStatus"`
    EnableCommandBlock       bool   `json:"EnableCommandBlock"`
    MaxPlayers               int    `json:"MaxPlayers"`
    MaxWorldSize             int    `json:"MaxWorldSize"`
    FunctionPermissionLevel  int    `json:"FunctionPermissionLevel"`
    SpawnNpcs                bool   `json:"SpawnNpcs"`
    AllowFlight              bool   `json:"AllowFlight"`
    LevelName                string `json:"LevelName"`
    ViewDistance             int    `json:"ViewDistance"`
    ResourcePack             string `json:"ResourcePack"`
    SpawnAnimals             bool   `json:"SpawnAnimals"`
    WhitLlist                bool   `json:"WhitLlist"`
    GenerateStructures       bool   `json:"GenerateStructures"`
    OnlineMode               bool   `json:"OnlineMode"`
    LevelSeed                string `json:"LevelSeed"`
    Motd                     string `json:"Motd"`
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

/*
 * Attach defaults and new values ​​from server.properties file
 *
 * data: new value send by the client
 * tmpFolder: folder where will the server.properties
 * containerId: the containerId of the client server
 * serverFile: the path where the server.properties is into the client server
 *
 * return: map with each line of the server.properties file updating
 */
func getValueToUpdate(data UpdateServerProperties, tmpFolder string, containerId string, serverFile string) map[int]interface{} {
    // create tmp folder
    if err := utils.CreateTmpFolder(tmpFolder); err != nil {
        fmt.Println(err)
		return map[int]interface{}{}
    }

    // copy server.properties from the client container into tmpfolder
	if err := utils.DockerCopy(serverFile, tmpFolder); err != nil {
        fmt.Println(err)
        
        if err := utils.DeleteTmpFolder(tmpFolder); err != nil {
            fmt.Println(err)
            return map[int]interface{}{}
        }

		return map[int]interface{}{}
    }

    // read the server.properties and attach defaults and new values
    fileServerProp, err := os.Open(tmpFolder + "/server.properties")
    if err != nil {
        fmt.Println(err)
        return map[int]interface{}{}
    }
    defer fileServerProp.Close()

    lines := make(map[int]interface{})
    i := 0

    scanner := bufio.NewScanner(fileServerProp)
    for scanner.Scan() {
        property := strings.Split(scanner.Text(), "=")

        if (len(property) > 1) {
            if (reflect.ValueOf(&data).Elem().FieldByName(kebabToCamelCase(property[0])) != reflect.Value{}) {
                switch v := reflect.ValueOf(&data).Elem().FieldByName(kebabToCamelCase(property[0])); v.Kind() {
                case reflect.Bool:
                    lines[i] = property[0] + "=" + strconv.FormatBool(v.Bool())
                case reflect.String:
                    lines[i] = property[0] + "=" + v.String()
                case reflect.Int:
                    lines[i] = property[0] + "=" + strconv.FormatInt(v.Int(), 10)
                }
            } else {
                lines[i] = property[0] + "=" + property[1]
            }
        }
        i++
    }

    if err := scanner.Err(); err != nil {
        fmt.Println(err)
        return map[int]interface{}{}
    }

    return lines
}

/*
 * Re create a new server.properties file with the new value sended by the client
 *
 * data: new value send by the client
 * containerId: the containerId of the client server
 *
 * return: error or nil if no error 
 */
func CreateNewServerProperties(data UpdateServerProperties, containerId string) error {
    tmpFolder := "./models/" + containerId
    serverFile := containerId + ":/data/server.properties"
    values := getValueToUpdate(data, tmpFolder, containerId, serverFile)

    f, err := os.Create(tmpFolder + "/server.properties")
    if err != nil {
        fmt.Println(err)
        f.Close()
        return err
    }

    for _, v := range values {
        fmt.Fprintln(f, v)
        if err != nil {
            fmt.Println(err)
            return err
        }
    }
    err = f.Close()
    if err != nil {
        fmt.Println(err)
        return err
    }

    // Copy the new server.properties file into the client container
    if err := utils.DockerCopy(tmpFolder + "/server.properties", containerId + ":/data"); err != nil {
        fmt.Println(err)
        
        if err := utils.DeleteTmpFolder(tmpFolder); err != nil {
            fmt.Println(err)
            return err
        }

		return err
    }

    // remove the tmp folder
    if err := utils.DeleteTmpFolder(tmpFolder); err != nil {
        fmt.Println(err)
        return err
    }

    return nil
}