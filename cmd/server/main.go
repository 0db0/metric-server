package main

import (
	"metric-server/config"
	"metric-server/internal/app"
)

func main() {
	cfg := config.MustLoad()
	app.Run(cfg)
}

//10,Alloc,gauge,"{""raw"": 872648}"
//11,BuckHashSys,gauge,"{""raw"": 3394}"
//12,Frees,gauge,"{""raw"": 168}"
//13,GCCPUFraction,gauge,"{""raw"": null}"
//14,GCSys,gauge,"{""raw"": 6904848}"
//15,HeapAlloc,gauge,"{""raw"": 872648}"
