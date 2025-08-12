package services

import (
    "encoding/json"
    "errors"
    "os"
    "taxi-service/models"
)

const dummyUserFile = "./data/dummy_users.json"

func readUsers() ([]models.DummyUser, error) {
    // Criar diretório se não existir
    if err := os.MkdirAll("./data", 0755); err != nil {
        return nil, err
    }

    file, err := os.Open(dummyUserFile)
    if err != nil {
        if os.IsNotExist(err) {
            return []models.DummyUser{}, nil
        }
        return nil, err
    }
    defer file.Close()

    var users []models.DummyUser
    err = json.NewDecoder(file).Decode(&users)
    if err != nil {
        return nil, err
    }
    return users, nil
}

func writeUsers(users []models.DummyUser) error {
    data, err := json.MarshalIndent(users, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(dummyUserFile, data, 0644)
}

func ListDummyUser() ([]models.DummyUser, error) {
    return readUsers()
}

func GetDummyUser(id int) (models.DummyUser, error) {
    users, err := readUsers()
    if err != nil {
        return models.DummyUser{}, err
    }
    for _, user := range users {
        if int(user.ID) == id {
            return user, nil
        }
    }
    return models.DummyUser{}, errors.New("user not found")
}

func CreateDummyUser(user *models.DummyUser) error {
    users, err := readUsers()
    if err != nil {
        return err
    }
    // Assign a new ID
    var maxID int = 0
    for _, u := range users {
        if u.ID > maxID {
            maxID = u.ID
        }
    }
    user.ID = maxID + 1
    users = append(users, *user)
    return writeUsers(users)
}

func UpdateDummyUser(id int, updateData *models.DummyUser) (models.DummyUser, error) {
    users, err := readUsers()
    if err != nil {
        return models.DummyUser{}, err
    }
    for i, user := range users {
        if int(user.ID) == id {
            if updateData.Name != "" {
                users[i].Name = updateData.Name
            }
            if updateData.Email != "" {
                users[i].Email = updateData.Email
            }
            // Add other fields as needed
            err = writeUsers(users)
            return users[i], err
        }
    }
    return models.DummyUser{}, errors.New("user not found")
}

func DeleteDummyUser(id int) error {
    users, err := readUsers()
    if err != nil {
        return err
    }
    newUsers := []models.DummyUser{}
    found := false
    for _, user := range users {
        if int(user.ID) != id {
            newUsers = append(newUsers, user)
        } else {
            found = true
        }
    }
    if !found {
        return errors.New("user not found")
    }
    return writeUsers(newUsers)
}