module github.com/vimeo/alog/emitter/gkelog/traceextractors/oc

go 1.18

require (
	github.com/vimeo/alog/v3 v3.5.0
	go.opencensus.io v0.23.0
)

require github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect

replace github.com/vimeo/alog/v3 => ../../../../
