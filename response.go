package main

type Host string

type Response struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}
