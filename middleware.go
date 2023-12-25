package day02

type Middleware func(next HandleFunc) HandleFunc