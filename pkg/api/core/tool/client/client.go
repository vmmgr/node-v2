package client

import (
	"bytes"
	"github.com/vmmgr/node/pkg/api/core/tool/config"
	"github.com/vmmgr/node/pkg/api/core/tool/hash"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func Post(url, body string) error {
	client := &http.Client{}
	client.Timeout = time.Second * 5

	//body, _ := json.Marshal(controller.Chat{Err: data.Err, CreatedAt: data.CreatedAt, UserID: data.UserID,
	//	GroupID: data.GroupID, Admin: data.Admin, Message: data.Message})

	//Header部分
	header := http.Header{}
	header.Set("Content-Length", "10000")
	header.Add("Content-Type", "application/json")
	header.Add("TOKEN_1", config.Conf.Controller.Auth.Token1)
	header.Add("TOKEN_2", hash.Generate(config.Conf.Controller.Auth.Token2+config.Conf.Controller.Auth.Token3))

	//リクエストの作成
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return err
	}

	req.Header = header

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func Get(url, body string) error {
	client := &http.Client{}
	client.Timeout = time.Second * 5

	//Header部分
	header := http.Header{}
	header.Set("Content-Length", "10000")
	header.Add("Content-Type", "application/json")
	header.Add("TOKEN_1", config.Conf.Controller.Auth.Token1)
	header.Add("TOKEN_2", hash.Generate(config.Conf.Controller.Auth.Token2+config.Conf.Controller.Auth.Token3))

	//リクエストの作成
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return err
	}

	req.Header = header

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func Put(url, body string) error {
	client := &http.Client{}
	client.Timeout = time.Second * 5

	//Header部分
	header := http.Header{}
	header.Set("Content-Length", "10000")
	header.Add("Content-Type", "application/json")
	header.Add("TOKEN_1", config.Conf.Controller.Auth.Token1)
	header.Add("TOKEN_2", hash.Generate(config.Conf.Controller.Auth.Token2+config.Conf.Controller.Auth.Token3))

	//リクエストの作成
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return err
	}

	req.Header = header

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func Delete(url, body string) error {
	client := &http.Client{}
	client.Timeout = time.Second * 5

	//Header部分
	header := http.Header{}
	header.Set("Content-Length", "10000")
	header.Add("Content-Type", "application/json")
	header.Add("TOKEN_1", config.Conf.Controller.Auth.Token1)
	header.Add("TOKEN_2", hash.Generate(config.Conf.Controller.Auth.Token2+config.Conf.Controller.Auth.Token3))

	//リクエストの作成
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return err
	}

	req.Header = header

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
