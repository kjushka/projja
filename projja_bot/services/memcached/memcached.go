package memcached

import (
	"github.com/bradfitz/gomemcache/memcache"
	"fmt"
	"strings"
)

var (
	Memcached = memcache.New("127.0.0.1:11211")
	cashedTime int32 = 600
)

// Кешируем выбраного пользователя, данные хранятся следующим образом
// ключ: имяПользователяРаботающегоСботом_member значение: имя выбранного пользователя_
func CacheMember(whoSelectedName string, selectedUser string) {
	key := fmt.Sprintf("%s_member", whoSelectedName)
	Memcached.Set(&memcache.Item{Key: key, Value: []byte(selectedUser), Expiration: cashedTime})
}

// Кешируем выбранный проект
// TODO переделать так, чтобы два аргумента не были одной строкой
func CacheProject(userName string, projectId string, projectName string) {
	key := fmt.Sprintf("%s_poject", userName)
	Memcached.Set(&memcache.Item{Key: key, Value: []byte(projectId + " " + projectName), Expiration: cashedTime})
}

func GetSelectedMember(userName string) (string, error) {
	member, err := Memcached.Get(fmt.Sprintf("%s_member", userName))
	if err != nil {
		return "", err
	}

	return string(member.Value), nil
}

func GetSelectedProject(userName string) (string, string, error) {
	projectForAdd, err := Memcached.Get(fmt.Sprintf("%s_poject", userName))
	if err != nil {
		return "", "", err
	}
	
	args := strings.Split(string(projectForAdd.Value), " ")
	projectId := args[0]
	projectName := args[1]

	return projectId, projectName, nil
}

