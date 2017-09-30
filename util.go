package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func readObjectString(object io.ReadCloser) string {
	defer object.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(object)
	return buf.String()
}
func readObjectBytes(object io.ReadCloser) []byte {
	defer object.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(object)
	return buf.Bytes()
}

var r *rand.Rand // Rand for this package.

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomString(size int) string {
	const chars = "23456789abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
	result := ""
	for i := 0; i < size; i++ {
		index := r.Intn(len(chars))
		result += chars[index : index+1]
	}
	return result
}

func httpPost(url string, requestBody interface{}) ([]byte, error) {
	b, err := json.Marshal(requestBody)
	if err != nil {
		log.Println("json err:", err)
		return nil, err
	}

	body := bytes.NewBuffer([]byte(b))
	log.Println("url:", url)
	resp, err := http.Post(url, "application/json;charset=utf-8", body)
	log.Println("resp:", resp, ",err:", err)
	if err != nil {
		return nil, err
	}

	respBody := readObjectBytes(resp.Body)
	return respBody, nil
}

func httpGet(url string) ([]byte, error) {
	log.Println("url:", url)
	resp, err := http.Get(url)
	log.Println("resp:", resp, ",err:", err)
	if err != nil {
		return nil, err
	}

	respBody := readObjectBytes(resp.Body)
	return respBody, nil
}
