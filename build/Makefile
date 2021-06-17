-include Makefile.options
#####################################################################################
dist_dir=$(CURDIR)/data
executable_name=mtservice
main_dir=../cmd/$(executable_name)
service=intelektikalt/music-transcription-service
version=0.1
go_build_cmd=CGO_ENABLED=0 go build
commit_count=$(shell git rev-list --count HEAD)
music_rec_dir?=../music
#####################################################################################
$(dist_dir):
	mkdir -p $@
$(dist_dir)/%: $(music_rec_dir)/% | $(dist_dir)
	cp $< $@
$(dist_dir)/$(executable_name): | $(dist_dir)
	cd $(main_dir) && $(go_build_cmd) -o $(dist_dir)/$(executable_name) 
#####################################################################################
dfiles=$(dist_dir)/$(executable_name) $(dist_dir)/musicrec $(dist_dir)/libaudiorecognizer.so
build: $(dfiles)
########### DOCKER ##################################################################
tag=$(service):$(version).$(commit_count)

dbuild: $(dfiles)
	docker build -t $(tag) ./

dpush: dbuild
	docker push $(tag)
#####################################################################################
clean:
	rm -rf $(dist_dir)

.PHONY:
	clean build dbuild dpush
